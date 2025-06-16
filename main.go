package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "URL not found", http.StatusNotFound)
		log.Fatal("URL not found: ", r.URL.Path, " (main.go, serveIndex()")
	}
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("HTTP method not allowed: ", r.Method, " (main.go, serveIndex()")
	}
	http.ServeFile(w, r, "index.html")
}

func serveRoom(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/room" {
		http.Error(w, "URL not found", http.StatusNotFound)
		log.Fatal("URL not found: ", r.URL.Path, " (main.go, serveIndex()")
	}
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("HTTP method not allowed: ", r.Method, " (main.go, serveRoom()")
	}
	http.ServeFile(w, r, "room.html")
}

func serveClientWs(w http.ResponseWriter, r *http.Request, wRoom *WaitingRoom) {
	clientName := r.URL.Query().Get("clientName")
	roomNo := r.URL.Query().Get("roomNo")
	roomNoInt, err := strconv.Atoi(roomNo)
	if err != nil {
		log.Fatal("roomNoInt error: ", err, " (main.go, /ws)")
	}

	joinReq := &JoinReq{ClientName: clientName, Room: roomNoInt}
	wRoom.newJoin(w, r, joinReq)
}

func serveClientJoin(w http.ResponseWriter, r *http.Request) {
	var joinReq *JoinReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&joinReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Fatal("Bad http request: ", err, " (main.go, /join)")
	}

	clientName := joinReq.ClientName
	roomNo := strconv.Itoa(joinReq.Room)
	redirectURL := "/room?clientName=" + url.QueryEscape(clientName) + "&roomNo=" + url.QueryEscape(roomNo)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func main() {
	wRoom := newWaitingRoom()
	go wRoom.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/room", serveRoom)
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		serveClientJoin(w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveClientWs(w, r, wRoom)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
