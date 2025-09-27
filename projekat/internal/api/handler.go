package api

import (
	"encoding/json"
	"log"
	"net/http"

	"mini-zanzibar/internal/storage"
	"mini-zanzibar/pkg/models"
)

// Handler handles HTTP requests for the Zanzibar API
type Handler struct {
	store storage.Store
}

// NewHandler creates a new API handler
func NewHandler(store storage.Store) *Handler {
	return &Handler{store: store}
}

// CheckPermission handles permission check requests
func (h *Handler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding check request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Checking permission: %+v", req)

	allowed, err := h.store.CheckRelation(req)
	if err != nil {
		log.Printf("Error checking relation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.CheckResponse{Allowed: allowed}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Permission check result: %t", allowed)
}

// WriteRelation handles requests to write relations
func (h *Handler) WriteRelation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Try to decode as a single relation first
	var singleRelation models.Relation
	var relations []models.Relation

	// Read the body into a buffer so we can try both formats
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	// Try single relation format first
	if err := json.Unmarshal(body, &singleRelation); err == nil {
		// If it's a valid single relation, use it
		if singleRelation.Object.Type != "" && singleRelation.Object.ID != "" {
			relations = []models.Relation{singleRelation}
			log.Printf("Decoded single relation: %+v", singleRelation)
		}
	}

	// If single relation didn't work, try the array format
	if len(relations) == 0 {
		var req models.WriteRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("Error decoding write request as both single relation and array: %v", err)
			http.Error(w, "Invalid request body - expected single relation or {\"relations\": [...]}", http.StatusBadRequest)
			return
		}
		relations = req.Relations
		log.Printf("Decoded %d relations from array format", len(relations))
	}

	log.Printf("Writing %d relations", len(relations))

	for _, relation := range relations {
		if err := h.store.WriteRelation(relation); err != nil {
			log.Printf("Error writing relation: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Printf("Written relation: %+v", relation)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "success"}`))
}

// ReadRelations handles requests to read relations for an object
func (h *Handler) ReadRelations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var obj models.Object
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		log.Printf("Error decoding read request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	relations, err := h.store.ReadRelations(obj)
	if err != nil {
		log.Printf("Error reading relations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relations); err != nil {
		log.Printf("Error encoding relations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Read %d relations for object %+v", len(relations), obj)
}
