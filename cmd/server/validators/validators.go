package validators

import (
	"errors"
	"hmb_fighting/cmd/server/types"
)

func ValidateRegisterInput(user types.User) error {
	if user.Name == "" || user.Email == "" {
		return errors.New("Name and email are required")
	}
	return nil
}

func ValidateRefreshInput(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("Refresh token is required")
	}
	return nil
}

func ValidateCheckClientInput(clientID, accessToken string) error {
	if clientID == "" || accessToken == "" {
		return errors.New("ClientID and accessToken are required")
	}
	return nil
}

func ValidateAccessToken(accessToken string) error {
	if accessToken == "" {
		return errors.New("AccessToken is required")
	}
	return nil
}

func ValidateRestartInput(accessToken, roomID string) error {
	if accessToken == "" || roomID == "" {
		return errors.New("AccessToken and RoomID are required")
	}
	return nil
}

func ValidateWebSocketInput(room, accessToken string) error {
	if room == "" || room == "null" {
		return errors.New("Room parameter is required")
	}
	if accessToken == "" {
		return errors.New("AccessToken is required")
	}
	return nil
}

func ValidateSetTeamInput(roomID string, realTeamID int, accessToken string) error {
	if roomID == "" || accessToken == "" {
		return errors.New("RoomID and AccessToken are required")
	}
	if realTeamID < 0 {
		return errors.New("Invalid RealTeamID")
	}
	return nil
}

func ValidateCheckTeamsInput(roomID, accessToken string) error {
	if roomID == "" || accessToken == "" {
		return errors.New("RoomID and AccessToken are required")
	}
	return nil
}

func ValidateLeaveRoomInput(accessToken, roomID string) error {
	if accessToken == "" || roomID == "" {
		return errors.New("AccessToken and RoomID are required")
	}
	return nil
}
