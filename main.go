package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "URL not found", http.StatusNotFound)
	}
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
	}
	http.ServeFile(w, r, "landing.html")
}

func main() {
	wRoom := newWaitingRoom()
	go wRoom.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(wr http.ResponseWriter, req *http.Request) {
		log.Println("openWsReq")
		// openWsReq(wRoom, wr, req)
	})
	http.HandleFunc("/join", func(wr http.ResponseWriter, req *http.Request) {
		var joinReq *JoinReq
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&joinReq)
		if err != nil {
			http.Error(wr, "Bad request", http.StatusBadRequest)
			return
		}
		newClient(wr, req, wRoom, joinReq)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
