package dtos

import "hmb_fighting/cmd/server/types"

type RegisterUserResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ClientID     string `json:"clientID"`
}

type SelectTeamResp struct {
	AvailableTeams map[int]types.TeamConfig  `json:"availableTeams"`
	Characters     map[int][]types.Character `json:"characters"`
}
