# ============================================================================
# Environment Variables
# ============================================================================

# PostgreSQL
export POSTGRES_DB = manifestodb
export POSTGRES_USER = manifesto
export POSTGRES_PASSWORD = supersecret
export POSTGRES_HOST = localhost
export POSTGRES_PORT = 5432

# Redis
export REDIS_HOST = localhost
export REDIS_PORT = 6379
export REDIS_PASSWORD =
export REDIS_DB = 0

# Server
export PORT = 8080
export ENVIRONMENT = development
export READ_TIMEOUT = 10s
export WRITE_TIMEOUT = 10s
export SHUTDOWN_TIMEOUT = 30s
export CORS_ALLOWED_ORIGINS = *
export ENCRYPTION_KEY = 4676f229448b490f8228ea4290f3e543

# Database (DB_* map to POSTGRES_* by default)
export DB_HOST = $(POSTGRES_HOST)
export DB_PORT = $(POSTGRES_PORT)
export DB_USER = $(POSTGRES_USER)
export DB_PASSWORD = $(POSTGRES_PASSWORD)
export DB_NAME = $(POSTGRES_DB)
export DB_SSLMODE = disable
export DB_MAX_OPEN_CONNS = 25
export DB_MAX_IDLE_CONNS = 5
export DB_CONN_MAX_LIFETIME = 5m

# Auth (JWT)
export JWT_SECRET = development-supersecret-32-characters-min-123456
export ACCESS_TOKEN_TTL = 15m
export REFRESH_TOKEN_TTL = 168h
export JWT_ISSUER = manifesto

# OAuth (optional; leave CLIENT_ID/SECRET empty if unused)
export GOOGLE_CLIENT_ID = -
export GOOGLE_CLIENT_SECRET = -
export GOOGLE_REDIRECT_URL = http://localhost:5173/auth/callback/?provider=google
export MICROSOFT_CLIENT_ID = -
export MICROSOFT_CLIENT_SECRET = -
export MICROSOFT_REDIRECT_URL = http://localhost:8080/auth/callback/microsoft

# Build PostgreSQL connection string
CONN_STRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
CONTAINER_NAME = manifesto-postgres
REDIS_CONTAINER_NAME = manifesto-redis

# ============================================================================
# Help
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================================================
# Development
# ============================================================================

.PHONY: dev
dev: ## Run the development server
	go mod tidy
	go run ./cmd

.PHONY: build
build: ## Build the application binary
	go mod tidy
	go build -o bin/server ./cmd

.PHONY: prod
prod: build ## Build and run production server
	./bin/server

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint: ## Run linter
	golangci-lint run

# ============================================================================
# Docker - All Services
# ============================================================================

.PHONY: up
up: ## Start all services (postgres + redis)
	docker compose up -d --remove-orphans
	@echo "Waiting for services to be ready..."
	@sleep 3
	@make health

.PHONY: down
down: ## Stop all services
	docker compose down

.PHONY: down-v
down-v: ## Stop all services and remove volumes
	docker compose down -v

.PHONY: restart
restart: down up ## Restart all services

.PHONY: logs
logs: ## Show logs for all services
	docker compose logs -f

.PHONY: health
health: ## Check health of all services
	@echo "Checking PostgreSQL..."
	@docker exec $(CONTAINER_NAME) pg_isready -U $(POSTGRES_USER) || echo "❌ PostgreSQL not ready"
	@echo "Checking Redis..."
	@docker exec $(REDIS_CONTAINER_NAME) redis-cli ping || echo "❌ Redis not ready"

# ============================================================================
# Docker - PostgreSQL Only
# ============================================================================

.PHONY: postgres-up
postgres-up: ## Start PostgreSQL only
	docker compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@docker exec $(CONTAINER_NAME) pg_isready -U $(POSTGRES_USER)

.PHONY: postgres-down
postgres-down: ## Stop PostgreSQL
	docker compose stop postgres

.PHONY: postgres-restart
postgres-restart: postgres-down postgres-up ## Restart PostgreSQL

.PHONY: postgres-logs
postgres-logs: ## Show PostgreSQL logs
	docker compose logs -f postgres

.PHONY: postgres-shell
postgres-shell: ## Open shell in PostgreSQL container
	docker exec -it $(CONTAINER_NAME) /bin/sh

# ============================================================================
# Docker - Redis Only
# ============================================================================

.PHONY: redis-up
redis-up: ## Start Redis only
	docker compose up -d redis
	@echo "Waiting for Redis to be ready..."
	@sleep 2
	@docker exec $(REDIS_CONTAINER_NAME) redis-cli ping

.PHONY: redis-down
redis-down: ## Stop Redis
	docker compose stop redis

.PHONY: redis-restart
redis-restart: redis-down redis-up ## Restart Redis

.PHONY: redis-logs
redis-logs: ## Show Redis logs
	docker compose logs -f redis

.PHONY: redis-cli
redis-cli: ## Open Redis CLI
	@if [ -z "$(REDIS_PASSWORD)" ]; then \
		docker exec -it $(REDIS_CONTAINER_NAME) redis-cli; \
	else \
		docker exec -it $(REDIS_CONTAINER_NAME) redis-cli -a $(REDIS_PASSWORD); \
	fi

