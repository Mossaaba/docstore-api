package models

import (
	"errors"
	"reflect"
	"sync"
)

type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DocumentStore struct {
	mu        sync.RWMutex
	documents map[string]Document
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: make(map[string]Document),
	}
}

func (s *DocumentStore) Create(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.documents[doc.ID]; exists {
		return errors.New("document already exists")
	}
	s.documents[doc.ID] = doc
	return nil
}

func (s *DocumentStore) Get(id string) (Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, exists := s.documents[id]
	if !exists {
		return Document{}, errors.New("document not found")
	}
	return doc, nil
}

func (s *DocumentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.documents[id]; !exists {
		return errors.New("document not found")
	}

	delete(s.documents, id)
	return nil
}

func (s *DocumentStore) List() []Document {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}
	return docs
}

func (s *DocumentStore) Update(id string, doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.documents[id]; !exists {
		return errors.New("document not found")
	}

	// Ensure the document ID matches the path parameter
	doc.ID = id
	s.documents[id] = doc
	return nil
}

func (s *DocumentStore) PartialUpdate(id string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, exists := s.documents[id]
	if !exists {
		return errors.New("document not found")
	}

	// Use reflection to automatically detect and update attributes
	docValue := reflect.ValueOf(&doc).Elem()
	docType := reflect.TypeOf(doc)

	for key, value := range updates {
		// Skip ID field to prevent modification
		if key == "id" {
			continue
		}

		// Find the field by JSON tag or field name
		fieldIndex := -1
		for i := 0; i < docType.NumField(); i++ {
			field := docType.Field(i)
			jsonTag := field.Tag.Get("json")
			
			// Check if the key matches the JSON tag or field name
			if jsonTag == key || field.Name == key {
				fieldIndex = i
				break
			}
		}

		// If field found, update it only if types match exactly
		if fieldIndex >= 0 {
			field := docValue.Field(fieldIndex)
			if field.CanSet() {
				valueReflect := reflect.ValueOf(value)
				
				// Only update if the types match exactly (no conversion)
				if valueReflect.Type() == field.Type() {
					field.Set(valueReflect)
				}
				// Invalid types are silently ignored
			}
		}
		// Unknown fields are silently ignored
	}

	s.documents[id] = doc
	return nil
}
