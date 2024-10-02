package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	ROOM_CODE_LENGTH = 6
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			/*origin := r.Header.Get("Origin")

			return origin == "http://localhost:5173"*/
			return true //this needs to be changed later on for security
		},
	}
)

type GameServer struct {
	RoomCodeToPlayer map[string][]Player //a mapping between room code and the users in the room
	PlayerList     map[*Player]bool
}

func NewGameServer() *GameServer {

	return &GameServer{
		RoomCodeToPlayer: make(map[string][]Player),
		PlayerList:     make(map[*Player]bool),
	}
}

func (gs *GameServer) AddPlayer(p *Player) {
	return
}


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
func generateRoomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	roomCode := make([]byte, length)
	for i := range roomCode {
		roomCode[i] = charset[rand.IntN(len(charset))]
	}
	return string(roomCode)
}
type ErrorResponse struct{
	Status string `json:"status"`
	Message string `json:"message"`
}

func SendJSONError(w http.ResponseWriter, message string, statusCode int){
	errorResponse:=ErrorResponse{
		Status:"error",
		Message: "message",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(errorResponse)
}
type MakeRoomResponse struct{
	Status   string `json:"status"`
	RoomCode string `json:"room_code"`
}
// REST handlers start here
// POST:make a room
func (gs *GameServer) MakeRoom(w http.ResponseWriter, r *http.Request) {
	roomCode:=generateRoomCode(ROOM_CODE_LENGTH)
	gs.RoomCodeToPlayer[roomCode]=make([]Player,5)
	log.Printf("Created a new room with code %s...",roomCode)
	response:=MakeRoomResponse{
		Status:"success",
		RoomCode: roomCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		SendJSONError(w,"Failed to encode response",http.StatusInternalServerError)
		return
	}	
}

// POST:join a room using the room code
func (gs *GameServer) JoinRoom(w http.ResponseWriter, r *http.Request) {

}

// test endpoints
func ApiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API root!")
}

func ApiTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API test!")
}
