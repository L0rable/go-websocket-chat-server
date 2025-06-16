package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WaitingRoom struct {
	clients   map[*Client]bool
	rooms     map[*Room]bool
	joinRoom  chan *JoinRoom
	leaveRoom chan int
}

type JoinReq struct {
	ClientName string `json:"clientName"`
	Room       int    `json:"room"`
}

type JoinRoom struct {
	client *Client
	room   *Room
}

func newWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		clients:   make(map[*Client]bool),
		rooms:     make(map[*Room]bool),
		joinRoom:  make(chan *JoinRoom),
		leaveRoom: make(chan int),
	}
}

func (wRoom *WaitingRoom) newJoin(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, joinReq *JoinReq) {
	clientName := joinReq.ClientName
	roomNo := joinReq.Room
	if clientName == "" || roomNo == 0 {
		log.Fatal("Invalid client join req (client.go, newJoin())")
	}

	newRoom := newRoom(roomNo)
	newClient := newClient(clientName, newRoom, conn)
	wRoom.joinRoom <- &JoinRoom{newClient, newRoom}
}

func (wRoom *WaitingRoom) run() {
	for {
		select {
		case join := <-wRoom.joinRoom:
			// Currenly doesn't check if room exists
			newRoom := newRoom(join.room.id)
			go newRoom.run()
			log.Println("newRoom", newRoom.id)

			join.client.room = newRoom
			newRoom.join <- join.client

			log.Println("client", join.client.name)
		}
	}
}
