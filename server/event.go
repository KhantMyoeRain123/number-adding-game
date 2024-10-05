package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"sort"
	"strconv"
	"time"
)

type EventHandler func(event Event, p *Player)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	//input
	START  = "start"
	ANSWER = "answer"

	//output
	QUESTION  = "question"
	CORRECT   = "correct"
	INCORRECT = "incorrect"
	RESULT    = "result"
)

// input
type StartPayload struct {
}
type AnswerPayload struct {
	Answer int `json:"answer"`
}

// output
type QuestionPayload struct {
	Question [QUESTION_LENGTH]int `json:"question"`
}
type CorrectPayload struct {
	Points int `json:"points"`
}

type IncorrectPayload struct {
	Points       int `json:"points"`
	ActualAnswer int `json:"actual_answer"`
}

type ResultPayload struct {
	Standing int `json:"standing"`
	Points   int `json:"points"`
}

func Broadcast(p *Player, event Event) {
	gs := p.GameServer
	roomPlayerList := gs.RoomCodeToState[p.RoomCode].RoomPlayerList

	for _, player := range roomPlayerList {
		player.Egress <- event
	}
}

func SendQuestionEvent(p *Player, roomState *RoomState) {
	var sum int
	currentQuestion := roomState.CurrentQuestion
	for i := 0; i < QUESTION_LENGTH; i++ {
		num := rand.IntN(10)
		currentQuestion[i] = num
		sum += num
	}
	log.Println("Sum: ", sum)
	roomState.CurrentQuestion = currentQuestion
	roomState.CurrentAnswer = sum

	log.Println("Current Question: ", roomState.CurrentQuestion)
	log.Println("Current Answer: ", roomState.CurrentAnswer)
	//broadcasts the QUESTION event with the question in it to all players in the host's room
	questionPayload := QuestionPayload{
		Question: currentQuestion,
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

func StartHandler(event Event, p *Player) {
	if !p.Host {
		log.Println("A non-host player cannot start a game!")
		return
	}
	gs := p.GameServer
	roomState := gs.RoomCodeToState[p.RoomCode]
	/*gs.mu.Lock()
	defer gs.mu.Unlock()*/ //should be ok to not lock
	//set room state to ROOM_STARTED
	roomState.State = ROOM_STARTED
	roomState.RoundNumber = 1
	//generates a question and the answer for the room the host is in
	SendQuestionEvent(p, roomState)
}

func AnswerHandler(event Event, p *Player) {
	gs := p.GameServer
	roomState := gs.RoomCodeToState[p.RoomCode]

	if roomState.State != ROOM_STARTED {
		log.Println("Room has not started...")
		return
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()

	roomState.NumberSubmitted++

	var answerPayload AnswerPayload

	err := json.Unmarshal(event.Payload, &answerPayload)

	if err != nil {
		log.Println("Could not decode bytes into answer payload...")
	}
	if answerPayload.Answer == roomState.CurrentAnswer {
		//send CORRECT event
		pointsReceived := (len(roomState.RoomPlayerList) - (roomState.NumberSubmitted - 1)) * 100
		p.Points += pointsReceived
		correctPayload := CorrectPayload{
			Points: p.Points,
		}
		correctPayloadBytes, err := json.Marshal(correctPayload)

		if err != nil {
			log.Println("Could not encode correct payload into json...")
		}

		correctEvent := Event{
			Type:    CORRECT,
			Payload: correctPayloadBytes,
		}
		p.Egress <- correctEvent

	} else {
		//send INCORRECT event
		pointsDeducted := (roomState.NumberSubmitted) * 50
		p.Points -= pointsDeducted
		incorrectPayload := IncorrectPayload{
			Points:       p.Points,
			ActualAnswer: roomState.CurrentAnswer,
		}
		incorrectPayloadBytes, err := json.Marshal(incorrectPayload)

		if err != nil {
			log.Println("Could not encode incorrect payload into json...")
		}

		incorrectEvent := Event{
			Type:    INCORRECT,
			Payload: incorrectPayloadBytes,
		}
		p.Egress <- incorrectEvent

	}
	if roomState.NumberSubmitted == len(roomState.RoomPlayerList) {
		//sleep for 5 seconds to wait for last player to get feedback
		gs.mu.Unlock() //unlock here because we are just sleeping
		start := time.Now()
		time.Sleep(5 * time.Second)
		duration := time.Since(start)
		log.Printf("Slept for: %v\n", duration)

		gs.mu.Lock() //lock again because we are accessing shared resource again
		roomState.NumberSubmitted = 0
		roomState.RoundNumber++

		if roomState.RoundNumber > MAX_NUMBER_OF_ROUNDS {
			roomState.State = ROOM_ENDED
			//extract the players from roomPlayerList
			roomPlayerList := roomState.RoomPlayerList
			players := make([]*Player, len(roomPlayerList))
			i := 0
			for _, player := range roomPlayerList {
				players[i] = player
				i++
			}
			//sort the players by points
			sort.Slice(players, func(i, j int) bool {
				return players[i].Points >= players[j].Points
			})
			//return appropriate results to each player
			previousPoints := -1
			standing := 1
			for index, player := range players {
				if index != 0 {
					if previousPoints > player.Points {
						standing++
					}
				}
				log.Println("Standing: " + strconv.Itoa(standing))
				log.Println("PlayerId: " + player.PlayerId)
				log.Println("Points: " + strconv.Itoa(player.Points))

				resultPayload := ResultPayload{
					Standing: standing,
					Points:   player.Points,
				}

				resultPayloadBytes, err := json.Marshal(resultPayload)

				if err != nil {
					log.Println("Could not encode result payload into json...")
				}
				resultEvent := Event{
					Type:    RESULT,
					Payload: resultPayloadBytes,
				}
				player.Egress <- resultEvent

				previousPoints = player.Points
			}

		} else {
			//send new QUESTION event
			SendQuestionEvent(p, roomState)
		}
	}

}
