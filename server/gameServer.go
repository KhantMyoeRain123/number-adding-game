package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	ROOM_CODE_LENGTH = 6
	ROOM_WAITING     = 0
	ROOM_STARTED     = 1
	ROOM_ENDED       = 2
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

type RoomState struct {
	RoomPlayerList map[*Player]bool
	State          int
}

type GameServer struct {
	RoomCodeToState map[string]RoomState //a mapping between room code and the users in the room
	PlayerList      map[*Player]bool
	PlayerCodes     map[string]bool
	mu              sync.Mutex
}

func NewGameServer() *GameServer {

	return &GameServer{
		RoomCodeToState: make(map[string]RoomState),
		PlayerList:      make(map[*Player]bool),
		PlayerCodes:     make(map[string]bool),
	}
}

func generateCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	roomCode := make([]byte, length)
	for i := range roomCode {
		roomCode[i] = charset[rand.IntN(len(charset))]
	}
	return string(roomCode)
}

func (gs *GameServer) AddPlayer(p *Player, roomCode string) error {
	//add player to server
	if _, ok := gs.PlayerCodes[p.PlayerId]; ok {
		gs.PlayerList[p] = true
	} else {
		return errors.New("cannot find player with given id")
	}

	//add player to room
	if roomCodeToState, ok := gs.RoomCodeToState[roomCode]; ok {
		roomCodeToState.RoomPlayerList[p] = true
	} else {
		return errors.New("cannot find room with given roomCode")
	}
	return nil
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type WebsocketUpgradeResponse struct {
	PlayerId   string
	PlayerName string
	RoomCode   string
	Host       bool
}

// upgrades to websocket
func (gs *GameServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	playerName := r.URL.Query().Get("playerName")
	playerId := r.URL.Query().Get("playerId")
	roomCode := r.URL.Query().Get("roomCode")
	host, err := strconv.ParseBool(r.URL.Query().Get("host"))

	if err != nil {
		log.Println("Couldn't convert host parameter to bool...")
		return
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := NewPlayer(playerId, playerName, host, conn, gs)

	log.Println("Added player...")

	err = gs.AddPlayer(player, roomCode)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(gs.PlayerList)
	log.Println(gs.RoomCodeToState[roomCode].RoomPlayerList)

	go player.ReadMessages()
	go player.WriteMessages()

}

func SendJSONError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(errorResponse)
}

type RoomCodeResponse struct {
	Status   string `json:"status"`
	PlayerId string `json:"playerId"`
	RoomCode string `json:"room_code"`
	State    int    `json:"room_state"`
	Host     bool   `json:"host"`
}

// REST handlers start here
// POST:make a room
func (gs *GameServer) MakeRoom(w http.ResponseWriter, r *http.Request) {
	roomCode := generateCode(ROOM_CODE_LENGTH)
	playerId := generateCode(4)
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.PlayerCodes[playerId] = true
	//create the new room
	gs.RoomCodeToState[roomCode] = RoomState{
		RoomPlayerList: make(map[*Player]bool),
		State:          ROOM_WAITING,
	}
	log.Printf("Created a new room with code %s...", roomCode)
	response := RoomCodeResponse{
		Status:   "success",
		PlayerId: playerId,
		RoomCode: roomCode,
		State:    gs.RoomCodeToState[roomCode].State,
		Host:     true,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		SendJSONError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET: checks whether a room with a certain room code exists
func (gs *GameServer) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomCode := r.PathValue("roomCode")
	playerId := generateCode(4)
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.PlayerCodes[playerId] = true
	if _, ok := gs.RoomCodeToState[roomCode]; ok {
		response := RoomCodeResponse{
			Status:   "success",
			PlayerId: playerId,
			RoomCode: roomCode,
			State:    gs.RoomCodeToState[roomCode].State,
			Host:     false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			SendJSONError(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		SendJSONError(w, "Room not found", http.StatusBadRequest)
		return
	}
}

// test endpoints
func ApiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API root!")
}

func ApiTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are at the API test!")
}
