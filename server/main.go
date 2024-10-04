package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	mux := http.NewServeMux()
	setupAPI(mux)
}

func setupAPI(mux *http.ServeMux) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	gameServer := NewGameServer()
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	defer func() {
		for playerId, player := range gameServer.PlayerList {
			log.Println("Closing connection for " + playerId)
			player.Connection.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "Closing connection for "+playerId),
			)
		}
	}()

	//register the handlers to each route
	mux.HandleFunc("/api", ApiRoot)
	mux.HandleFunc("/api/test", ApiTest)
	mux.HandleFunc("/api/ws", gameServer.ServeWS)
	//room routes
	mux.HandleFunc("POST /api/room/make", gameServer.MakeRoom)
	mux.HandleFunc("POST /api/room/join/{roomCode}", gameServer.JoinRoom)

	go func() {
		log.Println("Server is running on http://localhost:8080")
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-sigChan
}
