package database

import (
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     *config.Config
	logger     logger.Logger
}

func NewRabbitMQ(cfg *config.Config, log logger.Logger) (*RabbitMQ, error) {
	log.Info("connecting to RabbitMQ", "url", maskPassword(cfg.RabbitMQ.URL))
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Error("failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error("failed to create channel", "error", err)
		conn.Close()
		return nil, err
	}

	queues := []string{
		cfg.RabbitMQ.Queues.ClickEvents,
		cfg.RabbitMQ.Queues.DashboardRequest,
	}

	for _, queueName := range queues {
		log.Info("declaring queue", "queue", queueName)
		_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Error("failed to declare queue", "queue", queueName, "error", err)
			ch.Close()
			conn.Close()
			return nil, err
		}
	}

	log.Info("RabbitMQ connected successfully")
	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Config:     cfg,
		logger:     log,
	}, nil
}

func maskPassword(url string) string {
	// Simple masking for security
	return "amqp://***:***@..."
}

func (r *RabbitMQ) Close() error {
	r.logger.Info("closing RabbitMQ connection")

	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			r.logger.Warn("failed to close channel", "error", err)
		}
	}

	if r.Connection != nil {
		if err := r.Connection.Close(); err != nil {
			r.logger.Warn("failed to close connection", "error", err)
			return err
		}
	}

	r.logger.Info("RabbitMQ connection closed successfully")
	return nil
}
