package main

import "github.com/gorilla/websocket"

type Player struct {
	PlayerId   string
	PlayerName string
	Host       bool
	RoomCode   string
	Connection *websocket.Conn
	GameServer *GameServer
}

func NewPlayer(playerId string, playerName string, host bool, roomCode string, conn *websocket.Conn, gameServer *GameServer) *Player {
	return &Player{
		PlayerId:   playerId,
		PlayerName: playerName,
		Host:       host,
		RoomCode:   roomCode,
		Connection: conn,
		GameServer: gameServer,
	}
}

func (p *Player) ReadMessages() {
	return
}

func (p *Player) WriteMessages() {
	return
}
