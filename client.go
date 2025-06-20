package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	id       string
	name     string
	room     *Room
	roomConn *websocket.Conn
	sendBuff chan []byte
}

func newClient(id string, name string, room *Room) *Client {
	client := &Client{
		id:       id,
		name:     name,
		room:     room,
		roomConn: nil,
		sendBuff: nil,
	}
	return client
}

func (c *Client) readPump() {
	if c.room != nil {
		defer func() {
			c.room.leave <- c
			c.roomConn.Close()
			c.room = nil
		}()

		for {
			_, msg, err := c.roomConn.ReadMessage()
			if err != nil {
				log.Println("roomConn.ReadMessage() ", err, " (client.go, readPump())")
				return
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

func (c *Client) openWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Websocket upgrade failed: ", err, " (client.go, openWs())")
	}
	c.roomConn = conn
	c.sendBuff = make(chan []byte, 256)

	go c.readPump()
	go c.writePump()
}
