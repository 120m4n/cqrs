package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/120m4n/cqrs/events"
	"github.com/120m4n/cqrs/models"
	"github.com/120m4n/cqrs/repository"
)

type createFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createdFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()

	feed := models.Feed{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// transmit event on nats

	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("error publishing feed created event: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
