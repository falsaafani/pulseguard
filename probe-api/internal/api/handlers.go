package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/pulseguard/probe-api/internal/storage"
)

type Handler struct {
	db *storage.DB
}

func New(db *storage.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/targets", h.CreateTarget).Methods("POST")
	r.HandleFunc("/targets", h.ListTargets).Methods("GET")
	r.HandleFunc("/status", h.GetStatus).Methods("GET")
	r.HandleFunc("/metrics", h.GetMetrics).Methods("GET")
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type CreateTargetRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

func (h *Handler) CreateTarget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate request
	if req.Name == "" || req.URL == "" || req.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "name, url, and type are required"})
		return
	}

	// Create target
	target, err := h.db.CreateTarget(req.Name, req.URL, req.Type)
	if err != nil {
		log.Printf("Error creating target: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create target"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(target)
}

func (h *Handler) ListTargets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	targets, err := h.db.ListTargets()
	if err != nil {
		log.Printf("Error listing targets: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to list targets"})
		return
	}

	// Return empty array instead of null if no targets
	if targets == nil {
		targets = []storage.Target{}
	}

	json.NewEncoder(w).Encode(targets)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get optional target_id query parameter
	targetIDStr := r.URL.Query().Get("target_id")

	var probes []storage.Probe
	var err error

	if targetIDStr != "" {
		targetID, parseErr := strconv.Atoi(targetIDStr)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid target_id"})
			return
		}
		probes, err = h.db.GetRecentProbes(targetID, 100)
	} else {
		probes, err = h.db.GetAllRecentProbes(100)
	}

	if err != nil {
		log.Printf("Error fetching probes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch status"})
		return
	}

	// Return empty array instead of null if no probes
	if probes == nil {
		probes = []storage.Probe{}
	}

	json.NewEncoder(w).Encode(probes)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Prometheus metrics format
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("# PulseGuard Metrics\n# TODO: Implement Prometheus metrics\n"))
}
