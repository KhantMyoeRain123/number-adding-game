package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	setupAPI(mux)
}

func setupAPI(mux *http.ServeMux) {
	gameServer := NewGameServer()
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	//register the handlers to each route
	mux.HandleFunc("/api", ApiRoot)
	mux.HandleFunc("/api/test", ApiTest)
	mux.HandleFunc("/api/ws", gameServer.ServeWS)
	//room routes
	mux.HandleFunc("POST /api/room/make",gameServer.MakeRoom)
	mux.HandleFunc("POST /api/room/join",gameServer.JoinRoom)

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
