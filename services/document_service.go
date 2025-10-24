package services

import (
	"docstore-api/models"
)

type DocumentService interface {
	CreateDocument(doc models.Document) error
	GetDocument(id string) (models.Document, error)
	ListDocuments() []models.Document
	DeleteDocument(id string) error
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