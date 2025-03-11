package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hmb_fighting/cmd/server/types"
	"time"

	_ "github.com/lib/pq"
	redis "github.com/redis/go-redis/v9"
)

type PostgresDatabase struct {
	db      *sql.DB
	redis   *redis.Client
	baseTTL time.Duration
}

func NewPostgresDatabase(connStr string, redisAddr string) (*PostgresDatabase, error) {
	baseTTL := time.Hour * 24
	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr, // например, "localhost:6379"
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &PostgresDatabase{
		db:      db,
		redis:   redisClient,
		baseTTL: baseTTL,
	}, nil
}

// cacheGet пытается получить данные из Redis по ключу
func (p *PostgresDatabase) cacheGet(ctx context.Context, key string, dest interface{}) error {
	data, err := p.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil // Кеш пуст
	}
	if err != nil {
		return fmt.Errorf("failed to get from cache: %v", err)
	}
	return json.Unmarshal(data, dest)
}

// cacheSet сохраняет данные в Redis с заданным TTL
func (p *PostgresDatabase) cacheSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %v", err)
	}
	return p.redis.Set(ctx, key, data, ttl).Err()
}

// cacheDel удаляет ключ из кеша
func (p *PostgresDatabase) cacheDel(ctx context.Context, key string) error {
	return p.redis.Del(ctx, key).Err()
}

func (p *PostgresDatabase) GetWeapons() (map[string]types.Weapon, error) {
	ctx := context.Background()
	cacheKey := "weapons_"
	var weapons map[string]types.Weapon

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &weapons); err == nil && weapons != nil {
		return weapons, nil
	}

	// Если в кеше нет, идём в базу
	rows, err := p.db.Query(`
		SELECT name, display_name, range, is_two_handed, image_url, attack_bonus, grapple_bonus
		FROM weapons`)
	if err != nil {
		return nil, fmt.Errorf("failed to query weapons: %v", err)
	}
	defer rows.Close()

	weapons = make(map[string]types.Weapon)
	for rows.Next() {
		var w types.Weapon
		if err := rows.Scan(&w.Name, &w.DisplayName, &w.Range, &w.IsTwoHanded, &w.ImageURL, &w.AttackBonus, &w.GrappleBonus); err != nil {
			return nil, fmt.Errorf("failed to scan weapon: %v", err)
		}
		weapons[w.Name] = w
	}

	if err := p.cacheSet(ctx, cacheKey, weapons, p.baseTTL); err != nil {
		// Ошибка кеширования не критична, просто логируем (или игнорируем)
		fmt.Printf("failed to cache weapons: %v\n", err)
	}

	return weapons, nil
}

func (p *PostgresDatabase) GetShields() (map[string]types.Shield, error) {
	ctx := context.Background()
	cacheKey := "shields_"
	var shields map[string]types.Shield

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &shields); err == nil && shields != nil {
		return shields, nil
	}

	// Если кеша нет, идём в базу
	rows, err := p.db.Query(`
		SELECT name, display_name, defense_bonus, image_url, attack_bonus, grapple_bonus
		FROM shields`)
	if err != nil {
		return nil, fmt.Errorf("failed to query shields: %v", err)
	}
	defer rows.Close()

	shields = make(map[string]types.Shield)
	for rows.Next() {
		var s types.Shield
		if err := rows.Scan(&s.Name, &s.DisplayName, &s.DefenseBonus, &s.ImageURL, &s.AttackBonus, &s.GrappleBonus); err != nil {
			return nil, fmt.Errorf("failed to scan shield: %v", err)
		}
		shields[s.Name] = s
	}

	if err := p.cacheSet(ctx, cacheKey, shields, p.baseTTL); err != nil {
		fmt.Printf("failed to cache shields: %v\n", err)
	}

	return shields, nil
}

