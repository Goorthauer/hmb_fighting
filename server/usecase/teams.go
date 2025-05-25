package usecase

import (
	"fmt"
	"hmb_fighting/server/dtos"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/jwt"
	"hmb_fighting/server/types"
)

func (u *Usecase) SelectTeam() (*dtos.SelectTeamResp, error) {
	teams, err := u.db.GetTeams()
	if err != nil {
		return nil, fmt.Errorf("Teams not found: %v", err)
	}
	characters, err := u.db.GetCharacters()
	if err != nil {
		return nil, fmt.Errorf("Characters not found: %v", err)
	}

	outChars := make(map[int][]entities.Character)
	for _, char := range characters {
		if !char.IsActive {
			continue
		}
		outChars[char.TeamID] = append(outChars[char.TeamID], char)
	}

	return &dtos.SelectTeamResp{
		AvailableTeams: teams,
		Characters:     outChars,
	}, nil
}

func (u *Usecase) SetTeam(roomID string, realTeamID int, accessToken string) error {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return fmt.Errorf("Room not found")
	}

	teams, err := u.db.GetTeams()
	if err != nil {
		return fmt.Errorf("Teams not found: %v", err)
	}
	if _, ok := teams[realTeamID]; !ok {
		return fmt.Errorf("Team not found")
	}
	characters, err := u.db.GetCharacters()
	if err != nil {
		return fmt.Errorf("Characters not found: %v", err)
	}

	characterTeam := make([]entities.Character, 0)
	for _, char := range characters {
		if !char.IsActive || char.TeamID != realTeamID {
			continue
		}
		char.PrepareToFight(game.AbilitiesConfig)
		characterTeam = append(characterTeam, char)
	}
	if len(characterTeam) == 0 {
		return fmt.Errorf("Team has no active characters")
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if game.TeamsConfig == nil {
		game.TeamsConfig = make(map[int]entities.TeamConfig)
	}
	if game.Teams == nil {
		game.Teams = make(map[int]entities.Team)
	}
	for _, selectedTeam := range game.TeamsConfig {
		if realTeamID == selectedTeam.ID {
			return fmt.Errorf("Team already selected")
		}
	}
	teamID := -1
	if game.Players[0] == claims.ClientID {
		teamID = 0
	} else if game.Players[1] == claims.ClientID {
		teamID = 1
	} else if len(game.Players) < 2 && claims.Role == types.UserRoleSpectator {
		teamID = 1
		game.Players[1] = claims.ClientID
	}
	if teamID == -1 {
		return fmt.Errorf("Invalid team ID")
	}
	for i := range characterTeam {
		characterTeam[i].TeamID = teamID
		if shield, ok := game.ShieldsConfig[characterTeam[i].Shield]; ok {
			characterTeam[i].Defense += shield.DefenseBonus
		}
	}
	game.TeamsConfig[teamID] = teams[realTeamID]
	game.Teams[teamID] = entities.Team{Characters: characterTeam}

	if len(game.Players) == 2 {
		if len(game.InitialOrder) == 0 {
			game.InitTurnOrder()
		}
		game.Phase = types.GamePhaseSetup
	}

	return u.db.SetRoom(game)
}

func (u *Usecase) CheckTeams(roomID, accessToken string) (bool, error) {
	_, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return false, fmt.Errorf("Invalid token: %v", err)
	}

	game, err := u.db.GetRoom(roomID)
	if err != nil || game == nil {
		return false, fmt.Errorf("Room not found")
	}

	return game.Phase == types.GamePhaseSetup, nil
}
