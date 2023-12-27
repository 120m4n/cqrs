package database

import (
	"context"
	"database/sql"

	"github.com/120m4n/cqrs/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO feeds (title, url, description)
		VALUES ($1, $2, $3)
	`, feed.Title, feed.URL, feed.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, url, description, created_at
		FROM feeds
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := make([]*models.Feed, 0)
	for rows.Next() {
		feed := new(models.Feed)
		err := rows.Scan(&feed.ID, &feed.Title, &feed.URL, &feed.Description, &feed.CreatedAt)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}
