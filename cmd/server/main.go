package main

import (
	"log"
	"net/http"

	swagger "github.com/swaggo/http-swagger"
)

func main() {
	db := NewMockDatabase() // или NewRealDatabase(), в зависимости от реализации
	handler := NewHandler(db)

	http.Handle("/ws", enableCORS(http.HandlerFunc(handler.handleWebSocket)))
	http.Handle("/register", enableCORS(http.HandlerFunc(handler.handleRegister)))
	http.Handle("/refresh", enableCORS(http.HandlerFunc(handler.handleRefresh)))
	http.Handle("/check-client", enableCORS(http.HandlerFunc(handler.handleCheckClient)))
	http.Handle("/create-room", enableCORS(http.HandlerFunc(handler.handleCreateRoom)))
	http.Handle("/restart", enableCORS(http.HandlerFunc(handler.handleRestart)))
	http.Handle("/swagger/", swagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
