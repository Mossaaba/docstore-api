# Document Store API

A RESTful document storage API built with Go, featuring a clean layered architecture and comprehensive testing. This project demonstrates CRUD operations with proper concurrency control, HTTP endpoints, and Swagger documentation.

## Features

- **RESTful API**: HTTP endpoints for document management
- **Layered Architecture**: Clean separation of concerns (Models → Services → Controllers)
- **Thread-Safe Operations**: Uses RWMutex for concurrent read/write access
- **Full CRUD Operations**: Create, Read, Update (PUT/PATCH), Delete, and List documents
- **Swagger Documentation**: Auto-generated API documentation
- **Comprehensive Testing**: Unit tests for all layers including concurrency testing
- **Error Handling**: Proper HTTP status codes and error messages

## Project Structure

```
docstore-api/
├── main.go                           # Application entry point
├── models/
│   ├── document.go                   # Document model and store
│   └── document_test.go              # Model layer tests
├── services/
│   ├── document_service.go           # Business logic layer
│   └── document_service_test.go      # Service layer tests
├── controllers/
│   ├── document_controller.go        # HTTP handlers
│   └── document_controller_test.go   # Controller layer tests
├── docs/                             # Swagger generated docs
├── go.mod                            # Go module dependencies
└── README.md                         # This file
```

## Quick Start

### Prerequisites
- Go 1.19 or higher

### Installation

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
go run main.go
```

5. Access the API:
- **API Base URL**: http://localhost:8080/api/v1
- **Swagger UI**: http://localhost:8080/swagger/index.html

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests by layer
go test ./models
go test ./services  
go test ./controllers

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

## API Reference

### HTTP Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/documents` | Create a new document |
| GET | `/api/v1/documents` | List all documents |
| GET | `/api/v1/documents/{id}` | Get document by ID |
| PUT | `/api/v1/documents/{id}` | Update entire document |
| PATCH | `/api/v1/documents/{id}` | Partially update document |
| DELETE | `/api/v1/documents/{id}` | Delete document by ID |

### Document Structure
```json
{
    "id": "string",
    "name": "string", 
    "description": "string"
}
```

### Example API Calls

#### Create Document
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "id": "doc-1",
    "name": "My Document",
    "description": "A sample document"
  }'
```

#### Get Document
```bash
curl http://localhost:8080/api/v1/documents/doc-1
```

#### List All Documents
```bash
curl http://localhost:8080/api/v1/documents
```

#### Update Document (PUT)
```bash
curl -X PUT http://localhost:8080/api/v1/documents/doc-1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "doc-1",
    "name": "Updated Document Name",
    "description": "Updated description"
  }'
```

#### Partially Update Document (PATCH)
```bash
curl -X PATCH http://localhost:8080/api/v1/documents/doc-1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Only Update Name"
  }'
```

#### Delete Document
```bash
curl -X DELETE http://localhost:8080/api/v1/documents/doc-1
```

### Response Codes

- `200 OK` - Successful GET request
- `201 Created` - Document created successfully
- `204 No Content` - Document deleted successfully
- `400 Bad Request` - Invalid JSON or request format
- `404 Not Found` - Document not found
- `409 Conflict` - Document with ID already exists

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
- **DocumentController**: HTTP request handlers
- JSON serialization/deserialization
- HTTP status code management
- Swagger documentation annotations

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

### Test Commands

```bash
# Run all tests


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




# To DO : 
1- Catching the parametter in inout rahther then specifiey name and description
2-  

