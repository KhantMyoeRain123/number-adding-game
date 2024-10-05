package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	PlayerId   string
	PlayerName string
	Host       bool
	RoomCode   string
	Points     int
	Connection *websocket.Conn
	GameServer *GameServer
	Egress     chan Event
}

func NewPlayer(playerId string, playerName string, host bool, roomCode string, conn *websocket.Conn, gameServer *GameServer) *Player {
	return &Player{
		PlayerId:   playerId,
		PlayerName: playerName,
		Host:       host,
		RoomCode:   roomCode,
		Points:     0,
		Connection: conn,
		GameServer: gameServer,
		Egress:     make(chan Event),
	}
}

func (p *Player) ReadMessages() {
	defer func() {
		p.GameServer.mu.Lock()
		defer p.GameServer.mu.Unlock()
		p.GameServer.RemovePlayer(p)
	}()
	for {
		_, eventBytes, err := p.Connection.ReadMessage()

		if err != nil {

			if websocket.IsCloseError(err) {
				log.Println("Client sent a close message, closing connection gracefully.")
			}
			break
		}
		var event Event

		err = json.Unmarshal(eventBytes, &event)

		if err != nil {
			log.Println("Error unmarshalling read event bytes...")
		}
		routeEvent(event, p)
	}
}

func routeEvent(event Event, p *Player) {
	p.GameServer.Handlers[event.Type](event, p)
}

func (p *Player) WriteMessages() {
	for {
		select {
		case event, ok := <-p.Egress:
			message, err := json.Marshal(event)
			if err != nil {
				log.Println("Error marshalling event for writing...")
			}
			if !ok {
				if err = p.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			if err := p.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}
}
