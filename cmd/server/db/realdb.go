package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hmb_fighting/cmd/server/types"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return &PostgresDatabase{db: db}, nil
}

func (p *PostgresDatabase) GetWeapons() (map[string]types.Weapon, error) {
	rows, err := p.db.Query(`
		SELECT name, display_name, range, is_two_handed, image_url, attack_bonus, grapple_bonus
		FROM weapons`)
	if err != nil {
		return nil, fmt.Errorf("failed to query weapons: %v", err)
	}
	defer rows.Close()

	weapons := make(map[string]types.Weapon)
	for rows.Next() {
		var w types.Weapon
		if err := rows.Scan(&w.Name, &w.DisplayName, &w.Range, &w.IsTwoHanded, &w.ImageURL, &w.AttackBonus, &w.GrappleBonus); err != nil {
			return nil, fmt.Errorf("failed to scan weapon: %v", err)
		}
		weapons[w.Name] = w
	}
	return weapons, nil
}

func (p *PostgresDatabase) GetShields() (map[string]types.Shield, error) {
	rows, err := p.db.Query(`
		SELECT name, display_name, defense_bonus, image_url, attack_bonus, grapple_bonus
		FROM shields`)
	if err != nil {
		return nil, fmt.Errorf("failed to query shields: %v", err)
	}
	defer rows.Close()

	shields := make(map[string]types.Shield)
	for rows.Next() {
		var s types.Shield
		if err := rows.Scan(&s.Name, &s.DisplayName, &s.DefenseBonus, &s.ImageURL, &s.AttackBonus, &s.GrappleBonus); err != nil {
			return nil, fmt.Errorf("failed to scan shield: %v", err)
		}
		shields[s.Name] = s
	}
	return shields, nil
}

func (p *PostgresDatabase) GetTeamsConfig() (map[int]types.TeamConfig, error) {
	rows, err := p.db.Query(`
		SELECT id, name, icon_url, description
		FROM teams_config`)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams_config: %v", err)
	}
	defer rows.Close()

	teams := make(map[int]types.TeamConfig)
	for rows.Next() {
		var t types.TeamConfig
		if err := rows.Scan(&t.ID, &t.Name, &t.IconURL, &t.Description); err != nil {
			return nil, fmt.Errorf("failed to scan team_config: %v", err)
		}
		teams[t.ID] = t
	}
	return teams, nil
}

func (p *PostgresDatabase) GetCharacters() ([]types.Character, error) {
	rows, err := p.db.Query(`
		SELECT id, name, team_id, role_id, count_of_ability, image_url, is_active,
		       weapon, shield, is_titan_armour, height, weight, hp, stamina, initiative,
		       wrestling, attack, defense, attack_min, attack_max
		FROM characters`)
	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %v", err)
	}
	defer rows.Close()

	var characters []types.Character
	for rows.Next() {
		var c types.Character
		if err := rows.Scan(&c.ID, &c.Name, &c.TeamID, &c.RoleID, &c.CountOfAbility, &c.ImageURL, &c.IsActive,
			&c.Weapon, &c.Shield, &c.IsTitanArmour, &c.Height, &c.Weight, &c.HP, &c.Stamina, &c.Initiative,
			&c.Wrestling, &c.Attack, &c.Defense, &c.AttackMin, &c.AttackMax); err != nil {
			return nil, fmt.Errorf("failed to scan character: %v", err)
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func (p *PostgresDatabase) GetAbilities() (map[string]types.Ability, error) {
	rows, err := p.db.Query(`
		SELECT name, display_name, type, description, range, image_url
		FROM abilities`)
	if err != nil {
		return nil, fmt.Errorf("failed to query abilities: %v", err)
	}
	defer rows.Close()

	abilities := make(map[string]types.Ability)
	for rows.Next() {
		var a types.Ability
		if err := rows.Scan(&a.Name, &a.DisplayName, &a.Type, &a.Description, &a.Range, &a.ImageURL); err != nil {
			return nil, fmt.Errorf("failed to scan ability: %v", err)
		}
		abilities[a.Name] = a
	}
	return abilities, nil
}

func (p *PostgresDatabase) GetRoleConfig() (map[string]types.Role, error) {
	rows, err := p.db.Query(`
		SELECT id, name
		FROM roles`)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %v", err)
	}
	defer rows.Close()

	roles := make(map[string]types.Role)
	for rows.Next() {
		var r types.Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %v", err)
		}
		roles[r.ID] = r
	}
	return roles, nil
}

func (p *PostgresDatabase) SetUser(refreshToken string, user types.User) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = $2, email = $3, password = $4`,
		user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert/update user: %v", err)
	}

	if refreshToken != "" {
		_, err = tx.Exec(`
			INSERT INTO refresh_tokens (token, user_id)
			VALUES ($1, $2)
			ON CONFLICT (token) DO UPDATE
			SET user_id = $2`,
			refreshToken, user.ID)
		if err != nil {
			return fmt.Errorf("failed to insert/update refresh token: %v", err)
		}
	}

	return tx.Commit()
}

func (p *PostgresDatabase) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	err := p.db.QueryRow(`
		SELECT id, name, email, password
		FROM users
		WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return types.User{}, nil
	}
	if err != nil {
		return types.User{}, fmt.Errorf("failed to query user by email: %v", err)
	}
	return user, nil
}

