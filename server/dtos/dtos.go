package dtos

import (
	"hmb_fighting/server/entities"
)

type RegisterUserResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ClientID     string `json:"clientID"`
}

type SelectTeamResp struct {
	AvailableTeams map[int]entities.TeamConfig  `json:"availableTeams"`
	Characters     map[int][]entities.Character `json:"characters"`
}
