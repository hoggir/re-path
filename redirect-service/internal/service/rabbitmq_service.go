package service

import (
	"context"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService interface {
	PublishClickEvent(ctx context.Context, payload []byte) error
}

type rabbitMQService struct {
	rabbitmq *database.RabbitMQ
}

func NewRabbitMQService(rabbitmq *database.RabbitMQ) RabbitMQService {
	return &rabbitMQService{
		rabbitmq: rabbitmq,
	}
}

func (s *rabbitMQService) PublishClickEvent(ctx context.Context, payload []byte) error {
	queueName := s.rabbitmq.Config.RabbitMQ.Queues.ClickEvents

	err := s.rabbitmq.Channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("📤 Published click event to queue: %s (size: %d bytes)", queueName, len(payload))
	return nil
}
