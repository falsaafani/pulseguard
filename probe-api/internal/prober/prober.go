package prober

import (
	"encoding/json"
	"log"
	"net/http"
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

	log.Println("Prober started, checking targets every 30 seconds")

	// Run immediately on start
	p.probeAllTargets()

	for range ticker.C {
		p.probeAllTargets()
	}
}

func (p *Prober) probeAllTargets() {
	targets, err := p.db.GetEnabledTargets()
	if err != nil {
		log.Printf("Error fetching targets: %v", err)
		return
	}

	if len(targets) == 0 {
		log.Println("No enabled targets to probe")
		return
	}

	log.Printf("Probing %d targets...", len(targets))

	for _, target := range targets {
		result := p.probeTarget(target)

		// Store probe result in database
		if err := p.db.CreateProbe(result.TargetID, result.LatencyMS, result.StatusCode, result.OK); err != nil {
			log.Printf("Error storing probe result for target %d: %v", target.ID, err)
		}

		// Publish result to NATS
		if err := p.publishResult(result); err != nil {
			log.Printf("Error publishing probe result for target %d: %v", target.ID, err)
		}

		log.Printf("[%s] %s - %dms, status=%d, ok=%v",
			target.Type, target.Name, result.LatencyMS, result.StatusCode, result.OK)
	}
}

type ProbeResult struct {
	TargetID   int    `json:"target_id"`
	TargetName string `json:"target_name"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	LatencyMS  int    `json:"latency_ms"`
	StatusCode int    `json:"status_code"`
	OK         bool   `json:"ok"`
	Timestamp  string `json:"timestamp"`
}

func (p *Prober) probeTarget(target storage.Target) ProbeResult {
	result := ProbeResult{
		TargetID:   target.ID,
		TargetName: target.Name,
		URL:        target.URL,
		Type:       target.Type,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	switch target.Type {
	case "http", "https":
		result = p.probeHTTP(target, result)
	default:
		log.Printf("Unsupported target type: %s", target.Type)
		result.OK = false
		result.StatusCode = 0
		result.LatencyMS = 0
	}

	return result
}

func (p *Prober) probeHTTP(target storage.Target, result ProbeResult) ProbeResult {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Follow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	start := time.Now()

	resp, err := client.Get(target.URL)
	latency := time.Since(start)

	result.LatencyMS = int(latency.Milliseconds())

	if err != nil {
		log.Printf("Error probing %s: %v", target.URL, err)
		result.OK = false
		result.StatusCode = 0
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.OK = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}

func (p *Prober) publishResult(result ProbeResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return p.nats.Publish("probe.results", data)
}
