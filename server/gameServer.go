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
	RoomCodeToUser map[string][]int //a mapping between room code and the users in the room
}

func NewGameServer() *GameServer {

	return &GameServer{
		RoomCodeToUser: make(map[string][]int),
	}
}

// API handlers start here
// upgrades to websocket
func (gs *GameServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	_, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

// POST:make a room
func (gs *GameServer) MakeRoom(w http.ResponseWriter, r *http.Request) {
	//create a random player id
	//create a random room id
	//add the mapping to the game server room to player map
	//return json success message
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
