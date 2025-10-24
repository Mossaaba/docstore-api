# Document Store API

A thread-safe, in-memory document storage system built with Go. This project demonstrates CRUD operations with proper concurrency control using Go's sync package.

## Features

- **Thread-Safe Operations**: Uses RWMutex for concurrent read/write access
- **CRUD Operations**: Create, Read, Update, Delete, and List documents
- **Error Handling**: Proper error messages for edge cases
- **Comprehensive Tests**: Unit tests including concurrency testing
- **Performance Benchmarks**: Benchmark tests for performance analysis

## Project Structure

```
docstore-api/
├── main.go           # Demo application
├── document.go       # Document struct and DocumentStore implementation
├── document_test.go  # Comprehensive unit tests
├── Dockerfile        # Container configuration
└── README.md         # This file
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

2. Initialize Go module:
```bash
go mod init docstore-api
```

3. Run the demo:
```bash
go run .
```

4. Run tests:
```bash
go test -v
```

5. Run benchmarks:
```bash
go test -bench=.
```

## API Reference

### Document Structure
```go
type Document struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

### DocumentStore Methods

#### Create Document
```go
func (s *DocumentStore) Create(doc Document) error
```
- Creates a new document
- Returns error if document with same ID already exists
- Thread-safe with exclusive lock

#### Get Document
```go
func (s *DocumentStore) Get(id string) (Document, error)
```
- Retrieves document by ID
- Returns error if document not found
- Thread-safe with shared read lock

#### Delete Document
```go
func (s *DocumentStore) Delete(id string) error
```
- Removes document by ID
- Returns error if document not found
- Thread-safe with exclusive lock

#### List All Documents
```go
func (s *DocumentStore) List() []Document
```
- Returns all documents as a slice
- Thread-safe with shared read lock
- Creates a copy to prevent external modifications

## Usage Example

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Create new document store
    store := NewDocumentStore()
    
    // Create a document
    doc := Document{
        ID:          "1",
        Name:        "My Document",
        Description: "A sample document",
    }
    
    if err := store.Create(doc); err != nil {
        log.Fatal(err)
    }
    
    // Retrieve the document
    retrieved, err := store.Get("1")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Retrieved: %+v\n", retrieved)
    
    // List all documents
    docs := store.List()
    fmt.Printf("Total documents: %d\n", len(docs))
    
    // Delete the document
    if err := store.Delete("1"); err != nil {
        log.Fatal(err)
    }
}
```

## Testing

The project includes comprehensive tests covering:

- **Unit Tests**: All CRUD operations
- **Error Cases**: Invalid operations and edge cases
- **Concurrency Tests**: Thread safety verification


Run specific test types:
```bash
# Run all tests
go test -v

# Run only concurrent tests
go test -v -run TestDocumentStore_Concurrent

# Run benchmarks
go test -bench=. -benchmem

# Test coverage
go test -cover
```

## Thread Safety

The DocumentStore uses `sync.RWMutex` to ensure thread safety:

- **Read Operations** (Get, List): Use `RLock()` allowing multiple concurrent readers
- **Write Operations** (Create, Delete): Use `Lock()` for exclusive access
- **Automatic Cleanup**: `defer` statements ensure locks are always released

## Performance Characteristics

- **Create**: O(1) average case
- **Get**: O(1) lookup time
- **Delete**: O(1) removal time
- **List**: O(n) where n is number of documents

## Docker Support

Build and run with Docker:
```bash
# Build image
docker build -t docstore-api .

# Run container
docker run docstore-api
```

