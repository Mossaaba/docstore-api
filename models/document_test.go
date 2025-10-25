package models

import (
	"sync"
	"testing"
)

func TestNewDocumentStore(t *testing.T) {
	store := NewDocumentStore()

	if store == nil {
		t.Fatal("NewDocumentStore() returned nil")
	}

	if store.documents == nil {
		t.Error("documents map not initialized")
	}

	if len(store.documents) != 0 {
		t.Error("new store should be empty")
	}
}

func TestDocumentStore_Create(t *testing.T) {
	store := NewDocumentStore()

	doc := Document{
		ID:          "1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Test successful creation
	err := store.Create(doc)
	if err != nil {
		t.Errorf("Create() failed: %v", err)
	}

	// Verify document was stored
	if len(store.documents) != 1 {
		t.Errorf("expected 1 document, got %d", len(store.documents))
	}

	stored, exists := store.documents[doc.ID]
	if !exists {
		t.Error("document not found in store")
	}

	if stored != doc {
		t.Errorf("stored document doesn't match: got %+v, want %+v", stored, doc)
	}
}

func TestDocumentStore_CreateDuplicate(t *testing.T) {
	store := NewDocumentStore()

	doc := Document{
		ID:          "1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create first document
	err := store.Create(doc)
	if err != nil {
		t.Fatalf("first Create() failed: %v", err)
	}

	// Try to create duplicate
	err = store.Create(doc)
	if err == nil {
		t.Error("Create() should fail for duplicate ID")
	}

	expectedError := "document already exists"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentStore_Get(t *testing.T) {
	store := NewDocumentStore()

	doc := Document{
		ID:          "1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	store.Create(doc)

	// Test successful get
	retrieved, err := store.Get(doc.ID)
	if err != nil {
		t.Errorf("Get() failed: %v", err)
	}

	if retrieved != doc {
		t.Errorf("retrieved document doesn't match: got %+v, want %+v", retrieved, doc)
	}
}

func TestDocumentStore_GetNotFound(t *testing.T) {
	store := NewDocumentStore()

	// Try to get non-existent document
	_, err := store.Get("non-existent")
	if err == nil {
		t.Error("Get() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentStore_Delete(t *testing.T) {
	store := NewDocumentStore()

	doc := Document{
		ID:          "1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	store.Create(doc)

	// Test successful deletion
	err := store.Delete(doc.ID)
	if err != nil {
		t.Errorf("Delete() failed: %v", err)
	}

	// Verify document was deleted
	if len(store.documents) != 0 {
		t.Errorf("expected 0 documents after deletion, got %d", len(store.documents))
	}

	// Verify document can't be retrieved
	_, err = store.Get(doc.ID)
	if err == nil {
		t.Error("Get() should fail after document deletion")
	}
}

func TestDocumentStore_DeleteNotFound(t *testing.T) {
	store := NewDocumentStore()

	// Try to delete non-existent document
	err := store.Delete("non-existent")
	if err == nil {
		t.Error("Delete() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentStore_List(t *testing.T) {
	store := NewDocumentStore()

	// Test empty list
	docs := store.List()
	if len(docs) != 0 {
		t.Errorf("expected empty list, got %d documents", len(docs))
	}

	// Add some documents
	doc1 := Document{ID: "1", Name: "Doc 1", Description: "First document"}
	doc2 := Document{ID: "2", Name: "Doc 2", Description: "Second document"}

	store.Create(doc1)
	store.Create(doc2)

	// Test list with documents
	docs = store.List()
	if len(docs) != 2 {
		t.Errorf("expected 2 documents, got %d", len(docs))
	}

	// Verify all documents are present (order doesn't matter)
	found := make(map[string]bool)
	for _, doc := range docs {
		found[doc.ID] = true
	}

	if !found["1"] || !found["2"] {
		t.Error("not all documents found in list")
	}
}

func TestDocumentStore_ConcurrentAccess(t *testing.T) {
	store := NewDocumentStore()

	// Test concurrent writes
	var wg sync.WaitGroup
	numGoroutines := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			doc := Document{
				ID:          string(rune('0' + id)),
				Name:        "Concurrent Doc",
				Description: "Created concurrently",
			}
			store.Create(doc)
		}(i)
	}

	wg.Wait()

	// Verify all documents were created
	docs := store.List()
	if len(docs) != numGoroutines {
		t.Errorf("expected %d documents, got %d", numGoroutines, len(docs))
	}
}

func TestDocumentStore_ConcurrentReadWrite(t *testing.T) {
	store := NewDocumentStore()

	// Pre-populate with some documents
	for i := 0; i < 5; i++ {
		doc := Document{
			ID:          string(rune('0' + i)),
			Name:        "Initial Doc",
			Description: "Pre-populated",
		}
		store.Create(doc)
	}

	var wg sync.WaitGroup

	// Start concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				store.List()
				store.Get("0")
			}
		}()
	}

	// Start concurrent writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			doc := Document{
				ID:          "concurrent-" + string(rune('0'+id)),
				Name:        "Concurrent Write",
				Description: "Written during concurrent test",
			}
			store.Create(doc)
		}(i)
	}

	wg.Wait()

	// Verify final state
	docs := store.List()
	if len(docs) < 5 {
		t.Errorf("expected at least 5 documents, got %d", len(docs))
	}
}

func TestDocumentStore_Update(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test successful update
	updatedDoc := Document{
		ID:          "different-id", // This should be overridden by the method
		Name:        "Updated Name",
		Description: "Updated Description",
	}

	err = store.Update("test-1", updatedDoc)
	if err != nil {
		t.Errorf("Update() failed: %v", err)
	}

	// Verify document was updated
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after update: %v", err)
	}

	// Check that ID was preserved (overridden by method)
	if retrieved.ID != "test-1" {
		t.Errorf("expected ID 'test-1', got '%s'", retrieved.ID)
	}

	// Check that other fields were updated
	if retrieved.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", retrieved.Name)
	}

	if retrieved.Description != "Updated Description" {
		t.Errorf("expected description 'Updated Description', got '%s'", retrieved.Description)
	}
}

