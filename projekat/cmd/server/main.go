package main

import (
	"log"
	"net/http"

	"mini-zanzibar/internal/api"
	"mini-zanzibar/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize storage
	store := storage.NewMemoryStore()

	// Initialize handlers
	handler := api.NewHandler(store)

	// Setup routes with gorilla/mux for better routing
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/check", handler.CheckPermission).Methods("POST")
	r.HandleFunc("/write", handler.WriteRelation).Methods("POST")
	r.HandleFunc("/read", handler.ReadRelations).Methods("POST")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	// Add CORS middleware for development
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	log.Println("Mini Zanzibar server starting on :8080")
	log.Println("Available endpoints:")
	log.Println("  POST /check - Check permissions")
	log.Println("  POST /write - Write relations")
	log.Println("  POST /read - Read relations")
	log.Println("  GET /health - Health check")

	log.Fatal(http.ListenAndServe(":8080", r))
}
