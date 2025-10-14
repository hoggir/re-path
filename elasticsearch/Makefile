.PHONY: help up down restart logs logs-es logs-kibana ps health clean test-connection kibana

COMPOSE_FILE=docker-compose.yml

help:
	@echo "🔍 Elasticsearch Docker - Makefile Commands"
	@echo ""
	@echo "Basic Commands:"
	@echo "  make up              - Start Elasticsearch"
	@echo "  make down            - Stop Elasticsearch"
	@echo "  make restart         - Restart Elasticsearch"
	@echo "  make logs            - View all logs (follow mode)"
	@echo "  make logs-es         - View Elasticsearch logs only"
	@echo "  make ps              - Show running containers"
	@echo "  make health          - Check Elasticsearch health"
	@echo ""
	@echo "Kibana Commands:"
	@echo "  make kibana          - Start Elasticsearch + Kibana"
	@echo "  make logs-kibana     - View Kibana logs"
	@echo ""
	@echo "Maintenance Commands:"
	@echo "  make clean           - Stop and remove volumes (⚠️  deletes data)"
	@echo "  make reset           - Complete reset (stop, clean, restart)"
	@echo "  make test-connection - Test Elasticsearch connection"
	@echo ""
	@echo "Info Commands:"
	@echo "  make cluster-info    - Show cluster information"
	@echo "  make indices         - List all indices"
	@echo "  make nodes           - Show node information"
	@echo ""

# Basic Operations
up:
	@echo "🚀 Starting Elasticsearch..."
	@if ! docker info > /dev/null 2>&1; then \
		echo "❌ Docker is not running!"; \
		echo ""; \
		echo "Please start Docker Desktop first:"; \
		echo "  1. Open Docker Desktop application"; \
		echo "  2. Wait until Docker is running"; \
		echo "  3. Run 'make up' again"; \
		exit 1; \
	fi
	docker-compose -f $(COMPOSE_FILE) up -d elasticsearch
	@echo "✅ Elasticsearch started!"
	@echo "📍 Elasticsearch: http://localhost:9200"
	@echo "⏳ Waiting for Elasticsearch to be ready..."
	@sleep 5
	@make health

down:
	@echo "🛑 Stopping Elasticsearch..."
	docker-compose -f $(COMPOSE_FILE) down
	@echo "✅ Elasticsearch stopped!"

restart:
	@echo "🔄 Restarting Elasticsearch..."
	@make down
	@sleep 2
	@make up

# Logs
logs:
	@echo "📜 Viewing logs (Ctrl+C to exit)..."
	docker-compose -f $(COMPOSE_FILE) logs -f

logs-es:
	@echo "📜 Elasticsearch logs (Ctrl+C to exit)..."
	docker-compose -f $(COMPOSE_FILE) logs -f elasticsearch

logs-kibana:
	@echo "📜 Kibana logs (Ctrl+C to exit)..."
	docker-compose -f $(COMPOSE_FILE) logs -f kibana

# Status
ps:
	@echo "📋 Running containers:"
	@docker-compose -f $(COMPOSE_FILE) ps

health:
	@echo "🏥 Checking Elasticsearch health..."
	@curl -s http://localhost:9200/_cluster/health?pretty 2>/dev/null || echo "❌ Elasticsearch not responding"

# Kibana
kibana:
	@echo "🚀 Starting Elasticsearch + Kibana..."
	docker-compose -f $(COMPOSE_FILE) --profile kibana up -d
	@echo "✅ Elasticsearch + Kibana started!"
	@echo "📍 Elasticsearch: http://localhost:9200"
	@echo "📍 Kibana: http://localhost:5601"
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@make health

# Maintenance
clean:
	@echo "⚠️  This will delete all Elasticsearch data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "🧹 Cleaning up..."; \
		docker-compose -f $(COMPOSE_FILE) down -v; \
		echo "✅ Cleanup complete!"; \
	else \
		echo "❌ Cancelled"; \
	fi

reset:
	@echo "🔄 Resetting Elasticsearch..."
	docker-compose -f $(COMPOSE_FILE) down -v
	@sleep 2
	@make up

test-connection:
	@echo "🧪 Testing Elasticsearch connection..."
	@echo ""
	@echo "1. Cluster Health:"
	@curl -s http://localhost:9200/_cluster/health?pretty
	@echo ""
	@echo "2. Cluster Info:"
	@curl -s http://localhost:9200 | python3 -m json.tool 2>/dev/null || curl -s http://localhost:9200

# Info Commands
cluster-info:
	@echo "ℹ️  Cluster Information:"
	@curl -s http://localhost:9200 | python3 -m json.tool 2>/dev/null || curl -s http://localhost:9200

indices:
	@echo "📑 Indices:"
	@curl -s "http://localhost:9200/_cat/indices?v"

nodes:
	@echo "🖥️  Nodes:"
	@curl -s "http://localhost:9200/_cat/nodes?v"

# Development helpers
create-test-index:
	@echo "🧪 Creating test index..."
	@curl -X PUT "localhost:9200/test-index" -H 'Content-Type: application/json' -d'{"settings": {"number_of_shards": 1,"number_of_replicas": 0}}'
	@echo ""
	@echo "✅ Test index created!"

index-test-doc:
	@echo "📝 Indexing test document..."
	@curl -X POST "localhost:9200/test-index/_doc" -H 'Content-Type: application/json' -d'{"user": "test_user","message": "Hello from Makefile!","timestamp": "'$$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}'
	@echo ""
	@echo "✅ Document indexed!"

search-test:
	@echo "🔎 Searching test index..."
	@curl -s "localhost:9200/test-index/_search?pretty"

delete-test-index:
	@echo "🗑️  Deleting test index..."
	@curl -X DELETE "localhost:9200/test-index"
	@echo ""
	@echo "✅ Test index deleted!"

# Quick tests
test: create-test-index index-test-doc search-test
	@echo ""
	@echo "✅ All tests completed!"
	@echo "💡 Run 'make delete-test-index' to clean up"
