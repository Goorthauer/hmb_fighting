package main

import (
	"log"
	"net/http"

	swagger "github.com/swaggo/http-swagger"
)

func main() {
	http.Handle("/ws", enableCORS(http.HandlerFunc(handleWebSocket)))
	http.Handle("/register", enableCORS(http.HandlerFunc(handleRegister)))
	http.Handle("/refresh", enableCORS(http.HandlerFunc(handleRefresh)))
	http.Handle("/check-client", enableCORS(http.HandlerFunc(handleCheckClient)))
	http.Handle("/create-room", enableCORS(http.HandlerFunc(handleCreateRoom)))
	http.Handle("/restart", enableCORS(http.HandlerFunc(handleRestart)))
	http.Handle("/swagger/", swagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
