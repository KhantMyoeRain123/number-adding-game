package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func (r *http.Request) bool{
			origin:=r.Header.Get("Origin")

			return origin=="http://localhost:5173"			
		},
	}
)

func main() {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You are at the API root!")
	})
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You are at the API test!")
	})
	mux.HandleFunc("/api/ws", serveWS)

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")

	_, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
