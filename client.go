package main

import (
	"log"
	"net/http"

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

/*func openWsReq(wRoom *WaitingRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade err ", err)
	}

	id := uuid.New()
	client := &Client{id: id.String(), name: "", room: nil, roomConn: conn, sendBuff: make(chan []byte, 256)}
	client.room.join <- client

	go client.waiting(wRoom)
	go client.readPump()
	go client.writePump()
}*/

func newClient(w http.ResponseWriter, r *http.Request, wRoom *WaitingRoom, joinMsg *JoinReq) {
	id := uuid.New()
	clientName := joinMsg.ClientName
	roomNo := joinMsg.Room
	client := &Client{id: id.String(), name: clientName, room: nil, roomConn: nil, sendBuff: make(chan []byte, 256)}

	log.Println(joinMsg)

	if clientName == "" || roomNo == 0 {
		log.Println("Invalid client join req")
		return
	}

	wRoom.joinRoom <- &JoinRoom{client, roomNo}
}
