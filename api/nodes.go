package handler

import (
	"encoding/json"
	"net/http"
)

// Node is a minimal demo node returned by the API.
type Node struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

func init() {
	// The Vercel Go runtime serves each exported handler function as an individual route.
}

func NodesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method Not Allowed"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	nodes := []Node{
		{ID: "node-123", Title: "Welcome to KOS", Body: "This is a demo node served by the Vercel Go API."},
		{ID: "node-456", Title: "Knowledge Graph", Body: "A second node showing the project can scale to real data."},
	}

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
