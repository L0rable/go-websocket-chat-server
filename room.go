package main

import (
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
	clientId string
	text     string
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
			for _, msg := range room.messages {
				clientJoin.sendBuff <- []byte(msg.clientId + ": " + msg.text + "\n")
			}
		case clientLeave := <-room.leave:
			_, clientExists := room.clients[clientLeave]
			if clientExists {
				log.Println("client ", clientLeave)
				delete(room.clients, clientLeave)
				close(clientLeave.sendBuff)
			}

		case msg := <-room.broadcast:
			room.messages = append(room.messages, msg)
			log.Println(msg.text)
			for client := range room.clients {
				select {
				case client.sendBuff <- []byte(msg.clientId + ": " + msg.text + "\n"):
				default:
					close(client.sendBuff)
					delete(room.clients, client)
				}
			}
		}
	}
}
