package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/120m4n/cqrs/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/mitchellh/mapstructure"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElasticSearchRepository(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{
		client: client,
	}, nil
}

func (e *ElasticSearchRepository) Close() {
	// Compare this snippet from search/elastic.go:
}

func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := e.client.Index(
		"feeds",
		bytes.NewReader(body),
		e.client.Index.WithDocumentID(strconv.FormatInt(feed.ID, 10)),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("wait_for"),
	)

	return err
}

func (e *ElasticSearchRepository) SearchFeeds(ctx context.Context, query string) (results []*models.Feed, err error) {
	var buf bytes.Buffer

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
				"fields": []string{
					"title",
					"description",
				},
				"fuzziness": 3,
				// "cut_off_frequency": 0.0001, // Remove this line
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("feeds"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New("elasticsearch error " + res.String())
	}

	var eRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		if errors.Is(err, io.EOF) {
			// Handle EOF error specifically
			log.Println("Reached the end of the file or stream unexpectedly")
		} else {
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	feeds := make([]*models.Feed, 0, len(eRes))
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to convert _source to map[string]interface{}")
		}
		feed := new(models.Feed)
		if err := mapstructure.Decode(source, &feed); err != nil {
			return nil, fmt.Errorf("failed to decode source to feed: %w", err)
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}
