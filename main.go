package main

import (
	"fmt"
	"log"
	
	"docstore-api/models"
)

func main() {
	// Create a new document store
	store := models.NewDocumentStore()
	
	// Create some sample documents
	doc1 := models.Document{
		ID:          "1",
		Name:        "Getting Started",
		Description: "A guide to getting started with the document store",
	}
	
	doc2 := models.Document{
		ID:          "2", 
		Name:        "API Reference",
		Description: "Complete API reference for document operations",
	}
	
	// Add documents to store
	if err := store.Create(doc1); err != nil {
		log.Printf("Error creating doc1: %v", err)
	} else {
		fmt.Println("Created document:", doc1.Name)
	}
	
	if err := store.Create(doc2); err != nil {
		log.Printf("Error creating doc2: %v", err)
	} else {
		fmt.Println("Created document:", doc2.Name)
	}
	
	// List all documents
	fmt.Println("\nAll documents:")
	docs := store.List()
	for _, doc := range docs {
		fmt.Printf("- %s: %s\n", doc.Name, doc.Description)
	}
	
	// Get a specific document
	fmt.Println("\nRetrieving document with ID '1':")
	if doc, err := store.Get("1"); err != nil {
		log.Printf("Error getting document: %v", err)
	} else {
		fmt.Printf("Found: %s - %s\n", doc.Name, doc.Description)
	}
	
	// Try to get non-existent document
	fmt.Println("\nTrying to get non-existent document:")
	if _, err := store.Get("999"); err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
	
	// Delete a document
	fmt.Println("\nDeleting document with ID '1':")
	if err := store.Delete("1"); err != nil {
		log.Printf("Error deleting document: %v", err)
	} else {
		fmt.Println("Document deleted successfully")
	}
	
	// List documents after deletion
	fmt.Println("\nDocuments after deletion:")
	docs = store.List()
	for _, doc := range docs {
		fmt.Printf("- %s: %s\n", doc.Name, doc.Description)
	}
}