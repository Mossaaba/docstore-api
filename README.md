# Document Store API

A RESTful document storage API built with Go, featuring a clean layered architecture and comprehensive testing. This project demonstrates CRUD operations with proper concurrency control, HTTP endpoints, Swagger documentation, and monitoring stack using Grafana + Prometheus + Loki.

## Features

- **RESTful API**: HTTP endpoints for document management
- **JWT Authentication**: Secure API access with JSON Web Tokens
- **Layered Architecture**: Clean separation of concerns (Models → Services → Controllers)
- **Thread-Safe Operations**: Uses RWMutex for concurrent read/write access
- **Full CRUD Operations**: Create, Read, Update (PUT/PATCH), Delete, and List documents
- **Swagger Documentation**: Auto-generated API documentation with JWT security
- **Comprehensive Testing**: Unit tests for all layers including concurrency testing
- **Error Handling**: Proper HTTP status codes and error messages
- **Monitoring Tools**: Enabling metrics, health checks, and log collection

## Architecture

The application follows a clean 3-layer architecture:

### **Models Layer** (`models/`)
- **Document**: Core data structure
- **DocumentStore**: Thread-safe in-memory storage with full CRUD operations
- **Update Operations**: Full replacement (PUT) and partial updates (PATCH)
- Uses `sync.RWMutex` for concurrent access control

### **Services Layer** (`services/`)
- **DocumentService**: Business logic interface and implementation
- Abstracts storage operations from HTTP layer
- Handles business rules and validation

### **Controllers Layer** (`controllers/`)
- **DocumentController**: HTTP request handlers for documents
- **AuthController**: Authentication and JWT token management
- JSON serialization/deserialization
- HTTP status code management
- Swagger documentation annotations

### **Middleware Layer** (`middleware/`)
- **JWTAuthMiddleware**: JWT token validation and user context
- Token parsing and validation
- Authorization header processing

### Data Flow
```
HTTP Request → Controller → Service → Model → Storage
HTTP Response ← Controller ← Service ← Model ← Storage
```

## Testing

The project includes comprehensive tests for all layers:

### **Models Layer Tests** (`models/document_test.go`)
- CRUD operations testing
- Thread-safety and concurrency tests
- Edge cases (duplicates, not found)
- Concurrent read/write scenarios

### **Services Layer Tests** (`services/document_service_test.go`)
- Business logic validation
- Error handling verification
- Full workflow testing
- Service interface compliance

### **Controllers Layer Tests** (`controllers/document_controller_test.go`)
- HTTP endpoint testing
- JSON request/response validation
- Status code verification
- Complete API workflow tests

## Project Structure

```
docstore-api/
├── src/
│   ├── main.go                       # Application entry point
│   ├── config/
│   │   └── config.go                 # Configuration management
│   ├── models/
│   │   ├── document.go               # Document model and store
│   │   └── document_test.go          # Model layer tests
│   ├── services/
│   │   ├── document_service.go       # Business logic layer
│   │   └── document_service_test.go  # Service layer tests
│   ├── controllers/
│   │   ├── document_controller.go    # HTTP handlers
│   │   ├── auth_controller.go        # Authentication handlers
│   │   └── document_controller_test.go # Controller layer tests
│   ├── middleware/
│   │   └── jwt_middleware.go         # JWT authentication middleware
│   └── docs/                         # Swagger generated docs
├── config/
│   ├── .env.example                  # Environment template
│   ├── .env.development              # Development environment
│   └── .env.production               # Production environment
├── docker/                           # Docker configuration
├── go.mod                            # Go module dependencies
└── README.md                         # This file
```

## Quick Start

### Prerequisites
- Go 1.21 or higher (for local development)
- Docker and Docker Compose (recommended)

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd docstore-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Generate Swagger documentation:
```bash
swag init
```

4. Run the API server:
```bash
cd src
go run main.go
```

5. Use pre-commit to enhance code quality:
```bash
pre-commit install
pre-commit run --all-files
```

### Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific layer tests
go test ./models -v
go test ./services -v
go test ./controllers -v

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## API Usage Examples

### 1. Login (Get JWT Token)
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": "admin"
}
```

### 2. Create Document (Protected)
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id": "doc-1",
    "name": "My Document",
    "description": "A sample document"
  }'
```

### 3. Get Document (Protected)
```bash
curl http://localhost:8080/api/v1/documents/doc-1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. List All Documents (Protected)
```bash
curl http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Update Document - PUT (Protected)
```bash
curl -X PUT http://localhost:8080/api/v1/documents/doc-1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id": "doc-1",
    "name": "Updated Document Name",
    "description": "Updated description"
  }'
```

