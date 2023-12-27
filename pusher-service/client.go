package main

import (
	"time"

	"github.com/gorilla/websocket"
)

var newline = []byte{'\n'}

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait = 10 * time.Second
)

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	client := &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
	go client.writePump()
	return client
}

func (c *Client) startPinging() {
	ticker := time.NewTicker(pingPeriod)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) Close() {
	c.socket.Close()
	close(c.outbound)
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()
    for {
        select {
        case message, ok := <-c.outbound:
            c.socket.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // The hub closed the channel.
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.socket.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued chat messages to the current websocket message.
            n := len(c.outbound)
            for i := 0; i < n; i++ {
                w.Write(newline)
                w.Write(<-c.outbound)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.socket.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}