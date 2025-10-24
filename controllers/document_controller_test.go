package controllers

import (
	"bytes"
	"docstore-api/models"
	"docstore-api/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *DocumentController) {
	gin.SetMode(gin.TestMode)
	store := models.NewDocumentStore()
	service := services.NewDocumentService(store)
	controller := NewDocumentController(service)
	router := gin.New()
	return router, controller
}
func TestDocumentController_CreateDocument(t *testing.T) {
	router, controller := setupTestRouter()
	router.POST("/documents", controller.CreateDocument)
	t.Run("Valid document creation", func(t *testing.T) {
		doc := models.Document{
			ID:          "test-1",
			Name:        "Test Document",
			Description: "Test Description",
		}
		jsonData, _ := json.Marshal(doc)
		req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, response.ID)
		assert.Equal(t, doc.Name, response.Name)
		assert.Equal(t, doc.Description, response.Description)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Duplicate document", func(t *testing.T) {
		doc := models.Document{
			ID:          "duplicate-test",
			Name:        "Duplicate Test",
			Description: "Test Description",
		}
		
		jsonData, _ := json.Marshal(doc)
		
		// Create first document
		req1, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)
		
		// Try to create duplicate
		req2, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}
func TestDocumentController_GetDocument(t *testing.T) {
	router, controller := setupTestRouter()
	router.GET("/documents/:id", controller.GetDocument)
	router.POST("/documents", controller.CreateDocument)

	// Create a test document first
	doc := models.Document{
		ID:          "get-test-1",
		Name:        "Get Test Document",
		Description: "Test Description",
	}
	jsonData, _ := json.Marshal(doc)
	createReq, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	t.Run("Get existing document", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents/get-test-1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, response.ID)
		assert.Equal(t, doc.Name, response.Name)
	})

	t.Run("Get non-existent document", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents/non-existent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
func TestDocumentController_ListDocuments(t *testing.T) {
	router, controller := setupTestRouter()
	router.GET("/documents", controller.ListDocuments)
	router.POST("/documents", controller.CreateDocument)

	t.Run("Empty list", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response))
	})

	t.Run("List with documents", func(t *testing.T) {
		// Create test documents
		docs := []models.Document{
			{ID: "list-1", Name: "Doc 1", Description: "First doc"},
			{ID: "list-2", Name: "Doc 2", Description: "Second doc"},
		}
		
		for _, doc := range docs {
			jsonData, _ := json.Marshal(doc)
			createReq, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
			createReq.Header.Set("Content-Type", "application/json")
			createW := httptest.NewRecorder()
			router.ServeHTTP(createW, createReq)
		}
		
		req, _ := http.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
	})
}
func TestDocumentController_UpdateDocument(t *testing.T) {
	router, controller := setupTestRouter()
	router.PUT("/documents/:id", controller.UpdateDocument)
	router.POST("/documents", controller.CreateDocument)

	// Create a test document first
	originalDoc := models.Document{
		ID:          "update-test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}
	jsonData, _ := json.Marshal(originalDoc)
	createReq, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	t.Run("Valid update", func(t *testing.T) {
		updatedDoc := models.Document{
			ID:          "update-test-1",
			Name:        "Updated Name",
			Description: "Updated Description",
		}
		
		jsonData, _ := json.Marshal(updatedDoc)
		req, _ := http.NewRequest("PUT", "/documents/update-test-1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", response.Name)
		assert.Equal(t, "Updated Description", response.Description)
	})

	t.Run("Update non-existent document", func(t *testing.T) {
		doc := models.Document{
			ID:          "non-existent",
			Name:        "Test",
			Description: "Test",
		}
		
		jsonData, _ := json.Marshal(doc)
		req, _ := http.NewRequest("PUT", "/documents/non-existent", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/documents/update-test-1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
func TestDocumentController_PartialUpdateDocument(t *testing.T) {
	router, controller := setupTestRouter()
	router.PATCH("/documents/:id", controller.PartialUpdateDocument)
	router.POST("/documents", controller.CreateDocument)

	// Create a test document first
	originalDoc := models.Document{
		ID:          "patch-test-1",
		Name:        "Original Name",
		Description: "Original Description",
	}
	jsonData, _ := json.Marshal(originalDoc)
	createReq, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	t.Run("Partial update - name only", func(t *testing.T) {
		updates := map[string]interface{}{
			"name": "Updated Name Only",
		}
		
		jsonData, _ := json.Marshal(updates)
		req, _ := http.NewRequest("PATCH", "/documents/patch-test-1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name Only", response.Name)
		assert.Equal(t, "Original Description", response.Description) // Should remain unchanged
	})

	t.Run("Partial update - description only", func(t *testing.T) {
		updates := map[string]interface{}{
			"description": "Updated Description Only",
		}
		
		jsonData, _ := json.Marshal(updates)
		req, _ := http.NewRequest("PATCH", "/documents/patch-test-1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name Only", response.Name) // Should remain from previous test
		assert.Equal(t, "Updated Description Only", response.Description)
	})

	t.Run("Partial update - both fields", func(t *testing.T) {
		updates := map[string]interface{}{
			"name":        "Both Updated Name",
			"description": "Both Updated Description",
		}
		
		jsonData, _ := json.Marshal(updates)
		req, _ := http.NewRequest("PATCH", "/documents/patch-test-1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.Document
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Both Updated Name", response.Name)
		assert.Equal(t, "Both Updated Description", response.Description)
	})

	t.Run("Partial update non-existent document", func(t *testing.T) {
		updates := map[string]interface{}{
			"name": "Test",
		}
		
		jsonData, _ := json.Marshal(updates)
		req, _ := http.NewRequest("PATCH", "/documents/non-existent", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/documents/patch-test-1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
func TestDocumentController_DeleteDocument(t *testing.T) {
	router, controller := setupTestRouter()
	router.DELETE("/documents/:id", controller.DeleteDocument)
	router.POST("/documents", controller.CreateDocument)
	router.GET("/documents/:id", controller.GetDocument)

	// Create a test document first
	doc := models.Document{
		ID:          "delete-test-1",
		Name:        "Delete Test Document",
		Description: "Test Description",
	}
	jsonData, _ := json.Marshal(doc)
	createReq, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	t.Run("Delete existing document", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/documents/delete-test-1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNoContent, w.Code)
		
		// Verify document is deleted by trying to get it
		getReq, _ := http.NewRequest("GET", "/documents/delete-test-1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		assert.Equal(t, http.StatusNotFound, getW.Code)
	})

	t.Run("Delete non-existent document", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/documents/non-existent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}