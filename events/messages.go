package events

import "time"

type Message interface {
	Type() string
}

type CreatedFeedMessage struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (m *CreatedFeedMessage) Type() string {
	return "created_feed"
}