func TestDocumentStore_UpdateNotFound(t *testing.T) {
	store := NewDocumentStore()

	updatedDoc := Document{
		ID:          "non-existent",
		Name:        "Updated Name",
		Description: "Updated Description",
	}

	// Try to update non-existent document
	err := store.Update("non-existent", updatedDoc)
	if err == nil {
		t.Error("Update() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentStore_PartialUpdate(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test partial update - only name
	updates := map[string]interface{}{
		"name": "Updated Name Only",
	}

	err = store.PartialUpdate("test-1", updates)
	if err != nil {
		t.Errorf("PartialUpdate() failed: %v", err)
	}

	// Verify only name was updated
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after partial update: %v", err)
	}

	if retrieved.Name != "Updated Name Only" {
		t.Errorf("expected name 'Updated Name Only', got '%s'", retrieved.Name)
	}

	if retrieved.Description != "Original Description" {
		t.Errorf("description should remain unchanged, got '%s'", retrieved.Description)
	}

	if retrieved.ID != "test-1" {
		t.Errorf("ID should remain unchanged, got '%s'", retrieved.ID)
	}
}

func TestDocumentStore_PartialUpdateDescription(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test partial update - only description
	updates := map[string]interface{}{
		"description": "Updated Description Only",
	}

	err = store.PartialUpdate("test-1", updates)
	if err != nil {
		t.Errorf("PartialUpdate() failed: %v", err)
	}

	// Verify only description was updated
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after partial update: %v", err)
	}

	if retrieved.Description != "Updated Description Only" {
		t.Errorf("expected description 'Updated Description Only', got '%s'", retrieved.Description)
	}

	if retrieved.Name != "Original Name" {
		t.Errorf("name should remain unchanged, got '%s'", retrieved.Name)
	}

	if retrieved.ID != "test-1" {
		t.Errorf("ID should remain unchanged, got '%s'", retrieved.ID)
	}
}

func TestDocumentStore_PartialUpdateBothFields(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test partial update - both fields
	updates := map[string]interface{}{
		"name":        "Updated Name",
		"description": "Updated Description",
	}

	err = store.PartialUpdate("test-1", updates)
	if err != nil {
		t.Errorf("PartialUpdate() failed: %v", err)
	}

	// Verify both fields were updated
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after partial update: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", retrieved.Name)
	}

	if retrieved.Description != "Updated Description" {
		t.Errorf("expected description 'Updated Description', got '%s'", retrieved.Description)
	}

	if retrieved.ID != "test-1" {
		t.Errorf("ID should remain unchanged, got '%s'", retrieved.ID)
	}
}

func TestDocumentStore_PartialUpdateInvalidTypes(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test partial update with invalid types (should be ignored)
	updates := map[string]interface{}{
		"name":        123,     // Invalid type - should be ignored
		"description": true,    // Invalid type - should be ignored
		"invalid":     "field", // Unknown field - should be ignored
	}

	err = store.PartialUpdate("test-1", updates)
	if err != nil {
		t.Errorf("PartialUpdate() failed: %v", err)
	}

	// Verify document remains unchanged
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after partial update: %v", err)
	}

	if retrieved.Name != "Original Name" {
		t.Errorf("name should remain unchanged, got '%s'", retrieved.Name)
	}

	if retrieved.Description != "Original Description" {
		t.Errorf("description should remain unchanged, got '%s'", retrieved.Description)
	}
}

func TestDocumentStore_PartialUpdateNotFound(t *testing.T) {
	store := NewDocumentStore()

	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	// Try to partial update non-existent document
	err := store.PartialUpdate("non-existent", updates)
	if err == nil {
		t.Error("PartialUpdate() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentStore_PartialUpdateEmptyUpdates(t *testing.T) {
	store := NewDocumentStore()

	// Create initial document
	originalDoc := Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}

	err := store.Create(originalDoc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test partial update with empty updates map
	updates := map[string]interface{}{}

	err = store.PartialUpdate("test-1", updates)
	if err != nil {
		t.Errorf("PartialUpdate() failed: %v", err)
	}

	// Verify document remains unchanged
	retrieved, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get() failed after partial update: %v", err)
	}

	if retrieved != originalDoc {
		t.Errorf("document should remain unchanged, got %+v, want %+v", retrieved, originalDoc)
	}
}
