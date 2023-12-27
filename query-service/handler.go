package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/120m4n/cqrs/events"
	"github.com/120m4n/cqrs/models"
	"github.com/120m4n/cqrs/search"
	"github.com/120m4n/cqrs/repository"
)

func onCreatedFeed(m *events.CreatedFeedMessage) {
	feed := &models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Print(err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	feeds, err := repository.GetFeeds  (ctx)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)

}

func searchFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeeds(ctx, query)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)

}
