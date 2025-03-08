package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/ws", enableCORS(http.HandlerFunc(handleWebSocket)))
	http.Handle("/register", enableCORS(http.HandlerFunc(handleRegister)))
	http.Handle("/check-client", enableCORS(http.HandlerFunc(handleCheckClient)))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))
	http.Handle("/", http.FileServer(http.Dir("..")))

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
