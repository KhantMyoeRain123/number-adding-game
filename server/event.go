package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
)

type EventHandler func(event Event, p *Player)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	START = "start"

	//output
	QUESTION = "question"
)

// input
type StartPayload struct {
}

// output
type QuestionPayload struct {
	Question [QUESTION_LENGTH]int `json:"question"`
}

func Broadcast(p *Player, event Event) {
	gs := p.GameServer
	roomPlayerList := gs.RoomCodeToState[p.RoomCode].RoomPlayerList

	for _, player := range roomPlayerList {
		player.Egress <- event
	}
}

func StartHandler(event Event, p *Player) {
	if !p.Host {
		log.Println("A non-host player cannot start a game!")
		return
	}
	p.GameServer.mu.Lock()
	defer p.GameServer.mu.Unlock()
	//generates a question and the answer for the room the host is in
	var sum int
	p.GameServer.CurrentQuestionsMap[p.RoomCode] = [QUESTION_LENGTH]int{}
	question := p.GameServer.CurrentQuestionsMap[p.RoomCode]
	for i := 0; i < QUESTION_LENGTH; i++ {
		num := rand.IntN(10)
		question[i] = num
		sum += num
	}
	log.Println("Sum: ", sum)
	p.GameServer.CurrentQuestionsMap[p.RoomCode] = question
	p.GameServer.CurrentAnswersMap[p.RoomCode] = sum
	//broadcasts the QUESTION event with the question in it to all players in the host's room
	questionPayload := QuestionPayload{
		Question: p.GameServer.CurrentQuestionsMap[p.RoomCode],
	}
	questionPayloadBytes, err := json.Marshal(questionPayload)

	if err != nil {
		log.Println("Could not encode question payload into json...")
	}
	questionEvent := Event{
		Type:    QUESTION,
		Payload: questionPayloadBytes,
	}
	Broadcast(p, questionEvent)
}
