package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-zanzibar/internal/api"
	"mini-zanzibar/internal/storage"
	"mini-zanzibar/pkg/models"
)

func TestBasicPermissionCheck(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := api.NewHandler(store)

	// Write a relation
	writeReq := models.WriteRequest{
		Relations: []models.Relation{
			{
				Object:   models.Object{Type: "document", ID: "doc1"},
				Relation: "viewer",
				Subject:  models.Subject{Type: "user", ID: "alice"},
			},
		},
	}

	body, _ := json.Marshal(writeReq)
	req := httptest.NewRequest("POST", "/write", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.WriteRelation(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Check permission (should be allowed)
	checkReq := models.CheckRequest{
		Object:   models.Object{Type: "document", ID: "doc1"},
		Relation: "viewer",
		Subject:  models.Subject{Type: "user", ID: "alice"},
	}

	body, _ = json.Marshal(checkReq)
	req = httptest.NewRequest("POST", "/check", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.CheckPermission(w, req)

	var response models.CheckResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Allowed {
		t.Error("Expected permission to be allowed")
	}

	// Check permission for different user (should be denied)
	checkReq.Subject.ID = "bob"

	body, _ = json.Marshal(checkReq)
	req = httptest.NewRequest("POST", "/check", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.CheckPermission(w, req)

	json.NewDecoder(w.Body).Decode(&response)

	if response.Allowed {
		t.Error("Expected permission to be denied for bob")
	}
}

func TestMultipleRelations(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := api.NewHandler(store)

	// Write multiple relations
	writeReq := models.WriteRequest{
		Relations: []models.Relation{
			{
				Object:   models.Object{Type: "document", ID: "doc1"},
				Relation: "owner",
				Subject:  models.Subject{Type: "user", ID: "alice"},
			},
			{
				Object:   models.Object{Type: "document", ID: "doc1"},
				Relation: "viewer",
				Subject:  models.Subject{Type: "user", ID: "bob"},
			},
		},
	}

	body, _ := json.Marshal(writeReq)
	req := httptest.NewRequest("POST", "/write", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.WriteRelation(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Check owner permission for alice
	checkReq := models.CheckRequest{
		Object:   models.Object{Type: "document", ID: "doc1"},
		Relation: "owner",
		Subject:  models.Subject{Type: "user", ID: "alice"},
	}

	body, _ = json.Marshal(checkReq)
	req = httptest.NewRequest("POST", "/check", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.CheckPermission(w, req)

	var response models.CheckResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Allowed {
		t.Error("Expected alice to be owner")
	}

	// Check viewer permission for bob
	checkReq.Relation = "viewer"
	checkReq.Subject.ID = "bob"

	body, _ = json.Marshal(checkReq)
	req = httptest.NewRequest("POST", "/check", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	handler.CheckPermission(w, req)

	json.NewDecoder(w.Body).Decode(&response)

	if !response.Allowed {
		t.Error("Expected bob to be viewer")
	}
}
