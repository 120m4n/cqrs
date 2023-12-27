package events

import (
	"context"

	"github.com/120m4n/cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SubscribeCreatedFeed(ctx context.Context) (<-chan *CreatedFeedMessage, error)
	OnCreatedFeed(ctx context.Context, handler func(*CreatedFeedMessage)) error
}

var eventStore EventStore

func SetEventStore(es EventStore) {
	eventStore = es
}

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan *CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreatedFeed(ctx context.Context, handler func(*CreatedFeedMessage)) error {
	return eventStore.OnCreatedFeed(ctx, handler)
}
