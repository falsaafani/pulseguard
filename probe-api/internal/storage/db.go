package storage

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

type Target struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type Probe struct {
	ID         int       `json:"id"`
	TargetID   int       `json:"target_id"`
	Timestamp  time.Time `json:"timestamp"`
	LatencyMS  int       `json:"latency_ms"`
	StatusCode int       `json:"status_code"`
	OK         bool      `json:"ok"`
}

type Incident struct {
	ID        int       `json:"id"`
	TargetID  int       `json:"target_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Kind      string    `json:"kind"`
	Details   string    `json:"details"`
}

func NewDB(connectionString string) (*DB, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

// TODO: Implement CRUD methods for targets, probes, incidents
