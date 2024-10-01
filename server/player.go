package main

import "github.com/gorilla/websocket"

type Player struct {
	Connection *websocket.Conn
	GameServer *GameServer
}

func NewPlayer(conn *websocket.Conn, gameServer *GameServer) *Player {
	return &Player{
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
