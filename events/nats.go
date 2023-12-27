package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/120m4n/cqrs/models"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan *CreatedFeedMessage
}

func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{
		conn: conn,
	}, nil
}

func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}

	if n.feedCreatedChan != nil {
		close(n.feedCreatedChan)
	}
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (n *NatsEventStore) decodeMessage(data []byte, m interface{})  error {
	b := bytes.Buffer{}
	b.Write(data)

	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	m := &CreatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		URL:         feed.URL,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	return n.conn.Publish(m.Type(), data)
}

func (n *NatsEventStore) OnCreatedFeed(ctx context.Context, handler func(*CreatedFeedMessage)) (err error) {
  msg := &CreatedFeedMessage{}
  n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
	err = n.decodeMessage(m.Data, msg)
	handler(msg)
  })
  return
}

func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan *CreatedFeedMessage, error) {
	msg := &CreatedFeedMessage{}
	n.feedCreatedChan = make(chan *CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	
	n.feedCreatedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case m := <-ch:
				err = n.decodeMessage(m.Data, msg)
				if err != nil {
					continue
				}
				n.feedCreatedChan <- msg
			case <-ctx.Done():
				return
			}
		}
	}()

	return (<-chan *CreatedFeedMessage)(n.feedCreatedChan), nil
}

