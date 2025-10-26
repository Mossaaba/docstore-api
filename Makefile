# Docker commands for the Document Store API

.PHONY: help build run stop dev test  docker-test logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
########## ########## ########## ##########
########## Production Docker commands :
########## ########## ########## ##########

prod-build: ## Build production image with .env.production
	docker-compose -f docker/compose.prod.yml build

prod: ## Run with nginx reverse proxy using .env.production
	docker-compose -f docker/compose.prod.yml up -d

prod-logs: ## Show production logs
	docker-compose -f docker/compose.prod.yml logs -f docstore-api

prod-stop: ## Stop production setup with nginx
	docker-compose -f docker/compose.prod.yml down

prod-up: ## Start production with monitoring stack
	docker-compose -f docker/compose.prod.yml up -d

prod-monitoring: ## Start only monitoring services (Grafana, Prometheus, Loki)
	docker-compose -f docker/compose.prod.yml up -d grafana prometheus loki promtail

prod-monitoring-logs: ## Show monitoring services logs
	docker-compose -f docker/compose.prod.yml logs -f grafana prometheus loki promtail


########## ########## ########## ##########
########## Devlopement Docker commands :
########## ########## ########## ##########

# Development Docker commands
dev: ## Run the application in development mode with hot reload
	docker-compose -f docker/compose.dev.yml up --build -d
	@echo "Development environment started in background"
	@echo "Use 'make dev-logs' to view logs"
	@echo "Use 'make dev-stop' to stop the environment"

dev-stop: ## Stop development containers
	docker-compose -f docker/compose.dev.yml down

dev-logs: ## Show development logs
	docker-compose -f docker/compose.dev.yml logs -f docstore-api-dev

########## ########## ########## ##########
########## Devlopement Local commands :
########## ########## ########## ##########
# Local build
build-local: ## Build the application locally
	go build -mod=mod -o docstore-api ./src

run-local: build-local ## Build and run the application locally
	APP_ENV=development ./docstore-api

run-local-prod: build-local ## Build and run the application locally in production mode
	APP_ENV=production ./docstore-api

########## ########## ########## ##########
########## Testing commands :
########## ########## ########## ##########
# Testing
test: ## Run tests locally
	cd src && go test -v ./... -mod=mod

test-coverage: ## Run tests with coverage report and open in browser
	cd src && GOFLAGS="-mod=mod" go test -cover ./...
	cd src && GOFLAGS="-mod=mod" go test -coverprofile=../coverage.out ./...
	GOFLAGS="-mod=mod" go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Opening coverage report in browser..."
	open coverage.html

docker-test: ## Run tests in Docker container
	docker-compose -f docker/compose.dev.yml exec docstore-api-dev go test -v ./...

########## ########## ########## ##########
########## Swagger
########## ########## ########## ##########

swagger-dev: ## Generate swagger documentation for development environment
	docker-compose -f docker/compose.dev.yml exec docstore-api-dev sh -c "cd /app/src && swag init -g main.go --output docs --instanceName dev"

swagger-prod-rebuild: ## Rebuild production image with updated Swagger docs
	docker-compose -f docker/compose.prod.yml build --no-cache

########## ########## ########## ##########
########## Pre-commit commands
########## ########## ########## ##########

precommit-install: ## Install pre-commit hooks
	pre-commit install

precommit-run: ## Run pre-commit on all files
	pre-commit run --all-files

precommit-update: ## Update pre-commit hooks to latest versions
	pre-commit autoupdate

########## ########## ########## ##########
########## Utility commands
########## ########## ########## ##########

shell-dev: ## Get shell access to running container in dev
	docker-compose -f docker/compose.dev.yml exec docstore-api-dev sh

health: ## Check application health --> if not https
	curl -f http://localhost:8080/health || echo "Service is not healthy"

metrics: ## Check application metrics
	curl -f http://localhost:8080/metrics || echo "Metrics endpoint not available"

monitoring-status: ## Check status of all monitoring services
	@echo "=== Monitoring Services Status ==="
	@echo "Grafana (3000): $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 || echo 'DOWN')"
	@echo "Prometheus (9090): $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:9090 || echo 'DOWN')"
	@echo "Loki (3100): $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3100/ready || echo 'DOWN')"
	@echo "API Health (8080): $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/health || echo 'DOWN')"
	@echo "API Metrics (8080): $$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/metrics || echo 'DOWN')"