func (p *PostgresDatabase) GetTeams() (map[int]types.TeamConfig, error) {
	ctx := context.Background()
	cacheKey := "teams_"
	var teams map[int]types.TeamConfig

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &teams); err == nil && teams != nil {
		return teams, nil
	}

	// Если кеша нет, идём в базу
	rows, err := p.db.Query(`
		SELECT id, name, icon_url, description
		FROM teams`)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %v", err)
	}
	defer rows.Close()

	teams = make(map[int]types.TeamConfig)
	for rows.Next() {
		var t types.TeamConfig
		if err := rows.Scan(&t.ID, &t.Name, &t.IconURL, &t.Description); err != nil {
			return nil, fmt.Errorf("failed to scan team_config: %v", err)
		}
		teams[t.ID] = t
	}

	if err := p.cacheSet(ctx, cacheKey, teams, p.baseTTL); err != nil {
		fmt.Printf("failed to cache teams: %v\n", err)
	}

	return teams, nil
}

func (p *PostgresDatabase) GetCharacters() ([]types.Character, error) {
	ctx := context.Background()
	cacheKey := "characters_"
	var characters []types.Character

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &characters); err == nil && len(characters) > 0 {
		return characters, nil
	}

	// Если кеша нет, идём в базу
	rows, err := p.db.Query(`
		SELECT id, name, team_id, role_id, count_of_ability, image_url, is_active,
		       weapon, shield, is_titan_armour, height, weight, hp, stamina, initiative,
		       wrestling, attack, defense, attack_min, attack_max
		FROM characters`)
	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c types.Character
		if err := rows.Scan(&c.ID, &c.Name, &c.TeamID, &c.RoleID, &c.CountOfAbility, &c.ImageURL, &c.IsActive,
			&c.Weapon, &c.Shield, &c.IsTitanArmour, &c.Height, &c.Weight, &c.HP, &c.Stamina, &c.Initiative,
			&c.Wrestling, &c.Attack, &c.Defense, &c.AttackMin, &c.AttackMax); err != nil {
			return nil, fmt.Errorf("failed to scan character: %v", err)
		}
		characters = append(characters, c)
	}

	if err := p.cacheSet(ctx, cacheKey, characters, p.baseTTL); err != nil {
		fmt.Printf("failed to cache characters: %v\n", err)
	}

	return characters, nil
}

func (p *PostgresDatabase) GetAbilities() (map[string]types.Ability, error) {
	ctx := context.Background()
	cacheKey := "abilities_"
	var abilities map[string]types.Ability

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &abilities); err == nil && abilities != nil {
		return abilities, nil
	}

	// Если кеша нет, идём в базу
	rows, err := p.db.Query(`
		SELECT name, display_name, type, description, range, image_url
		FROM abilities`)
	if err != nil {
		return nil, fmt.Errorf("failed to query abilities: %v", err)
	}
	defer rows.Close()

	abilities = make(map[string]types.Ability)
	for rows.Next() {
		var a types.Ability
		if err := rows.Scan(&a.Name, &a.DisplayName, &a.Type, &a.Description, &a.Range, &a.ImageURL); err != nil {
			return nil, fmt.Errorf("failed to scan ability: %v", err)
		}
		abilities[a.Name] = a
	}

	if err := p.cacheSet(ctx, cacheKey, abilities, p.baseTTL); err != nil {
		fmt.Printf("failed to cache abilities: %v\n", err)
	}

	return abilities, nil
}

func (p *PostgresDatabase) GetRoleConfig() (map[string]types.Role, error) {
	ctx := context.Background()
	cacheKey := "role_config_"
	var roles map[string]types.Role

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &roles); err == nil && roles != nil {
		return roles, nil
	}

	// Если кеша нет, идём в базу
	rows, err := p.db.Query(`
		SELECT id, name
		FROM roles`)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %v", err)
	}
	defer rows.Close()

	roles = make(map[string]types.Role)
	for rows.Next() {
		var r types.Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %v", err)
		}
		roles[r.ID] = r
	}

	if err := p.cacheSet(ctx, cacheKey, roles, p.baseTTL); err != nil {
		fmt.Printf("failed to cache role_config: %v\n", err)
	}

	return roles, nil
}

