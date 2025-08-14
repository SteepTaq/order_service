package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"wb_order_service/internal/cache"
	"wb_order_service/internal/config"
	"wb_order_service/internal/database"
	"wb_order_service/internal/handlers"
	"wb_order_service/internal/kafka"
	"wb_order_service/internal/model"
	"wb_order_service/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log := logger.Setup(cfg.Logger.Level)
	slog.SetDefault(log)

	log.Info("starting order service")
	db, err := database.NewDatabase(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("failed to close database", "error", err)
		}
	}()
	log.Info("database connection open")
	cache := cache.NewCache()
	log.Info("cache initialized")
	log.Info("restoring cache from database")
	orders, err := db.GetAllOrders(context.Background())
	if err != nil {
		log.Warn("failed to restore cache", "error", err)
	} else {
		for _, order := range orders {
			cache.Set(order)
		}
		log.Info("orders restored to cache", "count", len(orders))
	}

	kafkaBrokers := strings.Split(cfg.Kafka.Brokers, ",")
	consumer := kafka.NewConsumer(kafkaBrokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, log)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error("failed to close kafka consumer", "error", err)
		}
	}()
	log.Info("kafka consumer initialized")

	handler := handlers.NewHandler(cache, db)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Get("/order/{order_uid}", handler.GetOrderByID)
	router.Get("/api-docs/*", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))
	router.Handle("/*", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("http server started", "port", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start http server", "error", err)
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		log.Info("starting kafka consumer")
		if err := consumer.ConsumeMessages(ctx, func(order *model.Order) error {
			if err := db.SaveOrder(ctx, order); err != nil {
				return fmt.Errorf("failed to save order to database: %v", err)
			}
			cache.Set(order)
			log.Info("order processed and saved", "order_uid", order.OrderUID)
			return nil
		}); err != nil {
			log.Error("failed to consume messages", "error", err)
		}
	}()

	<-done
	log.Info("get signal, stop service")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown http server", "error", err)
	}
	log.Info("service stopped")
}
