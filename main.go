package main

import (
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
	http.ServeFile(w, r, "index.html")
}

func main() {
	room := newRoom()
	go room.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/join", func(wr http.ResponseWriter, req *http.Request) {
		wsRequest(room, wr, req)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
