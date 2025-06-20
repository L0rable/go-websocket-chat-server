package main

import (
	"log"
	"net/http"
)

type WaitingRoom struct {
	clients map[string]*Client
	rooms   map[int]*Room
	// joinRoom  chan *JoinRoom
	// leaveRoom chan int
}

type ClientReq struct {
	ClientId   string `json:"clientId"`
	ClientName string `json:"clientName"`
	RoomNo     int    `json:"room"`
}

func newWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		clients: make(map[string]*Client),
		rooms:   make(map[int]*Room),
		// joinRoom:  make(chan *JoinRoom),
		// leaveRoom: make(chan int),
	}
}

func (wRoom *WaitingRoom) checkJoinRoom(roomNo int) *Room {
	var joinRoom *Room
	if wRoom.rooms[roomNo] == nil {
		joinRoom = newRoom(roomNo)
		go joinRoom.run()
	} else {
		joinRoom = wRoom.rooms[roomNo]
	}

	wRoom.rooms[roomNo] = joinRoom
	return joinRoom
}

func (wRoom *WaitingRoom) checkJoinClient(id string, name string, room *Room) *Client {
	var client *Client
	if wRoom.clients[id] == nil {
		client = newClient(id, name, room)
		wRoom.clients[id] = client
	} else {
		client = wRoom.clients[id]
		client.name = name
		client.room = room
	}
	return client
}

func (wRoom *WaitingRoom) newJoin(w http.ResponseWriter, r *http.Request, joinReq *ClientReq) {
	clientId := joinReq.ClientId
	clientName := joinReq.ClientName
	roomNo := joinReq.RoomNo
	if clientName == "" || roomNo == 0 {
		log.Fatal("Invalid client join req (waiting_room.go, newJoin())")
	}

	room := wRoom.checkJoinRoom(roomNo)
	client := wRoom.checkJoinClient(clientId, clientName, room)
	client.openWs(w, r)

	wRoom.clients[client.name] = client
	room.join <- client
}
