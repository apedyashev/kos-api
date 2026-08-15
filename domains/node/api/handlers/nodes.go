package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"kos-api/domains/node/api/dto"
	"kos-api/domains/node/service"
	"kos-api/domains/node/service/domain"
)

// NodesHandler handles the node collection HTTP endpoint.
type NodesHandler struct {
	nodeService service.Service
	debugErrors bool
}

// NewNodesHandler creates a node handler backed by the node service contract.
func NewNodesHandler(nodeService service.Service, debugErrors bool) *NodesHandler {
	return &NodesHandler{nodeService: nodeService, debugErrors: debugErrors}
}

// ServeHTTP handles GET and POST requests for /api/nodes.
func (h *NodesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *NodesHandler) listNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.nodeService.ListNodes(r.Context())
	if err != nil {
		log.Printf("list nodes failed: %v", err)
		if h.debugErrors {
			http.Error(w, fmt.Sprintf("failed to list nodes: %v", err), http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to list nodes", http.StatusInternalServerError)
		return
	}

	responses := make([]dto.NodeResponseDTO, 0, len(nodes))
	for _, currentNode := range nodes {
		responses = append(responses, MapDomainToResponse(currentNode))
	}
	writeJSON(w, http.StatusOK, responses)
}

func (h *NodesHandler) createNode(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateNodeRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	createdNode, err := h.nodeService.CreateNode(r.Context(), MapCreateRequestToDomain(request))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidNode) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("create node failed: %v", err)
		if h.debugErrors {
			http.Error(w, fmt.Sprintf("failed to create node: %v", err), http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to create node", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, MapDomainToResponse(createdNode))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