.PHONY: redis-shell
redis-shell: ## Open shell in Redis container
	docker exec -it $(REDIS_CONTAINER_NAME) /bin/sh

.PHONY: redis-flush
redis-flush: ## Flush all Redis data
	@echo "⚠️  Flushing all Redis data..."
	@docker exec $(REDIS_CONTAINER_NAME) redis-cli FLUSHALL
	@echo "✅ Redis data flushed"

# ============================================================================
# Database Operations
# ============================================================================

.PHONY: psql
psql: ## Open psql in the PostgreSQL container
	docker exec -it $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: conn
conn: ## Show the PostgreSQL connection string
	@echo $(CONN_STRING)

.PHONY: migrate
migrate: ## Run database migrations
	@echo "Running migrations..."
	@if [ ! -f migrations/001_genesis.sql ]; then \
		echo "❌ Migration file not found: migrations/001_genesis.sql"; \
		exit 1; \
	fi
	docker exec -i $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < migrations/001_genesis.sql
	@echo "✅ Migrations completed"

.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: name is required. Usage: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="migrations/$${timestamp}_$(name).sql"; \
	echo "-- Migration: $(name)" > $$filename; \
	echo "-- Created at: $$(date)" >> $$filename; \
	echo "" >> $$filename; \
	echo "-- Add your SQL here" >> $$filename; \
	echo "" >> $$filename; \
	echo "✅ Created migration: $$filename"

.PHONY: seed
seed: ## Seed test data
	@echo "Seeding test data..."
	@if [ -f migrations/seed_test_data.sql ]; then \
		docker exec -i $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < migrations/seed_test_data.sql; \
		echo "✅ Test data seeded"; \
	else \
		echo "⚠️  No seed file found (migrations/seed_test_data.sql)"; \
	fi

.PHONY: db-clean
db-clean: ## Clean database (drop all tables)
	@echo "⚠️  Cleaning database..."
	@read -p "Are you sure you want to drop all tables? (y/N) " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker exec -i $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"; \
		echo "✅ Database cleaned"; \
	else \
		echo "❌ Cancelled"; \
	fi

.PHONY: db-reset
db-reset: db-clean migrate seed ## Full database reset (clean + migrate + seed)
	@echo "✅ Database reset complete!"

.PHONY: db-backup
db-backup: ## Backup database to file
	@timestamp=$$(date +%Y%m%d_%H%M%S); \
	filename="backups/backup_$${timestamp}.sql"; \
	mkdir -p backups; \
	docker exec $(CONTAINER_NAME) pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) > $$filename; \
	echo "✅ Backup saved to $$filename"


# ============================================================================
# Setup & Initialization
# ============================================================================

.PHONY: setup
setup: up migrate seed ## Full setup (start services + migrate + seed)
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Services running:"
	@echo "  PostgreSQL: localhost:$(POSTGRES_PORT)"
	@echo "  Redis:      localhost:$(REDIS_PORT)"
	@echo "  Server:     make dev"

.PHONY: setup-postgres
setup-postgres: postgres-up migrate seed ## Setup PostgreSQL only
	@echo "✅ PostgreSQL setup complete!"

.PHONY: setup-redis
setup-redis: redis-up ## Setup Redis only
	@echo "✅ Redis setup complete!"

# ============================================================================
# Cleanup
# ============================================================================

.PHONY: clean
clean: down-v ## Stop services and remove volumes
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Cleanup complete"

.PHONY: clean-all
clean-all: clean ## Clean everything including Docker images
	docker compose down --rmi all --volumes --remove-orphans
	@echo "✅ Full cleanup complete"

# ============================================================================
# Utility
# ============================================================================

.PHONY: env
env: ## Show current environment variables
	@echo "PostgreSQL Configuration:"
	@echo "  POSTGRES_HOST:     $(POSTGRES_HOST)"
	@echo "  POSTGRES_PORT:     $(POSTGRES_PORT)"
	@echo "  POSTGRES_DB:       $(POSTGRES_DB)"
	@echo "  POSTGRES_USER:     $(POSTGRES_USER)"
	@echo ""
	@echo "Redis Configuration:"
	@echo "  REDIS_HOST:        $(REDIS_HOST)"
	@echo "  REDIS_PORT:        $(REDIS_PORT)"
	@echo "  REDIS_DB:          $(REDIS_DB)"
	@echo ""
	@echo "Server Configuration:"
	@echo "  PORT:              $(PORT)"
	@echo "  ENVIRONMENT:       $(ENVIRONMENT)"
	@echo ""
	@echo "Connection String:"
	@echo "  $(CONN_STRING)"

.PHONY: ps
ps: ## Show running containers
	docker compose ps

.PHONY: stats
stats: ## Show container resource usage
	docker stats $(CONTAINER_NAME) $(REDIS_CONTAINER_NAME)

# ============================================================================
# Default
# ============================================================================

.DEFAULT_GOAL := help