func (p *PostgresDatabase) SetUser(refreshToken string, user types.User) error {
	ctx := context.Background()
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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Инвалидируем кеш для пользователя по email и refresh token
	emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
	if err := p.cacheDel(ctx, emailCacheKey); err != nil {
		fmt.Printf("failed to invalidate cache for user by email %s: %v\n", user.Email, err)
	}
	if refreshToken != "" {
		tokenCacheKey := fmt.Sprintf("user:refresh:%s", refreshToken)
		if err := p.cacheDel(ctx, tokenCacheKey); err != nil {
			fmt.Printf("failed to invalidate cache for user by refresh %s: %v\n", refreshToken, err)
		}
	}

	return nil
}

func (p *PostgresDatabase) GetUserByEmail(email string) (types.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:email:%s", email)
	var user types.User

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &user); err == nil && user.ID != "" {
		return user, nil
	}

	// Если кеша нет, идём в базу
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

	if err := p.cacheSet(ctx, cacheKey, user, p.baseTTL); err != nil {
		fmt.Printf("failed to cache user by email %s: %v\n", email, err)
	}

	return user, nil
}

func (p *PostgresDatabase) GetUserByRefresh(token string) (types.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:refresh:%s", token)
	var user types.User

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &user); err == nil && user.ID != "" {
		return user, nil
	}

	// Если кеша нет, идём в базу
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

	if err := p.cacheSet(ctx, cacheKey, user, p.baseTTL); err != nil {
		fmt.Printf("failed to cache user by refresh %s: %v\n", token, err)
	}

	return user, nil
}

func (p *PostgresDatabase) GetRoom(roomID string) (*types.Game, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("room:%s", roomID)
	var game types.Game

	// Проверяем кеш
	if err := p.cacheGet(ctx, cacheKey, &game); err == nil && game.GameSessionId != "" {
		return &game, nil
	}

	// Логика получения из базы (без изменений до конца метода)
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

	game.Teams = make(map[int]types.Team)
	game.Players = make(map[int]string)

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

	game.WeaponsConfig, _ = p.GetWeapons()
	game.ShieldsConfig, _ = p.GetShields()
	game.AbilitiesConfig, _ = p.GetAbilities()
	game.RoleConfig, _ = p.GetRoleConfig()
	game.TeamsConfig, _ = p.GetTeams()

	if err := p.cacheSet(ctx, cacheKey, game, 3*time.Minute); err != nil {
		fmt.Printf("failed to cache room %s: %v\n", roomID, err)
	}

	return &game, nil
}

func (p *PostgresDatabase) SetRoom(game *types.Game) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("room:%s", game.GameSessionId)

	// Логика записи в базу (без изменений)
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

	_, err = tx.Exec(`DELETE FROM room_teams WHERE game_session_id = $1`, game.GameSessionId)
	if err != nil {
		return fmt.Errorf("failed to delete old room teams: %v", err)
	}
	_, err = tx.Exec(`DELETE FROM room_players WHERE game_session_id = $1`, game.GameSessionId)
	if err != nil {
		return fmt.Errorf("failed to delete old room players: %v", err)
	}

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

	for teamID, clientID := range game.Players {
		_, err = tx.Exec(`
			INSERT INTO room_players (game_session_id, team_id, client_id)
			VALUES ($1, $2, $3)`,
			game.GameSessionId, teamID, clientID)
		if err != nil {
			return fmt.Errorf("failed to insert room player for team %d: %v", teamID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Инвалидируем кеш после обновления
	if err := p.cacheDel(ctx, cacheKey); err != nil {
		fmt.Printf("failed to invalidate cache for room %s: %v\n", game.GameSessionId, err)
	}

	return nil
}
