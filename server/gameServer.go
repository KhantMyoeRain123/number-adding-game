package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			return origin == "http://localhost:5173"
		},
	}
)

type GameServer struct {
	RoomCodeToUser map[string][]Player //a mapping between room code and the users in the room
	PlayerList     map[Player]bool
}

func NewGameServer() *GameServer {

	return &GameServer{
		RoomCodeToUser: make(map[string][]Player),
		PlayerList:     make(map[Player]bool),
	}
}

func (gs *GameServer) AddPlayer(p *Player) {
	return
}

// API handlers start here
// upgrades to websocket
func (gs *GameServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	defer func() {
		log.Println("Closing Connection...")
		err := conn.Close()

		if err != nil {
			log.Println("Could not close connection: ", err)
		}
	}()
	if err != nil {
		log.Println(err)
		return
	}
	player := NewPlayer(conn, gs)

	gs.AddPlayer(player)

	go player.ReadMessages()
	go player.WriteMessages()
}

// POST:make a room
func (gs *GameServer) MakeRoom(w http.ResponseWriter, r *http.Request) {
	//create a random room code
	//return json success message with room code
}

// POST:join a room using the room code
func (gs *GameServer) JoinRoom(w http.ResponseWriter, r *http.Request) {
	//index into the room to player map
	//join the room and return success message if the room exists, error if not
}

// test endpoints
func ApiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API root!")
}

func ApiTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API test!")
}
