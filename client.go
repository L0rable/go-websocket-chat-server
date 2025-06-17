package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	id   uuid.UUID
	name string
	// ref to connected room
	room *Room
	// client websocket connection to room
	roomConn *websocket.Conn
	// message buffer channel, send to room
	sendBuff chan []byte
}

func openWs(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Websocket upgrade failed: ", err, " (client.go, openWs())")
	}
	return conn
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
				log.Println(err, "(client.go, readPump())")
				// return
			}
			c.room.broadcast <- &Message{ClientName: c.name, Text: string(msg)}
		}
	}
}

func (c *Client) writePump() {
	if c.room != nil {
		for message := range c.sendBuff {
			w, err := c.roomConn.NextWriter(websocket.TextMessage)
			if err != nil {
				// log.Fatal("roomConn.Nextwriter(): ", err, " (client.go, writePump())")
				log.Println("roomConn.Nextwriter(): ", err, " (client.go, writePump())")
			}
			// log.Println(string(message))
			w.Write(message)

			n := len(c.sendBuff)
			for i := 0; i < n; i++ {
				w.Write(<-c.sendBuff)
			}

			if err := w.Close(); err != nil {
				// log.Fatal("w.Close(): ", err, " (client.go, writePump())")
				log.Println("w.Close(): ", err, " (client.go, writePump())")
			}
		}
	}
}

func newClient(clientName string, room *Room, conn *websocket.Conn) *Client {
	id := uuid.New()
	client := &Client{
		id:       id,
		name:     clientName,
		room:     room,
		roomConn: conn,
		sendBuff: make(chan []byte, 256),
	}

	go client.readPump()
	go client.writePump()

	return client
}