func (p *PostgresDatabase) GetUserByRefresh(token string) (types.User, error) {
	var user types.User
	err := p.db.QueryRow(`
		SELECT u.id, u.name, u.email, u.password
		FROM users u
		JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token = $1`, token).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return types.User{}, nil
	}
	if err != nil {
		return types.User{}, fmt.Errorf("failed to query user by refresh token: %v", err)
	}
	return user, nil
}

func (p *PostgresDatabase) GetRoom(roomID string) (*types.Game, error) {
	var game types.Game
	var boardJSON []byte
	err := p.db.QueryRow(`
		SELECT game_session_id, current_turn, phase, board, winner
		FROM rooms
		WHERE game_session_id = $1`, roomID).Scan(&game.GameSessionId, &game.CurrentTurn, &game.Phase, &boardJSON, &game.Winner)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query room: %v", err)
	}

	if err := json.Unmarshal(boardJSON, &game.Board); err != nil {
		return nil, fmt.Errorf("failed to unmarshal board: %v", err)
	}

	// Инициализация полей
	game.Teams = make(map[int]types.Team)
	game.Players = make(map[int]string)

	// Получение команд
	rows, err := p.db.Query(`
		SELECT team_id, characters
		FROM room_teams
		WHERE game_session_id = $1`, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query room teams: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var teamID int
		var charactersJSON []byte
		if err := rows.Scan(&teamID, &charactersJSON); err != nil {
			return nil, fmt.Errorf("failed to scan room team: %v", err)
		}
		var team types.Team
		if err := json.Unmarshal(charactersJSON, &team.Characters); err != nil {
			return nil, fmt.Errorf("failed to unmarshal characters: %v", err)
		}
		game.Teams[teamID] = team
	}

	// Получение игроков
	rows, err = p.db.Query(`
		SELECT team_id, client_id
		FROM room_players
		WHERE game_session_id = $1`, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query room players: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var teamID int
		var clientID string
		if err := rows.Scan(&teamID, &clientID); err != nil {
			return nil, fmt.Errorf("failed to scan room player: %v", err)
		}
		game.Players[teamID] = clientID
	}

	// Загружаем конфигурации (предполагаем, что они уже есть в памяти или кэше)
	game.WeaponsConfig, _ = p.GetWeapons()
	game.ShieldsConfig, _ = p.GetShields()
	game.AbilitiesConfig, _ = p.GetAbilities()
	game.RoleConfig, _ = p.GetRoleConfig()
	game.TeamsConfig, _ = p.GetTeamsConfig()

	return &game, nil
}

func (p *PostgresDatabase) SetRoom(game *types.Game) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return fmt.Errorf("failed to marshal board: %v", err)
	}

	_, err = tx.Exec(`
		INSERT INTO rooms (game_session_id, current_turn, phase, board, winner)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (game_session_id) DO UPDATE
		SET current_turn = $2, phase = $3, board = $4, winner = $5`,
		game.GameSessionId, game.CurrentTurn, game.Phase, boardJSON, game.Winner)
	if err != nil {
		return fmt.Errorf("failed to insert/update room: %v", err)
	}

	// Очистка старых данных
	_, err = tx.Exec(`DELETE FROM room_teams WHERE game_session_id = $1`, game.GameSessionId)
	if err != nil {
		return fmt.Errorf("failed to delete old room teams: %v", err)
	}
	_, err = tx.Exec(`DELETE FROM room_players WHERE game_session_id = $1`, game.GameSessionId)
	if err != nil {
		return fmt.Errorf("failed to delete old room players: %v", err)
	}

	// Сохранение команд
	for teamID, team := range game.Teams {
		charactersJSON, err := json.Marshal(team.Characters)
		if err != nil {
			return fmt.Errorf("failed to marshal characters for team %d: %v", teamID, err)
		}
		_, err = tx.Exec(`
			INSERT INTO room_teams (game_session_id, team_id, characters)
			VALUES ($1, $2, $3)`,
			game.GameSessionId, teamID, charactersJSON)
		if err != nil {
			return fmt.Errorf("failed to insert room team %d: %v", teamID, err)
		}
	}

	// Сохранение игроков
	for teamID, clientID := range game.Players {
		_, err = tx.Exec(`
			INSERT INTO room_players (game_session_id, team_id, client_id)
			VALUES ($1, $2, $3)`,
			game.GameSessionId, teamID, clientID)
		if err != nil {
			return fmt.Errorf("failed to insert room player for team %d: %v", teamID, err)
		}
	}

	return tx.Commit()
}
