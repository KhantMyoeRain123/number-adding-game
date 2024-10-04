package main

import "github.com/gorilla/websocket"

type Player struct {
	PlayerId   string
	PlayerName string
	Host       bool
	Connection *websocket.Conn
	GameServer *GameServer
}

func NewPlayer(playerId string, playerName string, host bool, conn *websocket.Conn, gameServer *GameServer) *Player {
	return &Player{
		PlayerId:   playerId,
		PlayerName: playerName,
		Host:       host,
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
