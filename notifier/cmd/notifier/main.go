package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/pulseguard/notifier/internal/listener"
	"github.com/yourusername/pulseguard/notifier/internal/slack"
)

func main() {
	// Initialize Slack client
	slackClient := slack.NewClient(os.Getenv("SLACK_WEBHOOK_URL"))

	// Initialize incident listener
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	listener, err := listener.New(natsURL, slackClient)
	if err != nil {
		log.Fatalf("Failed to initialize listener: %v", err)
	}
	defer listener.Close()

	// Start listening for incidents
	if err := listener.Start(); err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	log.Println("Notifier started, listening for incidents...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down notifier...")
}
