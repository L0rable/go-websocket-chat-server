package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	id   string
	name string
	// ref to connected room
	room *Room
	// client websocket connection to room
	roomConn *websocket.Conn
	// message buffer channel, send to room
	sendBuff chan []byte
}

func (c *Client) readPump() {
	if c.room != nil {
		defer func() {
			c.room = nil
			c.roomConn.Close()
			c.room.leave <- c
		}()

		for {
			_, msg, err := c.roomConn.ReadMessage()
			if err != nil {
				break
			}
			c.room.broadcast <- &Message{ClientId: c.id, Text: string(msg)}
		}
	}
}

func (c *Client) writePump() {
	if c.room != nil {
		for message := range c.sendBuff {
			w, err := c.roomConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// log.Println(string(message))
			w.Write(message)

			n := len(c.sendBuff)
			for i := 0; i < n; i++ {
				w.Write(<-c.sendBuff)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func newClient(clientName string, room *Room, conn *websocket.Conn) *Client {
	id := uuid.New()
	client := &Client{id: id.String(), name: clientName, room: room, roomConn: conn, sendBuff: make(chan []byte, 256)}

	go client.readPump()
	go client.writePump()

	return client
}
