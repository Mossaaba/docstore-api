package services

import (
	"docstore-api/src/models"
)

type DocumentService interface {
	CreateDocument(doc models.Document) error
	GetDocument(id string) (models.Document, error)
	ListDocuments() []models.Document
	DeleteDocument(id string) error
	UpdateDocument(id string, doc models.Document) error
	PartialUpdateDocument(id string, updates map[string]interface{}) error
}

type documentService struct {
	store *models.DocumentStore
}

func NewDocumentService(store *models.DocumentStore) DocumentService {
	return &documentService{
		store: store,
	}
}

func (s *documentService) CreateDocument(doc models.Document) error {
	return s.store.Create(doc)
}

func (s *documentService) GetDocument(id string) (models.Document, error) {
	return s.store.Get(id)
}

func (s *documentService) ListDocuments() []models.Document {
	return s.store.List()
}

func (s *documentService) DeleteDocument(id string) error {
	return s.store.Delete(id)
}

func (s *documentService) UpdateDocument(id string, doc models.Document) error {
	return s.store.Update(id, doc)
}

func (s *documentService) PartialUpdateDocument(id string, updates map[string]interface{}) error {
	return s.store.PartialUpdate(id, updates)
}
