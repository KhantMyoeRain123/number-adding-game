package main

import (
	"fmt"
	"log"
	"net/http"
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
	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
