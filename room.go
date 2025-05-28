package main

import (
	"encoding/json"
	"log"
)

type Room struct {
	// map (hash table), key pointer to client and value is bool
	clients map[*Client]bool
	// buffer to handle client join
	join chan *Client
	// buffer to handle client leave
	leave chan *Client
	// byte slice for incoming client broadcasts
	broadcast chan *Message
	// slice to hold the message history
	messages []*Message
}

type Message struct {
	ClientId string `json:"clientId"`
	Text     string `json:"text"`
}

func newRoom() *Room {
	return &Room{
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *Message),
	}
}

func (room *Room) run() {
	for {
		select {
		case clientJoin := <-room.join:
			log.Println("client join")
			room.clients[clientJoin] = true
			msgs, err := json.Marshal(room.messages)
			if err != nil {
				log.Println("room.run() case clientJoin, data: ", err)
			}
			clientJoin.sendBuff <- msgs
		case clientLeave := <-room.leave:
			_, clientExists := room.clients[clientLeave]
			if clientExists {
				log.Println("client ", clientLeave)
				delete(room.clients, clientLeave)
				close(clientLeave.sendBuff)
			}

		case msg := <-room.broadcast:
			room.messages = append(room.messages, msg)
			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Println("room.run() case room.broadcast, data: ", err)
			}
			for client := range room.clients {
				select {
				case client.sendBuff <- msgJson:
				default:
					close(client.sendBuff)
					delete(room.clients, client)
				}
			}
		}
	}
}
