package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	// ref to connected room
	room *Room
	// client websocket connection to room
	roomConn *websocket.Conn
	// message buffer channel, send to room
	sendBuff chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.roomConn.Close()
		c.room.leave <- c
	}()

	for {
		_, msg, err := c.roomConn.ReadMessage()
		if err != nil {
			break
		}
		c.room.broadcast <- msg
	}
}

func (c *Client) writePump() {
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

func openWsReq(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade err ", err)
	}

	client := &Client{room: room, roomConn: conn, sendBuff: make(chan []byte, 128)}
	client.room.join <- client

	for _, msg := range room.messages {
		conn.WriteMessage(websocket.TextMessage, msg)
	}

	go client.readPump()
	go client.writePump()
}
