package repository

import (
	"context"

	"github.com/120m4n/cqrs/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
	GetFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}

func Close() {
	repository.Close()
}

func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

func GetFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.GetFeeds(ctx)
}


