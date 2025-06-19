package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "URL not found", http.StatusNotFound)
		log.Fatal("URL not found: ", r.URL.Path, " (main.go, serveIndex())")
	}
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("HTTP method not allowed: ", r.Method, " (main.go, serveIndex())")
	}
	http.ServeFile(w, r, "index.html")
}

func serveRoom(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/room" {
		http.Error(w, "URL not found", http.StatusNotFound)
		log.Fatal("URL not found: ", r.URL.Path, " (main.go, serveIndex())")
	}
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("HTTP method not allowed: ", r.Method, " (main.go, serveRoom())")
	}
	http.ServeFile(w, r, "room.html")
}

func serveClientWs(w http.ResponseWriter, r *http.Request, wRoom *WaitingRoom) {
	clientId := r.URL.Query().Get("clientId")
	clientName := r.URL.Query().Get("clientName")
	roomNo := r.URL.Query().Get("roomNo")
	roomNoInt, err := strconv.Atoi(roomNo)
	if err != nil {
		log.Fatal("roomNoInt error: ", err, " (main.go, serveClientWs())")
	}

	joinReq := &JoinReq{ClientId: clientId, ClientName: clientName, RoomNo: roomNoInt}
	wRoom.newJoin(w, r, joinReq)
}

func serveClientJoin(w http.ResponseWriter, r *http.Request) {
	var joinReq *JoinReq
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

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func serveClientLeave(w http.ResponseWriter, r *http.Request) {
	var leaveReq *LeaveReq
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

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func main() {
	wRoom := newWaitingRoom()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/room", serveRoom)
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		serveClientJoin(w, r)
	})
	http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
		serveClientLeave(w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveClientWs(w, r, wRoom)
	})

	err := http.ListenAndServe(":8080", nil)
	// log.Println("Go to localhost:8080 to access the chat application")
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
