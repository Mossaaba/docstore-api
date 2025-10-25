package controllers

import (
	"docstore-api/src/models"
	"docstore-api/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	service services.DocumentService
}

func NewDocumentController(service services.DocumentService) *DocumentController {
	return &DocumentController{
		service: service,
	}
}

// CreateDocument godoc
// @Summary Create a new document
// @Description Create a new document with the provided information
// @Tags documents
// @Accept json
// @Produce json
// @Param document body models.Document true "Document to create"
// @Success 201 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents [post]
func (ctrl *DocumentController) CreateDocument(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.CreateDocument(doc); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// GetDocument godoc
// @Summary Get a document by ID
// @Description Get a document by its ID
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.Document
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents/{id} [get]
func (ctrl *DocumentController) GetDocument(c *gin.Context) {
	id := c.Param("id")

	doc, err := ctrl.service.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments godoc
// @Summary List all documents
// @Description Get a list of all documents
// @Tags documents
// @Accept json
// @Produce json
// @Success 200 {array} models.Document
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents [get]
func (ctrl *DocumentController) ListDocuments(c *gin.Context) {
	docs := ctrl.service.ListDocuments()
	c.JSON(http.StatusOK, docs)
}

// UpdateDocument godoc
// @Summary Update a document (PUT)
// @Description Replace an entire document with new data
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param document body models.Document true "Document data to update"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents/{id} [put]
func (ctrl *DocumentController) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	var doc models.Document

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateDocument(id, doc); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the updated document
	updatedDoc, _ := ctrl.service.GetDocument(id)
	c.JSON(http.StatusOK, updatedDoc)
}

// PartialUpdateDocument godoc
// @Summary Partially update a document (PATCH)
// @Description Update specific fields of a document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents/{id} [patch]
func (ctrl *DocumentController) PartialUpdateDocument(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.PartialUpdateDocument(id, updates); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the updated document
	updatedDoc, _ := ctrl.service.GetDocument(id)
	c.JSON(http.StatusOK, updatedDoc)
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document by its ID
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/documents/{id} [delete]
func (ctrl *DocumentController) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteDocument(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
