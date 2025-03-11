package main

import (
	"hmb_fighting/cmd/server/db"
	"hmb_fighting/cmd/server/handlers"
	"hmb_fighting/cmd/server/usecase"
	"log"
	"net/http"

	swagger "github.com/swaggo/http-swagger"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	database := db.NewMockDatabase()
	uc := usecase.NewUsecase(database)
	handler := handlers.NewHandler(uc)

	http.Handle("/ws", enableCORS(http.HandlerFunc(handler.HandleWebSocket)))
	http.Handle("/register", enableCORS(http.HandlerFunc(handler.HandleRegister)))
	http.Handle("/refresh", enableCORS(http.HandlerFunc(handler.HandleRefresh)))
	http.Handle("/check-client", enableCORS(http.HandlerFunc(handler.HandleCheckClient)))
	http.Handle("/create-room", enableCORS(http.HandlerFunc(handler.HandleCreateRoom)))
	http.Handle("/select-team", enableCORS(http.HandlerFunc(handler.HandleSelectTeam)))
	http.Handle("/check-teams", enableCORS(http.HandlerFunc(handler.HandleCheckTeams)))
	http.Handle("/set-team", enableCORS(http.HandlerFunc(handler.HandleSetTeam)))
	http.Handle("/restart", enableCORS(http.HandlerFunc(handler.HandleRestart)))
	http.Handle("/leave-room", enableCORS(http.HandlerFunc(handler.HandleLeaveRoom)))
	http.Handle("/swagger/", swagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
