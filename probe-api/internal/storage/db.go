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

// Target operations

func (db *DB) CreateTarget(name, url, targetType string) (*Target, error) {
	query := `
		INSERT INTO targets (name, url, type, enabled)
		VALUES ($1, $2, $3, true)
		RETURNING id, name, url, type, enabled, created_at
	`

	target := &Target{}
	err := db.conn.QueryRow(query, name, url, targetType).Scan(
		&target.ID,
		&target.Name,
		&target.URL,
		&target.Type,
		&target.Enabled,
		&target.CreatedAt,
	)

	return target, err
}

func (db *DB) ListTargets() ([]Target, error) {
	query := `
		SELECT id, name, url, type, enabled, created_at
		FROM targets
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []Target
	for rows.Next() {
		var target Target
		err := rows.Scan(
			&target.ID,
			&target.Name,
			&target.URL,
			&target.Type,
			&target.Enabled,
			&target.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, rows.Err()
}

func (db *DB) GetEnabledTargets() ([]Target, error) {
	query := `
		SELECT id, name, url, type, enabled, created_at
		FROM targets
		WHERE enabled = true
		ORDER BY id
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []Target
	for rows.Next() {
		var target Target
		err := rows.Scan(
			&target.ID,
			&target.Name,
			&target.URL,
			&target.Type,
			&target.Enabled,
			&target.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, rows.Err()
}

// Probe operations

func (db *DB) CreateProbe(targetID, latencyMS, statusCode int, ok bool) error {
	query := `
		INSERT INTO probes (target_id, ts, latency_ms, status_code, ok)
		VALUES ($1, NOW(), $2, $3, $4)
	`

	_, err := db.conn.Exec(query, targetID, latencyMS, statusCode, ok)
	return err
}

func (db *DB) GetRecentProbes(targetID int, limit int) ([]Probe, error) {
	query := `
		SELECT id, target_id, ts, latency_ms, status_code, ok
		FROM probes
		WHERE target_id = $1
		ORDER BY ts DESC
		LIMIT $2
	`

	rows, err := db.conn.Query(query, targetID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var probes []Probe
	for rows.Next() {
		var probe Probe
		err := rows.Scan(
			&probe.ID,
			&probe.TargetID,
			&probe.Timestamp,
			&probe.LatencyMS,
			&probe.StatusCode,
			&probe.OK,
		)
		if err != nil {
			return nil, err
		}
		probes = append(probes, probe)
	}

	return probes, rows.Err()
}

func (db *DB) GetAllRecentProbes(limit int) ([]Probe, error) {
	query := `
		SELECT id, target_id, ts, latency_ms, status_code, ok
		FROM probes
		ORDER BY ts DESC
		LIMIT $1
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var probes []Probe
	for rows.Next() {
		var probe Probe
		err := rows.Scan(
			&probe.ID,
			&probe.TargetID,
			&probe.Timestamp,
			&probe.LatencyMS,
			&probe.StatusCode,
			&probe.OK,
		)
		if err != nil {
			return nil, err
		}
		probes = append(probes, probe)
	}

	return probes, rows.Err()
}
