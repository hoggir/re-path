package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQRPCService interface {
	Call(ctx context.Context, queueName string, payload interface{}, timeout time.Duration) ([]byte, error)
}

type rabbitMQRPCService struct {
	rabbitmq *database.RabbitMQ
	logger   logger.Logger
}

func NewRabbitMQRPCService(rabbitmq *database.RabbitMQ, log logger.Logger) RabbitMQRPCService {
	return &rabbitMQRPCService{
		rabbitmq: rabbitmq,
		logger:   log,
	}
}

func (s *rabbitMQRPCService) Call(ctx context.Context, queueName string, payload interface{}, timeout time.Duration) ([]byte, error) {
	// Serialize payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Declare a temporary exclusive queue for receiving response
	replyQueue, err := s.rabbitmq.Channel.QueueDeclare(
		"",    // name (empty = auto-generated)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %w", err)
	}

	// Generate unique correlation ID for this request
	correlationID := uuid.New().String()

	// Register consumer for reply queue
	msgs, err := s.rabbitmq.Channel.Consume(
		replyQueue.Name, // queue
		"",              // consumer (empty = auto-generated)
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	// Publish request message
	err = s.rabbitmq.Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			ReplyTo:       replyQueue.Name,
			Body:          body,
			DeliveryMode:  amqp.Transient,
			Timestamp:     time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish RPC request: %w", err)
	}

	s.logger.DebugContext(ctx, "RPC request sent", "queue", queueName, "correlationId", correlationID)

	select {
	case msg := <-msgs:
		if msg.CorrelationId == correlationID {
			// log.Printf("📦 Response body: %s", string(msg.Body))
			return msg.Body, nil
		}
		return nil, fmt.Errorf("received message with mismatched correlation ID")

	case <-time.After(timeout):
		return nil, fmt.Errorf("RPC call timeout after %v", timeout)

	case <-ctx.Done():
		return nil, fmt.Errorf("RPC call cancelled: %w", ctx.Err())
	}
}
