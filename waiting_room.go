package main

type WaitingRoom struct {
	clients   map[*Client]bool
	rooms     map[*Room]bool
	joinRoom  chan *JoinRoom
	leaveRoom chan int
}

type JoinReq struct {
	ClientName string `json:"clientName"`
	Room       int    `json:"room"`
}

type JoinRoom struct {
	Client *Client
	RoomNo int
}

func newWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		clients:   make(map[*Client]bool),
		rooms:     make(map[*Room]bool),
		joinRoom:  make(chan *JoinRoom),
		leaveRoom: make(chan int),
	}
}

func (wRoom *WaitingRoom) run() {
	for {
		select {
		case join := <-wRoom.joinRoom:
			// Currenly doesn't check if room exists
			newRoom := newRoom()
			go newRoom.run()

			join.Client.room = newRoom
			newRoom.join <- join.Client
		}
	}
}
