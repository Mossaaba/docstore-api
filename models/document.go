package models

import (
	"errors"
	"sync"
) 

/*
This is documnet struct : 
*/
type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
DocumentStore that implements a thread-safe document storage system
DocumentStore - A struct type that acts as a container for documents
mu sync.RWMutex - A read-write mutex for thread safety. This allows multiple concurrent readers OR one exclusive writer
documents map[string]Document - A map that stores documents with string keys (likely document IDs) and Document values


Multiple goroutines can read from the store simultaneously
Only one goroutine can write at a time
Reads are blocked during writes, but writes wait for all reads to complete
*/
type DocumentStore struct {
	mu        sync.RWMutex
	documents map[string]Document
}


/*
Create a constructor that returns a pointer to a new DocumentStore :
make(map[string]Document) - Creates an empty map to store documents
&DocumentStore{...} - Creates a new DocumentStore struct and returns its memory address (pointer)
*/
func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: make(map[string]Document),
	}
}


/*
Checks if document ID already exists before creating
Returns descriptive error if duplicate found
Maintains data integrity

Uses sync.RWMutex to prevent race conditions
Lock() ensures exclusive access during write operations
defer Unlock() guarantees lock release even if function panics

*/
func (s *DocumentStore) Create(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.documents[doc.ID]; exists {
		return errors.New("document already exists")
	}
	s.documents[doc.ID] = doc
	return nil
}

/* 

Uses RLock() instead of Lock() - allows multiple simultaneous reads
More efficient than exclusive locks for read operations
Multiple goroutines can read at the same time

Uses map lookup with existence check: doc, exists := s.documents[id]
Prevents panic if key doesn't exist
Returns appropriate error for missing documents
*/
func (s *DocumentStore) Get(id string) (Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, exists := s.documents[id]
	if !exists {
		return Document{}, errors.New("document not found")
	}
	return doc, nil
}

/*
Uses Lock() (not RLock()) because it's a write operation
Blocks all other reads and writes during deletion
Ensures data consistency

Verifies document exists before attempting deletion
Returns meaningful error if document not found
Prevents silent failures
*/

func (s *DocumentStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.documents[id]; !exists {
		return errors.New("document not found")
	}

	delete(s.documents, id)
	return nil
}

/*

Uses RLock() for shared access - multiple List() calls can run simultaneously
Doesn't block other read operations (Get, List)
Only blocks during write operations (Create, Delete)


make([]Document, 0, len(s.documents)) pre-allocates slice capacity
Avoids multiple memory reallocations during append operations
Performance optimization for large document collections

*/
func (s *DocumentStore) List() []Document {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}
	return docs
}