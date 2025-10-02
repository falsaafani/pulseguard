-- Initialize PulseGuard database schema

CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(1024) NOT NULL,
    type VARCHAR(50) NOT NULL, -- http, tcp, icmp
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS probes (
    id SERIAL PRIMARY KEY,
    target_id INTEGER NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    latency_ms INTEGER NOT NULL,
    status_code INTEGER,
    ok BOOLEAN NOT NULL
);

-- Index for time-series queries
CREATE INDEX IF NOT EXISTS idx_probes_target_ts ON probes(target_id, ts DESC);

CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    target_id INTEGER NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    kind VARCHAR(100) NOT NULL, -- latency_spike, error_rate, downtime
    details TEXT
);

CREATE INDEX IF NOT EXISTS idx_incidents_target ON incidents(target_id, started_at DESC);

-- Insert sample targets for testing
INSERT INTO targets (name, url, type) VALUES
    ('Google', 'https://google.com', 'http'),
    ('GitHub', 'https://github.com', 'http'),
    ('Example', 'https://example.com', 'http')
ON CONFLICT DO NOTHING;
