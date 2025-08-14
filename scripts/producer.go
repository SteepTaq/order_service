package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"wb_order_service/internal/config"
	"wb_order_service/internal/model"

	"github.com/segmentio/kafka-go"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	brokersEnv := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersEnv, ",")
	topic := getEnv("KAFKA_TOPIC", "orders")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	defer writer.Close()

	log.Println("kafka producer started")

	for i := 1; i <= 5; i++ {
		order := generateTestOrder(i)

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("failed to serialize order: %v", err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: orderJSON,
		})
		if err != nil {
			log.Printf("failed to send message: %v", err)
		} else {
			log.Printf("order %d sent: %s", i, order.OrderUID)
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("sending completed")
}

func generateTestOrder(index int) *model.Order {
	orderUID := fmt.Sprintf("test-order-%d-%d", index, time.Now().Unix())

	return &model.Order{
		OrderUID:          orderUID,
		TrackNumber:       fmt.Sprintf("WBILMTESTTRACK%d", index),
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("customer-%d", index),
		DeliveryService:   "meest",
		ShardKey:          fmt.Sprintf("%d", index%10),
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: model.Delivery{
			Name:    fmt.Sprintf("Test User %d", index),
			Phone:   fmt.Sprintf("+972000000%d", index),
			Zip:     fmt.Sprintf("263980%d", index),
			City:    "Test City",
			Address: fmt.Sprintf("Test Address %d", index),
			Region:  "Test Region",
			Email:   fmt.Sprintf("test%d@gmail.com", index),
		},
		Payment: model.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000 + rand.Intn(1000),
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   500 + rand.Intn(500),
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930 + index,
				TrackNumber: fmt.Sprintf("WBILMTESTTRACK%d", index),
				Price:       453 + rand.Intn(100),
				Rid:         fmt.Sprintf("ab4219087a764ae0btest%d", index),
				Name:        fmt.Sprintf("Test Product %d", index),
				Sale:        30,
				Size:        "0",
				TotalPrice:  317 + rand.Intn(100),
				NmID:        2389212 + index,
				Brand:       "Test Brand",
				Status:      202,
			},
		},
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
