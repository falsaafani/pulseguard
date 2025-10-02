package api

import (
	"encoding/json"
	"net/http"

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
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) CreateTarget(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement target creation
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) ListTargets(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement target listing
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement status retrieval
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Prometheus metrics
	w.WriteHeader(http.StatusNotImplemented)
}
