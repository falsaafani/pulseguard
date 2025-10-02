package listener

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/yourusername/pulseguard/notifier/internal/slack"
)

type Listener struct {
	nc          *nats.Conn
	slackClient *slack.Client
}

type IncidentEvent struct {
	TargetID  int    `json:"target_id"`
	TargetName string `json:"target_name"`
	Kind      string `json:"kind"`
	Details   string `json:"details"`
	StartedAt string `json:"started_at"`
}

func New(natsURL string, slackClient *slack.Client) (*Listener, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &Listener{
		nc:          nc,
		slackClient: slackClient,
	}, nil
}

func (l *Listener) Start() error {
	_, err := l.nc.Subscribe("incidents.*", func(msg *nats.Msg) {
		var event IncidentEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling incident: %v", err)
			return
		}

		log.Printf("Received incident: %+v", event)

		// Send notification to Slack
		if err := l.slackClient.SendIncidentAlert(event.TargetName, event.Kind, event.Details); err != nil {
			log.Printf("Error sending Slack notification: %v", err)
		}
	})

	return err
}

func (l *Listener) Close() {
	if l.nc != nil {
		l.nc.Close()
	}
}
