package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
	"wb_order_service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewConsumer(brokers []string, topic string, groupID string, logger *slog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader: reader,
		logger: logger,
	}
}

func (c *Consumer) Close() error { return c.reader.Close() }

func (c *Consumer) ConsumeMessages(ctx context.Context, messageHandler func(*model.Order) error) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("failed to read message", "error", err)
				if err.Error() == "EOF" || err.Error() == "connection refused" {
					c.logger.Warn("connection error, waiting before retry", "error", err)
					time.Sleep(5 * time.Second)
					continue
				}
				continue
			}
			c.logger.Debug("message received", "offset", m.Offset, "partition", m.Partition)
			var order model.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				c.logger.Error("failed to parse JSON", "error", err, "message", string(m.Value))
				continue
			}
			if err := validate.Struct(&order); err != nil {
				c.logger.Error("failed to validate order", "error", err)
				continue
			}

			if err := messageHandler(&order); err != nil {
				c.logger.Error("failed to process order", "error", err)
				continue
			}

			c.logger.Debug("message processed successfully", "order_uid", order.OrderUID)
		}
	}
}
