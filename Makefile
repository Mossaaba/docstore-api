# Docker commands for the Document Store API

.PHONY: help build run stop clean dev test  docker-test logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
########## ########## ########## ########## 
########## Production Docker commands : 
########## ########## ########## ########## 

prod: ## Run with nginx reverse proxy using .env.production
	docker-compose -f docker/docker-compose.prod.yml up -d

prod-stop: ## Stop production setup with nginx
	docker-compose -f docker/docker-compose.prod.yml down

prod-logs: ## Show production logs
	docker-compose -f docker/docker-compose.prod.yml logs -f docstore-api

prod-build: ## Build production image with .env.production
	docker-compose -f docker/docker-compose.prod.yml build
########## ########## ########## ########## 
########## Devlopement Docker commands : 
########## ########## ########## ########## 

# Development Docker commands
dev: ## Run the application in development mode with hot reload
	docker-compose -f docker/docker-compose.dev.yml up --build -d
	@echo "Development environment started in background"
	@echo "Use 'make dev-logs' to view logs"
	@echo "Use 'make dev-stop' to stop the environment"

dev-stop: ## Stop development containers
	docker-compose -f docker/docker-compose.dev.yml down

dev-logs: ## Show development logs
	docker-compose -f docker/docker-compose.dev.yml logs -f docstore-api-dev

########## ########## ########## ########## 
########## Devlopement Local commands : 
########## ########## ########## ########## 
# Local build
build-local: ## Build the application locally
	go build -mod=mod -o docstore-api ./src

run-local: build-local ## Build and run the application locally
	./docstore-api

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
	docker-compose -f docker/docker-compose.dev.yml exec docstore-api-dev go test -v ./...

########## ########## ########## ########## 
########## Swagger
########## ########## ########## ########## 

swagger-dev: ## Generate swagger documentation for development environment
	docker-compose -f docker/docker-compose.dev.yml exec docstore-api-dev sh -c "cd /app/src && swag init -g main.go --output docs --instanceName dev"

swagger-prod: ## Generate swagger documentation for production environment (built into image)
	@echo "Production Swagger docs are generated during Docker build process"
	@echo "To regenerate, rebuild the production image with: make build"

swagger-prod-rebuild: ## Rebuild production image with updated Swagger docs
	docker-compose -f docker/docker-compose.yml build --no-cache

########## ########## ########## ########## 
########## Utility commands
########## ########## ########## ########## 

shell-dev: ## Get shell access to running container in dev
	docker-compose -f docker/docker-compose.dev.yml exec docstore-api-dev sh

health: ## Check application health
	curl -f http://localhost:8080/api/v1/documents || echo "Service is not healthy"
