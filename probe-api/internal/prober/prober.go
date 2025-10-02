package prober

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yourusername/pulseguard/probe-api/internal/storage"
)

type Prober struct {
	db   *storage.DB
	nats *nats.Conn
}

func New(db *storage.DB, nc *nats.Conn) *Prober {
	return &Prober{
		db:   db,
		nats: nc,
	}
}

func NewNATSConnection(url string) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	return nats.Connect(url)
}

func (p *Prober) Start() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Prober started")

	for range ticker.C {
		// TODO: Fetch targets from DB
		// TODO: Probe each target
		// TODO: Publish results to NATS
		// TODO: Store results in DB
		log.Println("Probing targets...")
	}
}

// TODO: Implement HTTP/TCP/ICMP probe functions
