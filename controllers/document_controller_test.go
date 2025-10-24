package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"docstore-api/models"
	"docstore-api/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *DocumentController) {
	gin.SetMode(gin.TestMode)

	store := models.NewDocumentStore()
	service := services.NewDocumentService(store)
	controller := NewDocumentController(service)

	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/documents", controller.CreateDocument)
		v1.GET("/documents", controller.ListDocuments)
		v1.GET("/documents/:id", controller.GetDocument)
		v1.DELETE("/documents/:id", controller.DeleteDocument)
	}

	return router, controller
}

func TestNewDocumentController(t *testing.T) {
	store := models.NewDocumentStore()
	service := services.NewDocumentService(store)
	controller := NewDocumentController(service)

	assert.NotNil(t, controller)
	assert.NotNil(t, controller.service)
}

func TestDocumentController_CreateDocument(t *testing.T) {
	router, _ := setupTestRouter()

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	jsonData, _ := json.Marshal(doc)

	req, _ := http.NewRequest("POST", "/api/v1/documents", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Document
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, doc, response)
}

func TestDocumentController_CreateDocumentInvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	invalidJSON := `{"id": "test-1", "name": "Test", "description":}`

	req, _ := http.NewRequest("POST", "/api/v1/documents", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid")
}

func TestDocumentController_CreateDocumentDuplicate(t *testing.T) {
	router, _ := setupTestRouter()

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	jsonData, _ := json.Marshal(doc)

	// Create first document
	req1, _ := http.NewRequest("POST", "/api/v1/documents", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Try to create duplicate
	req2, _ := http.NewRequest("POST", "/api/v1/documents", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)

	var response map[string]string
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "document already exists", response["error"])
}

func TestDocumentController_GetDocument(t *testing.T) {
	router, controller := setupTestRouter()

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	controller.service.CreateDocument(doc)

	req, _ := http.NewRequest("GET", "/api/v1/documents/test-1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Document
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, doc, response)
}

func TestDocumentController_GetDocumentNotFound(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/documents/non-existent", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "document not found", response["error"])
}

func TestDocumentController_ListDocuments(t *testing.T) {
	router, controller := setupTestRouter()

	// Test empty list
	req, _ := http.NewRequest("GET", "/api/v1/documents", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Document
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response))

	// Add some documents
	docs := []models.Document{
		{ID: "1", Name: "Doc 1", Description: "First document"},
		{ID: "2", Name: "Doc 2", Description: "Second document"},
	}

	for _, doc := range docs {
		controller.service.CreateDocument(doc)
	}

	// Test list with documents
	req2, _ := http.NewRequest("GET", "/api/v1/documents", nil)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response2 []models.Document
	err = json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response2))

	// Verify all documents are present (order doesn't matter)
	found := make(map[string]bool)
	for _, doc := range response2 {
		found[doc.ID] = true
	}

	assert.True(t, found["1"])
	assert.True(t, found["2"])
}

func TestDocumentController_DeleteDocument(t *testing.T) {
	router, controller := setupTestRouter()

	doc := models.Document{
		ID:          "test-1",
		Name:        "Test Document",
		Description: "A test document",
	}

	// Create document first
	controller.service.CreateDocument(doc)

	req, _ := http.NewRequest("DELETE", "/api/v1/documents/test-1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	// Verify document was deleted by trying to get it
	req2, _ := http.NewRequest("GET", "/api/v1/documents/test-1", nil)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestDocumentController_DeleteDocumentNotFound(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "/api/v1/documents/non-existent", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "document not found", response["error"])
}

func TestDocumentController_FullAPIWorkflow(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Create a document
	doc := models.Document{
		ID:          "workflow-1",
		Name:        "Workflow Document",
		Description: "Testing full workflow",
	}

	jsonData, _ := json.Marshal(doc)

	req1, _ := http.NewRequest("POST", "/api/v1/documents", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Get the document
	req2, _ := http.NewRequest("GET", "/api/v1/documents/workflow-1", nil)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var retrievedDoc models.Document
	err := json.Unmarshal(w2.Body.Bytes(), &retrievedDoc)
	assert.NoError(t, err)
	assert.Equal(t, doc, retrievedDoc)

	// 3. List documents (should contain our document)
	req3, _ := http.NewRequest("GET", "/api/v1/documents", nil)

	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var docs []models.Document
	err = json.Unmarshal(w3.Body.Bytes(), &docs)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(docs))
	assert.Equal(t, doc, docs[0])

	// 4. Delete the document
	req4, _ := http.NewRequest("DELETE", "/api/v1/documents/workflow-1", nil)

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusNoContent, w4.Code)

	// 5. Verify document is gone
	req5, _ := http.NewRequest("GET", "/api/v1/documents/workflow-1", nil)

	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusNotFound, w5.Code)

	// 6. List should be empty
	req6, _ := http.NewRequest("GET", "/api/v1/documents", nil)

	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)
	assert.Equal(t, http.StatusOK, w6.Code)

	var finalDocs []models.Document
	err = json.Unmarshal(w6.Body.Bytes(), &finalDocs)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(finalDocs))
}
