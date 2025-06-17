package main

import (
	"log"
	"net/http"
)

type WaitingRoom struct {
	clients map[*Client]bool
	rooms   map[int]*Room
	// joinRoom  chan *JoinRoom
	// leaveRoom chan int
}

type JoinReq struct {
	ClientName string `json:"clientName"`
	Room       int    `json:"room"`
}

func newWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		clients: make(map[*Client]bool),
		rooms:   make(map[int]*Room),
		// joinRoom:  make(chan *JoinRoom),
		// leaveRoom: make(chan int),
	}
}

func (wRoom *WaitingRoom) checkJoinRoom(roomNo int) *Room {
	var joinRoom *Room
	if wRoom.rooms[roomNo] == nil {
		joinRoom = newRoom(roomNo)
	} else {
		joinRoom = wRoom.rooms[roomNo]
	}

	wRoom.rooms[roomNo] = joinRoom
	return joinRoom
}

func (wRoom *WaitingRoom) newJoin(w http.ResponseWriter, r *http.Request, joinReq *JoinReq) {
	clientName := joinReq.ClientName
	roomNo := joinReq.Room
	if clientName == "" || roomNo == 0 {
		log.Fatal("Invalid client join req (client.go, newJoin())")
	}

	conn := openWs(w, r)
	room := wRoom.checkJoinRoom(roomNo)
	newClient := newClient(clientName, room, conn)
	room.join <- newClient
}
