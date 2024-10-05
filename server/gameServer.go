package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	ROOM_CODE_LENGTH = 6
	ROOM_WAITING     = 0
	ROOM_STARTED     = 1
	ROOM_ENDED       = 2

	QUESTION_LENGTH      = 10
	MAX_NUMBER_OF_ROUNDS = 5
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
	RoomPlayerList  map[string]*Player
	State           int
	CurrentQuestion [QUESTION_LENGTH]int
	CurrentAnswer   int
	RoundNumber     int
	NumberSubmitted int
}

type GameServer struct {
	RoomCodeToState map[string]*RoomState //a mapping between room code and the users in the room
	PlayerList      map[string]*Player
	Handlers        map[string]EventHandler
	mu              sync.Mutex
}

func NewGameServer() *GameServer {

	return &GameServer{
		RoomCodeToState: make(map[string]*RoomState),
		PlayerList:      make(map[string]*Player),
		Handlers: map[string]EventHandler{
			START:  StartHandler,
			ANSWER: AnswerHandler,
		},
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

func (gs *GameServer) AddPlayer(player *Player) {
	//adds player to playerList
	gs.PlayerList[player.PlayerId] = player
	//adds the player to the room
	gs.RoomCodeToState[player.RoomCode].RoomPlayerList[player.PlayerId] = player
}

func (gs *GameServer) RemovePlayer(player *Player) {
	log.Println("Closing connection for " + player.PlayerId)
	if _, ok := gs.PlayerList[player.PlayerId]; ok {
		player.Connection.Close()
		delete(gs.PlayerList, player.PlayerId)

		roomPlayerList := gs.RoomCodeToState[player.RoomCode].RoomPlayerList
		//if player is host close connections of other  also delete the room
		if player.Host {
			for playerId, player := range roomPlayerList {
				log.Println("Host has exited...Closing connection for " + playerId)
				player.Connection.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Host has exited...Closing connection for "+playerId),
				)
			}
			delete(gs.RoomCodeToState, player.RoomCode)
			log.Println("Removed room " + player.RoomCode)
			log.Println(gs.PlayerList)
			log.Println(gs.RoomCodeToState)
		} else {
			delete(roomPlayerList, player.PlayerId)
			log.Println(gs.PlayerList)
			log.Println(roomPlayerList)
		}

	}

}

// upgrades to websocket
func (gs *GameServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	playerId := r.URL.Query().Get("playerId")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()
	player, ok := gs.PlayerList[playerId]
	if !ok {
		log.Println("Player does not exist...")
		conn.Close()
		return
	}
	player.Connection = conn
	player.Egress = make(chan Event)

	log.Println(gs.PlayerList)
	log.Println(gs.RoomCodeToState[player.RoomCode].RoomPlayerList)

	go player.ReadMessages()
	go player.WriteMessages()

}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
}

// REST handlers start here
// POST:make a room
func (gs *GameServer) MakeRoom(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("playerName")
	if playerName == "" {
		SendJSONError(w, "Need playerName to proceed", http.StatusBadRequest)
		return
	}
	roomCode := generateCode(ROOM_CODE_LENGTH)
	playerId := generateCode(4)
	gs.mu.Lock()
	defer gs.mu.Unlock()
	//create the new room
	gs.RoomCodeToState[roomCode] = &RoomState{
		RoomPlayerList:  make(map[string]*Player),
		State:           ROOM_WAITING,
		CurrentQuestion: [QUESTION_LENGTH]int{},
		CurrentAnswer:   -1,
		RoundNumber:     0,
		NumberSubmitted: 0,
	}
	log.Printf("Created a new room with code %s...", roomCode)
	player := NewPlayer(playerId, playerName, true, roomCode, nil, gs)

	gs.AddPlayer(player)

	response := RoomCodeResponse{
		Status:   "success",
		PlayerId: playerId,
		RoomCode: roomCode,
		State:    gs.RoomCodeToState[roomCode].State,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		SendJSONError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// POST: checks whether a room with a certain room code exists
func (gs *GameServer) JoinRoom(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("playerName")
	if playerName == "" {
		SendJSONError(w, "Need playerName to proceed", http.StatusBadRequest)
		return
	}
	roomCode := r.PathValue("roomCode")
	playerId := generateCode(4)
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if _, ok := gs.RoomCodeToState[roomCode]; ok {
		response := RoomCodeResponse{
			Status:   "success",
			PlayerId: playerId,
			RoomCode: roomCode,
			State:    gs.RoomCodeToState[roomCode].State,
		}

		player := NewPlayer(playerId, playerName, false, roomCode, nil, gs)
		gs.AddPlayer(player)

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
