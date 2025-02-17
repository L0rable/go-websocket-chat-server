package main

type Room struct {
	// map (hash table), key pointer to client and value is bool
	clients map[*Client]bool
	// byte slice for incoming client broadcasts
	clientBroadcast chan []byte
	// 2D byte slice to hold the message history
	messages [][]byte
}

func newRoom() *Room {
	return &Room{
		clients:         make(map[*Client]bool),
		clientBroadcast: make(chan []byte),
		messages:        make([][]byte, 0),
	}
}
