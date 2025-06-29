package main

import (
	"log"
	"net/http"
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
	http.ServeFile(w, r, "client/index.html")
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
	http.ServeFile(w, r, "client/room.html")
}

func main() {
	lobbyRoom := newLobbyRoom()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("client/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("client/js"))))

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/room", serveRoom)
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		redirectURL := lobbyRoom.handleJoinReq(w, r)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	})
	http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
		redirectURL := lobbyRoom.handleLeaveReq(w, r)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		lobbyRoom.handleWsReq(w, r)
	})

	log.Println("About to listen on 8080. Go to https://localhost:8080/")
	err := http.ListenAndServeTLS(":8080", "tls/server.crt", "tls/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
