package validators

import (
	"errors"
	"hmb_fighting/server/types"
	"regexp"
)

func ValidateRegisterInput(user types.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func ValidateLoginInput(email, password string) error {
	if password == "" || email == "" {
		return errors.New("password and email are required")
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
