package main

import (
	"fmt"
	"sync"
	"testing"
	
	"docstore-api/models"
)

// Test helper function to create a sample document
func createSampleDocument(id, name, description string) models.Document {
	return models.Document{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

// TestNewDocumentStore tests the constructor
func TestNewDocumentStore(t *testing.T) {
	store := models.NewDocumentStore()
	
	if store == nil {
		t.Fatal("NewDocumentStore() returned nil")
	}
	
	// Test that the store is properly initialized by checking it's empty
	docs := store.List()
	if len(docs) != 0 {
		t.Errorf("expected empty store, got %d documents", len(docs))
	}
	
	// Test that we can perform operations on the new store
	doc := createSampleDocument("test", "Test", "Test document")
	err := store.Create(doc)
	if err != nil {
		t.Errorf("failed to create document in new store: %v", err)
	}
}

// TestDocumentStore_Create tests document creation
func TestDocumentStore_Create(t *testing.T) {
	store := models.NewDocumentStore()
	doc := createSampleDocument("1", "Test Doc", "A test document")
	
	// Test successful creation
	err := store.Create(doc)
	if err != nil {
		t.Errorf("Create() failed: %v", err)
	}
	
	// Verify document was stored using public methods
	docs := store.List()
	if len(docs) != 1 {
		t.Errorf("expected 1 document, got %d", len(docs))
	}
	
	storedDoc, err := store.Get("1")
	if err != nil {
		t.Errorf("document was not stored: %v", err)
	}
	
	if storedDoc != doc {
		t.Errorf("stored document doesn't match original. Got %+v, want %+v", storedDoc, doc)
	}
}

// TestDocumentStore_Create_Duplicate tests duplicate document creation
func TestDocumentStore_Create_Duplicate(t *testing.T) {
	store := models.NewDocumentStore()
	doc := createSampleDocument("1", "Test Doc", "A test document")
	
	// Create first document
	err := store.Create(doc)
	if err != nil {
		t.Fatalf("First Create() failed: %v", err)
	}
	
	// Try to create duplicate
	err = store.Create(doc)
	if err == nil {
		t.Error("Create() should have failed for duplicate document")
	}
	
	expectedError := "document already exists"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
	
	// Verify only one document exists
	docs := store.List()
	if len(docs) != 1 {
		t.Errorf("expected 1 document after duplicate attempt, got %d", len(docs))
	}
}

// TestDocumentStore_Get tests document retrieval
func TestDocumentStore_Get(t *testing.T) {
	store := models.NewDocumentStore()
	doc := createSampleDocument("1", "Test Doc", "A test document")
	
	// Store a document first
	store.Create(doc)
	
	// Test successful retrieval
	retrievedDoc, err := store.Get("1")
	if err != nil {
		t.Errorf("Get() failed: %v", err)
	}
	
	if retrievedDoc != doc {
		t.Errorf("retrieved document doesn't match original. Got %+v, want %+v", retrievedDoc, doc)
	}
}

// TestDocumentStore_Get_NotFound tests retrieval of non-existent document
func TestDocumentStore_Get_NotFound(t *testing.T) {
	store := models.NewDocumentStore()
	
	// Try to get non-existent document
	doc, err := store.Get("nonexistent")
	if err == nil {
		t.Error("Get() should have failed for non-existent document")
	}
	
	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
	
	// Verify empty document is returned
	emptyDoc := models.Document{}
	if doc != emptyDoc {
		t.Errorf("expected empty document, got %+v", doc)
	}
}

// TestDocumentStore_Delete tests document deletion
func TestDocumentStore_Delete(t *testing.T) {
	store := models.NewDocumentStore()
	doc := createSampleDocument("1", "Test Doc", "A test document")
	
	// Store a document first
	store.Create(doc)
	
	// Test successful deletion
	err := store.Delete("1")
	if err != nil {
		t.Errorf("Delete() failed: %v", err)
	}
	
	// Verify document was deleted
	docs := store.List()
	if len(docs) != 0 {
		t.Errorf("expected 0 documents after deletion, got %d", len(docs))
	}
	
	// Verify document can't be retrieved
	_, err = store.Get("1")
	if err == nil {
		t.Error("Get() should fail after document deletion")
	}
}

// TestDocumentStore_Delete_NotFound tests deletion of non-existent document
func TestDocumentStore_Delete_NotFound(t *testing.T) {
	store := models.NewDocumentStore()
	
	// Try to delete non-existent document
	err := store.Delete("nonexistent")
	if err == nil {
		t.Error("Delete() should have failed for non-existent document")
	}
	
	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestDocumentStore_List tests listing all documents
func TestDocumentStore_List(t *testing.T) {
	store := models.NewDocumentStore()
	
	// Test empty list
	docs := store.List()
	if len(docs) != 0 {
		t.Errorf("expected empty list, got %d documents", len(docs))
	}
	
	// Add some documents
	doc1 := createSampleDocument("1", "Doc 1", "First document")
	doc2 := createSampleDocument("2", "Doc 2", "Second document")
	doc3 := createSampleDocument("3", "Doc 3", "Third document")
	
	store.Create(doc1)
	store.Create(doc2)
	store.Create(doc3)
	
	// Test list with documents
	docs = store.List()
	if len(docs) != 3 {
		t.Errorf("expected 3 documents, got %d", len(docs))
	}
	
	// Verify all documents are present (order doesn't matter for maps)
	docMap := make(map[string]models.Document)
	for _, doc := range docs {
		docMap[doc.ID] = doc
	}
	
	if docMap["1"] != doc1 {
		t.Errorf("doc1 not found in list or doesn't match")
	}
	if docMap["2"] != doc2 {
		t.Errorf("doc2 not found in list or doesn't match")
	}
	if docMap["3"] != doc3 {
		t.Errorf("doc3 not found in list or doesn't match")
	}
}
// Test DocumentStore_ConcurrentOperations tests thread safety
func TestDocumentStore_ConcurrentOperations(t *testing.T) {
	store := models.NewDocumentStore()
	numGoroutines := 10
	numOperations := 100
	
	var wg sync.WaitGroup
	
	// Concurrent creates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				docID := fmt.Sprintf("%d-%d", id, j)
				doc := createSampleDocument(docID, fmt.Sprintf("Doc %s", docID), "Concurrent test")
				store.Create(doc)
			}
		}(i)
	}
	wg.Wait()
	
	// Verify all documents were created
	docs := store.List()
	expectedCount := numGoroutines * numOperations
	if len(docs) != expectedCount {
		t.Errorf("expected %d documents, got %d", expectedCount, len(docs))
	}
	
	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				docID := fmt.Sprintf("%d-%d", id, j)
				doc, err := store.Get(docID)
				if err != nil {
					t.Errorf("Get() failed for %s: %v", docID, err)
					return
				}
				if doc.ID != docID {
					t.Errorf("wrong document retrieved: got %s, want %s", doc.ID, docID)
				}
			}
		}(i)
	}
	wg.Wait()
	
	// Concurrent deletes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				docID := fmt.Sprintf("%d-%d", id, j)
				err := store.Delete(docID)
				if err != nil {
					t.Errorf("Delete() failed for %s: %v", docID, err)
				}
			}
		}(i)
	}
	wg.Wait()
	
	// Verify all documents were deleted
	docs = store.List()
	if len(docs) != 0 {
		t.Errorf("expected 0 documents after deletion, got %d", len(docs))
	}
}
