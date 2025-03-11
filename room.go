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
	broadcast chan []byte
	// 2D byte slice to hold the message history
	messages [][]byte
}

func newRoom() *Room {
	return &Room{
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte),
		messages:  make([][]byte, 0),
	}
}

func (room *Room) run() {
	for {
		select {
		case clientJoin := <-room.join:
			log.Println("client join")
			room.clients[clientJoin] = true

		case clientLeave := <-room.leave:
			_, clientExists := room.clients[clientLeave]
			if clientExists {
				log.Println("client ", clientLeave)
				delete(room.clients, clientLeave)
				close(clientLeave.sendBuff)
			}

		case msg := <-room.broadcast:
			room.messages = append(room.messages, msg)
			log.Println(string(msg))
			for client := range room.clients {
				select {
				case client.sendBuff <- msg:
				default:
					close(client.sendBuff)
					delete(room.clients, client)
				}
			}
		}
	}
}
