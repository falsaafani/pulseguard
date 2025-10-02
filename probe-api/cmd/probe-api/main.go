package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/yourusername/pulseguard/probe-api/internal/api"
	"github.com/yourusername/pulseguard/probe-api/internal/prober"
	"github.com/yourusername/pulseguard/probe-api/internal/storage"
)

func main() {
	// Initialize database connection
	db, err := storage.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize NATS connection
	nc, err := prober.NewNATSConnection(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Initialize prober
	p := prober.New(db, nc)
	go p.Start()

	// Initialize API server
	r := mux.NewRouter()
	apiHandler := api.New(db)
	apiHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting probe-api on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
