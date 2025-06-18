package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "URL not found", http.StatusNotFound)
	// 	log.Fatal("URL not found: ", r.URL.Path, " (main.go, serveIndex())")
	// }
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
	clientName := r.URL.Query().Get("clientName")
	roomNo := r.URL.Query().Get("roomNo")
	roomNoInt, err := strconv.Atoi(roomNo)
	if err != nil {
		log.Fatal("roomNoInt error: ", err, " (main.go, /ws)")
	}

	joinReq := &JoinReq{ClientName: clientName, RoomNo: roomNoInt}
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

	clientName := url.QueryEscape(joinReq.ClientName)
	roomNo := url.QueryEscape(strconv.Itoa(joinReq.RoomNo))
	redirectURL := "/room?clientName=" + clientName + "&roomNo=" + roomNo

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func serveClientLeave(w http.ResponseWriter, r *http.Request) {
	var leaveReq *LeaveReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&leaveReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Fatal("Bad http request: ", err, " (main.go, /join)")
	}

	clientName := url.QueryEscape(leaveReq.ClientName)
	roomNo := url.QueryEscape(strconv.Itoa(leaveReq.RoomNo))
	redirectURL := "/?clientName=" + clientName + "&roomNo=" + roomNo

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
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
