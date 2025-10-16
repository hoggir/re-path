.PHONY: rabbitmq-up rabbitmq-down rabbitmq-logs rabbitmq-restart rabbitmq-clean rabbitmq-ui

# Start RabbitMQ
up:
	@echo "🚀 Starting RabbitMQ..."
	docker compose -f docker-compose.yml up -d
	@echo "✅ RabbitMQ started"
	@echo "📊 Management UI: http://localhost:15672"
	@echo "👤 Username: repath"
	@echo "🔑 Password: repath123"

# Stop RabbitMQ
down:
	@echo "🛑 Stopping RabbitMQ..."
	docker-compose -f docker-compose.yml down
	@echo "✅ RabbitMQ stopped"

# View RabbitMQ logs
logs:
	docker-compose -f docker-compose.yml logs -f rabbitmq

# Restart RabbitMQ
restart:
	@echo "🔄 Restarting RabbitMQ..."
	docker-compose -f docker-compose.yml restart
	@echo "✅ RabbitMQ restarted"

# Clean RabbitMQ (remove volumes)
clean:
	@echo "🧹 Cleaning RabbitMQ data..."
	docker-compose -f docker-compose.yml down -v
	@echo "✅ RabbitMQ cleaned"

# Open RabbitMQ Management UI
ui:
	@echo "🌐 Opening RabbitMQ Management UI..."
	xdg-open http://localhost:15672 2>/dev/null || open http://localhost:15672 2>/dev/null || echo "Open http://localhost:15672 in your browser"

# Show RabbitMQ status
status:
	@echo "📊 RabbitMQ Status:"
	@docker ps --filter name=repath-rabbitmq --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
