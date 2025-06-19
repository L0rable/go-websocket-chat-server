package main

import (
	"encoding/json"
	"log"
)

type Room struct {
	id        int
	clients   map[string]*Client
	join      chan *Client
	leave     chan *Client
	broadcast chan *Message
	messages  []*Message
}

type Message struct {
	ClientName string `json:"clientName"`
	Text       string `json:"text"`
}

func newRoom(roomNo int) *Room {
	room := &Room{
		id:        roomNo,
		clients:   make(map[string]*Client),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *Message),
	}

	go room.run()

	return room
}

func (room *Room) run() {
	for {
		select {
		case clientJoin := <-room.join:
			log.Println("client join:", clientJoin.name)
			room.clients[clientJoin.id] = clientJoin
			msgs, err := json.Marshal(room.messages)
			if err != nil {
				log.Fatal("json.Marshal:", err, "(room.go, run())")
			}
			clientJoin.sendBuff <- msgs

		case clientLeave := <-room.leave:
			if room.clients[clientLeave.id] != nil {
				log.Println("client leave: ", clientLeave.name)
				delete(room.clients, clientLeave.id)
				close(clientLeave.sendBuff)
			}

		case msg := <-room.broadcast:
			room.messages = append(room.messages, msg)
			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Fatal("json.Marshal:", err, "(room.go, run())")
			}
			for clientId, client := range room.clients {
				select {
				case client.sendBuff <- msgJson:
				default:
					close(client.sendBuff)
					delete(room.clients, clientId)
				}
			}
		}
	}
}
