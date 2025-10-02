package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	webhookURL string
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
}

func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
	}
}

func (c *Client) SendIncidentAlert(targetName, kind, details string) error {
	message := SlackMessage{
		Text: "🚨 Incident Detected",
		Attachments: []Attachment{
			{
				Color:  "danger",
				Title:  fmt.Sprintf("%s - %s", targetName, kind),
				Text:   details,
				Footer: "PulseGuard",
			},
		},
	}

	return c.sendMessage(message)
}

func (c *Client) sendMessage(message SlackMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d", resp.StatusCode)
	}

	return nil
}
