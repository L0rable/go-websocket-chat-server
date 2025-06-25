package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

type LobbyRoom struct {
	clients map[string]*Client
	rooms   map[int]*Room
}

type ClientReq struct {
	ClientId   string `json:"clientId"`
	ClientName string `json:"clientName"`
	RoomNo     int    `json:"room"`
}

func newLobbyRoom() *LobbyRoom {
	return &LobbyRoom{
		clients: make(map[string]*Client),
		rooms:   make(map[int]*Room),
	}
}

func (wRoom *LobbyRoom) checkJoinRoom(roomNo int) *Room {
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

func (wRoom *LobbyRoom) checkJoinClient(id string, name string, room *Room) *Client {
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

func (wRoom *LobbyRoom) handleJoinReq(w http.ResponseWriter, r *http.Request) string {
	var joinReq *ClientReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&joinReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Fatal("Bad http request: ", err, " (main.go, serveClientJoin())")
	}

	clientId := url.QueryEscape(joinReq.ClientId)
	if clientId == "" {
		clientId = uuid.New().String()
	}
	clientName := url.QueryEscape(joinReq.ClientName)
	roomNo := url.QueryEscape(strconv.Itoa(joinReq.RoomNo))
	redirectURL := "/room?clientId=" + clientId + "&clientName=" + clientName + "&roomNo=" + roomNo

	return redirectURL
}

func (wRoom *LobbyRoom) handleWsReq(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	clientName := r.URL.Query().Get("clientName")
	roomNo, err := strconv.Atoi(r.URL.Query().Get("roomNo"))
	if err != nil {
		log.Fatal("roomNoInt error: ", err, " (main.go, serveClientWs())")
	}
	if clientName == "" || roomNo == 0 {
		log.Fatal("Invalid client join req (waiting_room.go, newJoin())")
	}

	room := wRoom.checkJoinRoom(roomNo)
	client := wRoom.checkJoinClient(clientId, clientName, room)
	client.openWs(w, r)

	wRoom.clients[client.name] = client
	room.join <- client
}

func (wRoom *LobbyRoom) handleLeaveReq(w http.ResponseWriter, r *http.Request) string {
	var leaveReq *ClientReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&leaveReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Fatal("Bad http request: ", err, " (main.go, serveClientLeave())")
	}

	clientId := url.QueryEscape(leaveReq.ClientId)
	clientName := url.QueryEscape(leaveReq.ClientName)
	roomNo := url.QueryEscape(strconv.Itoa(leaveReq.RoomNo))
	redirectURL := "/?clientId=" + clientId + "&clientName=" + clientName + "&roomNo=" + roomNo

	return redirectURL
}
