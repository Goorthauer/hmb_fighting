package handlers

import (
	"encoding/json"
	"hmb_fighting/cmd/server/types"
	"hmb_fighting/cmd/server/usecase"
	"hmb_fighting/cmd/server/validators"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(uc *usecase.Usecase) *Handler { // Изменено с db.Database на *usecase.Usecase
	return &Handler{usecase: uc}
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var currentUser types.User
	if err := json.NewDecoder(r.Body).Decode(&currentUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateRegisterInput(currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.usecase.RegisterUser(currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Registered user: %s (%s) with ClientID: %s", currentUser.Name, currentUser.Email, response.ClientID)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := validators.ValidateLoginInput(req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.usecase.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateRefreshInput(req.RefreshToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.usecase.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleCheckClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID    string `json:"clientID"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateCheckClientInput(req.ClientID, req.AccessToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := h.usecase.CheckClient(req.ClientID, req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"valid":   valid,
	})
	log.Printf("Checked clientID: %s, valid: %v", req.ClientID, valid)
}

func (h *Handler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateAccessToken(req.AccessToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomID, err := h.usecase.CreateRoom(req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"roomID":  roomID,
	})
	log.Printf("Room %s created", roomID)
}

func (h *Handler) HandleRestart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"accessToken"`
		RoomID      string `json:"roomID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateRestartInput(req.AccessToken, req.RoomID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.usecase.RestartRoom(req.AccessToken, req.RoomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Room restarted",
	})
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	accessToken := r.URL.Query().Get("accessToken")

	log.Printf("WebSocket request received: room=%s, accessToken=%s", room, accessToken)

	if err := validators.ValidateWebSocketInput(room, accessToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	err = h.usecase.HandleWebSocket(conn, room, accessToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
}

func (h *Handler) HandleSelectTeam(w http.ResponseWriter, r *http.Request) {
	response, err := h.usecase.SelectTeam()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleSetTeam(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID      string `json:"roomID"`
		RealTeamID  int    `json:"realTeamID"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateSetTeamInput(req.RoomID, req.RealTeamID, req.AccessToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.usecase.SetTeam(req.RoomID, req.RealTeamID, req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Team assigned",
	})
}

func (h *Handler) HandleCheckTeams(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID      string `json:"roomID"`
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validators.ValidateCheckTeamsInput(req.RoomID, req.AccessToken); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allTeamsSelected, err := h.usecase.CheckTeams(req.RoomID, req.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		AllTeamsSelected bool `json:"allTeamsSelected"`
	}{
		AllTeamsSelected: allTeamsSelected,
	})
}

func (h *Handler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"accessToken"`
		RoomID      string `json:"roomID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateLeaveRoomInput(req.AccessToken, req.RoomID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.usecase.LeaveRoom(req.AccessToken, req.RoomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully left the room",
	})
}
