package services

import (
	"docstore-api/models"
	"testing"
)

func TestDocumentService_CreateDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "Test Description",
	}

	err := service.CreateDocument(doc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test duplicate creation
	err = service.CreateDocument(doc)
	if err == nil {
		t.Error("Expected error for duplicate document, got nil")
	}
}

func TestDocumentService_GetDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "Test Description",
	}

	// Create document first
	service.CreateDocument(doc)

	// Test getting existing document
	retrieved, err := service.GetDocument("test-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.ID != doc.ID || retrieved.Name != doc.Name || retrieved.Description != doc.Description {
		t.Errorf("Retrieved document doesn't match original. Got %+v, want %+v", retrieved, doc)
	}

	// Test getting non-existent document
	_, err = service.GetDocument("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent document, got nil")
	}
}

func TestDocumentService_ListDocuments(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Test empty list
	docs := service.ListDocuments()
	if len(docs) != 0 {
		t.Errorf("Expected empty list, got %d documents", len(docs))
	}

	// Add some documents
	doc1 := models.Document{ID: "1", Name: "Doc 1", Description: "First doc"}
	doc2 := models.Document{ID: "2", Name: "Doc 2", Description: "Second doc"}

	service.CreateDocument(doc1)
	service.CreateDocument(doc2)

	docs = service.ListDocuments()
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(docs))
	}
}

func TestDocumentService_UpdateDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Create initial document
	doc := models.Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}
	service.CreateDocument(doc)

	// Update document
	updatedDoc := models.Document{
		ID:          "test-1",
		Name:        "Updated Name",
		Description: "Updated Description",
	}

	err := service.UpdateDocument("test-1", updatedDoc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, _ := service.GetDocument("test-1")
	if retrieved.Name != "Updated Name" || retrieved.Description != "Updated Description" {
		t.Errorf("Document not updated correctly. Got %+v", retrieved)
	}

	// Test updating non-existent document
	err = service.UpdateDocument("non-existent", updatedDoc)
	if err == nil {
		t.Error("Expected error for non-existent document, got nil")
	}
}

func TestDocumentService_PartialUpdateDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Create initial document
	doc := models.Document{
		ID:          "test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}
	service.CreateDocument(doc)

	// Partial update - only name
	updates := map[string]interface{}{
		"name": "Updated Name Only",
	}

	err := service.PartialUpdateDocument("test-1", updates)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify partial update
	retrieved, _ := service.GetDocument("test-1")
	if retrieved.Name != "Updated Name Only" {
		t.Errorf("Name not updated. Got %s, want %s", retrieved.Name, "Updated Name Only")
	}
	if retrieved.Description != "Original Description" {
		t.Errorf("Description should remain unchanged. Got %s, want %s", retrieved.Description, "Original Description")
	}

	// Partial update - only description
	updates = map[string]interface{}{
		"description": "Updated Description Only",
	}

	err = service.PartialUpdateDocument("test-1", updates)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify description update
	retrieved, _ = service.GetDocument("test-1")
	if retrieved.Description != "Updated Description Only" {
		t.Errorf("Description not updated. Got %s, want %s", retrieved.Description, "Updated Description Only")
	}
	if retrieved.Name != "Updated Name Only" {
		t.Errorf("Name should remain unchanged. Got %s, want %s", retrieved.Name, "Updated Name Only")
	}

	// Test partial update on non-existent document
	err = service.PartialUpdateDocument("non-existent", updates)
	if err == nil {
		t.Error("Expected error for non-existent document, got nil")
	}
}

func TestDocumentService_DeleteDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "Test Description",
	}

	// Create document first
	service.CreateDocument(doc)

	// Delete document
	err := service.DeleteDocument("test-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = service.GetDocument("test-1")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}

	// Test deleting non-existent document
	err = service.DeleteDocument("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent document, got nil")
	}
}