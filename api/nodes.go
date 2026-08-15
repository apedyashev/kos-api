package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"kos-api/internal/config"
	"kos-api/internal/domain"
	"kos-api/internal/repository/postgres"
	"kos-api/internal/service"
)

// NodeResponseDTO is the public representation of a node.
type NodeResponseDTO struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Content   json.RawMessage `json:"content"`
	IsStub    bool            `json:"isStub"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type createNodeRequestDTO struct {
	Title   string          `json:"title"`
	Content json.RawMessage `json:"content"`
	IsStub  bool            `json:"isStub"`
}

// NodesHandler serves the node collection endpoint.
type nodesHandler struct {
	nodeService *service.NodeService
}

// NewNodesHandler creates an HTTP handler backed by the node service.
func NewNodesHandler(nodeService *service.NodeService) http.Handler {
	return &nodesHandler{nodeService: nodeService}
}

// ServeHTTP handles GET and POST requests for /api/nodes.
func (h *nodesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listNodes(w, r)
	case http.MethodPost:
		h.createNode(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *nodesHandler) listNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.nodeService.ListNodes(r.Context())
	if err != nil {
		log.Printf("list nodes failed: %v", err)
		http.Error(w, "failed to list nodes", http.StatusInternalServerError)
		return
	}

	responses := make([]NodeResponseDTO, 0, len(nodes))
	for _, node := range nodes {
		responses = append(responses, mapDomainToNodeResponse(node))
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *nodesHandler) createNode(w http.ResponseWriter, r *http.Request) {
	var request createNodeRequestDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if request.Content == nil {
		request.Content = json.RawMessage(`{"type":"doc","content":[]}`)
	}

	node, err := h.nodeService.CreateNode(r.Context(), domain.CreateNodeInput{
		Title:   strings.TrimSpace(request.Title),
		Content: request.Content,
		IsStub:  request.IsStub,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidNode) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("create node failed: %v", err)
		http.Error(w, "failed to create node", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, mapDomainToNodeResponse(node))
}

var (
	defaultNodesHandler http.Handler
	defaultNodesOnce    sync.Once
	defaultNodesErr     error
)

// NodesHandler is the Vercel serverless function entry point.
func NodesHandler(w http.ResponseWriter, r *http.Request) {
	defaultNodesOnce.Do(func() {
		appConfig, err := config.Load()
		if err != nil {
			defaultNodesErr = err
			return
		}
		nodeRepository, err := postgres.NewNodeRepository(appConfig.DatabaseURL)
		if err != nil {
			defaultNodesErr = err
			return
		}
		defaultNodesHandler = NewNodesHandler(service.NewNodeService(nodeRepository))
	})
	if defaultNodesErr != nil {
		http.Error(w, "database configuration is unavailable", http.StatusInternalServerError)
		return
	}
	defaultNodesHandler.ServeHTTP(w, r)
}

func mapDomainToNodeResponse(node domain.NodeDomain) NodeResponseDTO {
	return NodeResponseDTO{
		ID:        node.ID,
		Title:     node.Title,
		Content:   node.Content,
		IsStub:    node.IsStub,
		CreatedAt: node.CreatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
		UpdatedAt: node.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