### 6. Partially Update Document - PATCH (Protected)
```bash
curl -X PATCH http://localhost:8080/api/v1/documents/doc-1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Only Update Name"
  }'
```

### 7. Delete Document (Protected)
```bash
curl -X DELETE http://localhost:8080/api/v1/documents/doc-1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Response Codes

- `200 OK` - Successful GET request or login
- `201 Created` - Document created successfully
- `204 No Content` - Document deleted successfully
- `400 Bad Request` - Invalid JSON or request format
- `401 Unauthorized` - Missing, invalid, or expired JWT token
- `404 Not Found` - Document not found
- `409 Conflict` - Document with ID already exists



## API Reference

### HTTP Endpoints

#### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/login` | User login (get JWT token) | No |

#### Documents (Protected)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/documents` | Create a new document | Yes |
| GET | `/api/v1/documents` | List all documents | Yes |
| GET | `/api/v1/documents/{id}` | Get document by ID | Yes |
| PUT | `/api/v1/documents/{id}` | Update entire document | Yes |
| PATCH | `/api/v1/documents/{id}` | Partially update document | Yes |
| DELETE | `/api/v1/documents/{id}` | Delete document by ID | Yes |

### Document Structure
```json
{
    "id": "string",
    "name": "string",
    "description": "string"
}
```

### Authentication

The API uses JWT (JSON Web Tokens) for authentication. All document endpoints require a valid JWT token.

#### Configuration

The API uses environment variables for configuration with support for multiple environments:


## Environment Configuration

### Development Setup
```bash
cp config/.env.example config/.env.development
# Edit config/.env.development with your development settings
```

### Production Setup
```bash
cp config/.env.example config/.env.production
# Edit config/.env.production with secure production values
```

### Environment File Loading Priority
The configuration system loads files in this order:

1. **Environment Variables** (highest priority)
   - Set via Docker, system, or command line

2. **Environment-Specific File**
   - `config/.env.development` (when `APP_ENV=development`)
   - `config/.env.production` (when `APP_ENV=production`)

3. **General Environment File** (fallback)
   - `config/.env` (for local overrides)

**Example**: In development (`APP_ENV=development`):
- Loads `config/.env.development` first
- Then `config/.env` as fallback
- Environment variables override everything

### Default Credentials (configurable via config/.env)
- **Username**: `admin` (set via `ADMIN_USERNAME`)
- **Password**: `password` (set via `ADMIN_PASSWORD`)
- **JWT Secret**: Configurable via `JWT_SECRET`



## Docker Usage

### Docker Features

- **Multi-stage build**: Optimized production image (~10MB)
- **Security**: Non-root user, minimal attack surface
- **Hot reload**: Development mode with automatic restart
- **Health checks**: Built-in container health monitoring
- **Resource limits**: CPU and memory constraints
- **Nginx proxy**: Optional reverse proxy for production
- **SSL ready**: HTTPS configuration template included

### Common Commands

```bash
make help          # Show all available commands
make build         # Build production Docker image
make clean         # Remove containers, networks, images, and volumes
make prune         # Clean up unused Docker resources
make health        # Check application health
```

### Development Commands

```bash
make dev           # Start development with hot reload (detached)
make swagger-dev   # Generate swagger documentation for dev
make docker-test   # Run tests in development container
make dev-logs      # View development logs
make dev-stop      # Stop development environment
make test          # Run tests locally
make test-coverage # Run tests with coverage locally
```

### Production Commands

```bash
make prod-build    # Build with production environment
make prod-up       # Start production with monitoring stack
make prod-logs     # Show production logs
make prod-stop     # Stop production setup
```

### Utility Commands

```bash
make shell-dev     # Get shell access to running container
make health        # Get health check of the application
make metrics       # Get metrics of the application
```

## Security

### JWT Authentication
- **Token Expiration**: 24 hours
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Header Format**: `Authorization: Bearer <token>`
- **Secret Key**: Configurable via environment variable (defaults to demo key)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT signing | `your-secret-key-change-in-production` |
| `ADMIN_USERNAME` | Admin username | `admin` |
| `ADMIN_PASSWORD` | Admin password | `password` |
| `SERVER_PORT` | Server port | `8080` |
| `APP_ENV` | Environment (development/production) | `development` |

### Production Security Notes
- **Always** change the JWT secret key in production
- Use strong, randomly generated passwords
- Set secure environment variables in your deployment
- Consider implementing refresh tokens for better security
- Add rate limiting for authentication endpoints
- Never commit `.env` files to version control

## To Do

- GitHub Actions (Build and push Docker / run tests with coverage / generate release)
- Final code review and validation
- Code packaging and distribution
