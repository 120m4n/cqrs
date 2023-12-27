package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/120m4n/cqrs/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS" default:"localhost:4222"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	hub := NewHub()

	n, err := events.NewNatsEventStore(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}
	err = n.OnCreatedFeed(context.Background(), func(m *events.CreatedFeedMessage) {
		hub.Broadcast(newCreatedFeedMessage(strconv.FormatInt(m.ID, 10), m.Title, m.Description, m.CreatedAt), nil)
	})
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)
	defer events.Close()

	go hub.Run()
	http.HandleFunc("/ws", hub.HandleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
