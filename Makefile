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
build: ## Build the production Docker image
	docker-compose -f docker/docker-compose.yml build

run: ## Run the application in production mode
	docker-compose -f docker/docker-compose.yml up -d

stop: ## Stop the running containers
	docker-compose -f docker/docker-compose.yml down

clean: ## Remove containers, networks, and images
	docker-compose -f docker/docker-compose.yml down --rmi all --volumes --remove-orphans

logs: ## Show application logs
	docker-compose -f docker/docker-compose.yml logs -f docstore-api

### -->>>>>>>>>>> Production with nginx
prod: ## Run with nginx reverse proxy
	docker-compose -f docker/docker-compose.yml --profile production up -d

prod-stop: ## Stop production setup with nginx
	docker-compose -f docker/docker-compose.yml --profile production down
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

# Local build
build-local: ## Build the application locally
	go build -mod=mod -o docstore-api ./src

run-local: build-local ## Build and run the application locally
	./docstore-api

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
########## Utility commands
########## ########## ########## ########## 

swagger: ## Generate swagger documentation
	docker-compose -f docker/docker-compose.dev.yml exec docstore-api-dev swag init

shell: ## Get shell access to running container
	docker-compose -f docker/docker-compose.yml exec docstore-api sh

health: ## Check application health
	curl -f http://localhost:8080/api/v1/documents || echo "Service is not healthy"

# Docker image management
image-size: ## Show Docker image size
	docker images docstore-api

prune: ## Clean up unused Docker resources
	docker system prune -f