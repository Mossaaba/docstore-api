package services

import (
	"docstore-api/models"
	"testing"
)

func TestNewDocumentService(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	if service == nil {
		t.Fatal("NewDocumentService() returned nil")
	}
}

func TestDocumentService_CreateDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Test successful creation
	err := service.CreateDocument(doc)
	if err != nil {
		t.Errorf("CreateDocument() failed: %v", err)
	}

	// Verify document was created by retrieving it
	retrieved, err := service.GetDocument(doc.ID)
	if err != nil {
		t.Errorf("GetDocument() failed after creation: %v", err)
	}

	if retrieved != doc {
		t.Errorf("retrieved document doesn't match: got %+v, want %+v", retrieved, doc)
	}
}

func TestDocumentService_CreateDocumentDuplicate(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create first document
	err := service.CreateDocument(doc)
	if err != nil {
		t.Fatalf("first CreateDocument() failed: %v", err)
	}

	// Try to create duplicate
	err = service.CreateDocument(doc)
	if err == nil {
		t.Error("CreateDocument() should fail for duplicate ID")
	}

	expectedError := "document already exists"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentService_GetDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	service.CreateDocument(doc)

	// Test successful get
	retrieved, err := service.GetDocument(doc.ID)
	if err != nil {
		t.Errorf("GetDocument() failed: %v", err)
	}

	if retrieved != doc {
		t.Errorf("retrieved document doesn't match: got %+v, want %+v", retrieved, doc)
	}
}

func TestDocumentService_GetDocumentNotFound(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Try to get non-existent document
	_, err := service.GetDocument("non-existent")
	if err == nil {
		t.Error("GetDocument() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentService_ListDocuments(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Test empty list
	docs := service.ListDocuments()
	if len(docs) != 0 {
		t.Errorf("expected empty list, got %d documents", len(docs))
	}

	// Add some documents
	doc1 := models.Document{ID: "1", Name: "Doc 1", Description: "First document"}
	doc2 := models.Document{ID: "2", Name: "Doc 2", Description: "Second document"}

	service.CreateDocument(doc1)
	service.CreateDocument(doc2)

	// Test list with documents
	docs = service.ListDocuments()
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

func TestDocumentService_DeleteDocument(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	service.CreateDocument(doc)

	// Test successful deletion
	err := service.DeleteDocument(doc.ID)
	if err != nil {
		t.Errorf("DeleteDocument() failed: %v", err)
	}

	// Verify document was deleted
	_, err = service.GetDocument(doc.ID)
	if err == nil {
		t.Error("GetDocument() should fail after document deletion")
	}

	// Verify list is empty
	docs := service.ListDocuments()
	if len(docs) != 0 {
		t.Errorf("expected 0 documents after deletion, got %d", len(docs))
	}
}

func TestDocumentService_DeleteDocumentNotFound(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Try to delete non-existent document
	err := service.DeleteDocument("non-existent")
	if err == nil {
		t.Error("DeleteDocument() should fail for non-existent document")
	}

	expectedError := "document not found"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDocumentService_FullWorkflow(t *testing.T) {
	store := models.NewDocumentStore()
	service := NewDocumentService(store)

	// Create multiple documents
	docs := []models.Document{
		{ID: "1", Name: "Doc 1", Description: "First document"},
		{ID: "2", Name: "Doc 2", Description: "Second document"},
		{ID: "3", Name: "Doc 3", Description: "Third document"},
	}

	// Create all documents
	for _, doc := range docs {
		err := service.CreateDocument(doc)
		if err != nil {
			t.Errorf("CreateDocument() failed for %s: %v", doc.ID, err)
		}
	}

	// Verify all documents exist
	allDocs := service.ListDocuments()
	if len(allDocs) != 3 {
		t.Errorf("expected 3 documents, got %d", len(allDocs))
	}

	// Get each document individually
	for _, doc := range docs {
		retrieved, err := service.GetDocument(doc.ID)
		if err != nil {
			t.Errorf("GetDocument() failed for %s: %v", doc.ID, err)
		}
		if retrieved != doc {
			t.Errorf("retrieved document doesn't match for %s: got %+v, want %+v", doc.ID, retrieved, doc)
		}
	}

	// Delete one document
	err := service.DeleteDocument("2")
	if err != nil {
		t.Errorf("DeleteDocument() failed: %v", err)
	}

	// Verify document was deleted
	_, err = service.GetDocument("2")
	if err == nil {
		t.Error("GetDocument() should fail for deleted document")
	}

	// Verify remaining documents
	remainingDocs := service.ListDocuments()
	if len(remainingDocs) != 2 {
		t.Errorf("expected 2 documents after deletion, got %d", len(remainingDocs))
	}

	// Verify correct documents remain
	found := make(map[string]bool)
	for _, doc := range remainingDocs {
		found[doc.ID] = true
	}

	if !found["1"] || !found["3"] || found["2"] {
		t.Error("incorrect documents remaining after deletion")
	}
}
