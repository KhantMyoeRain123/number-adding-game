package main

import "encoding/json"

type EventHandler func(event Event, p *Player)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}



