package database

import (
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     *config.Config
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	log.Printf("🔌 Connecting to RabbitMQ: %s", maskPassword(cfg.RabbitMQ.URL))
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Printf("❌ Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("❌ Failed to create channel: %v", err)
		conn.Close()
		return nil, err
	}

	queues := []string{
		cfg.RabbitMQ.Queues.ClickEvents,
	}

	for _, queueName := range queues {
		log.Printf("📋 Declaring queue: %s", queueName)
		_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Printf("❌ Failed to declare queue %s: %v", queueName, err)
			ch.Close()
			conn.Close()
			return nil, err
		}
	}

	log.Printf("✅ RabbitMQ connected successfully")
	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Config:     cfg,
	}, nil
}

func maskPassword(url string) string {
	// Simple masking for security
	return "amqp://***:***@..."
}

func (r *RabbitMQ) Close() error {
	log.Println("🔌 Closing RabbitMQ connection...")

	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("⚠️  Failed to close channel: %v", err)
		}
	}

	if r.Connection != nil {
		if err := r.Connection.Close(); err != nil {
			log.Printf("⚠️  Failed to close connection: %v", err)
			return err
		}
	}

	log.Println("✅ RabbitMQ connection closed successfully")
	return nil
}
