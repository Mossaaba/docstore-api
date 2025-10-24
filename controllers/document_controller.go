package controllers

import (
	"net/http"
	"docstore-api/services"
	"docstore-api/models"
	
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
// @Failure 409 {object} map[string]string
// @Router /documents [post]
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
// @Failure 404 {object} map[string]string
// @Router /documents/{id} [get]
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
// @Router /documents [get]
func (ctrl *DocumentController) ListDocuments(c *gin.Context) {
	docs := ctrl.service.ListDocuments()
	c.JSON(http.StatusOK, docs)
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document by its ID
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /documents/{id} [delete]
func (ctrl *DocumentController) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	
	if err := ctrl.service.DeleteDocument(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}