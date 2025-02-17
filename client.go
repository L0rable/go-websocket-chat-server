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

func wsRequest(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade err ", err)
	}

	client := &Client{room: room, roomConn: conn, sendBuff: make(chan []byte, 128)}
	client.room.join <- client
}
