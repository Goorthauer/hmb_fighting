-- +goose Up
CREATE TABLE users (
                       id VARCHAR(50) PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL
);

CREATE TABLE refresh_tokens (
                                token VARCHAR(255) PRIMARY KEY,
                                user_id VARCHAR(50) NOT NULL,
                                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE TABLE weapons (
                         name VARCHAR(50) PRIMARY KEY,
                         display_name VARCHAR(100) NOT NULL,
                         range INT NOT NULL,
                         is_two_handed BOOLEAN NOT NULL,
                         image_url VARCHAR(255),
                         attack_bonus INT NOT NULL,
                         grapple_bonus INT NOT NULL
);

CREATE TABLE shields (
                         name VARCHAR(50) PRIMARY KEY,
                         display_name VARCHAR(100) NOT NULL,
                         defense_bonus INT NOT NULL,
                         image_url VARCHAR(255),
                         attack_bonus INT NOT NULL,
                         grapple_bonus INT NOT NULL
);

CREATE TABLE teams (
                       id INT PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       icon_url VARCHAR(255),
                       description TEXT
);

CREATE TABLE roles (
                       id VARCHAR(10) PRIMARY KEY,
                       name VARCHAR(50) NOT NULL
);

CREATE TABLE characters (
                            id INT PRIMARY KEY,
                            name VARCHAR(100) NOT NULL,
                            team_id INT NOT NULL,
                            role_id VARCHAR(10) NOT NULL,
                            count_of_ability INT NOT NULL,
                            image_url VARCHAR(255),
                            is_active BOOLEAN NOT NULL DEFAULT TRUE,
                            weapon VARCHAR(50),
                            shield VARCHAR(50),
                            is_titan_armour BOOLEAN NOT NULL DEFAULT FALSE,
                            height INT NOT NULL,
                            weight INT NOT NULL,
                            hp INT NOT NULL,
                            stamina INT NOT NULL,
                            initiative INT NOT NULL,
                            wrestling INT NOT NULL,
                            attack INT NOT NULL,
                            defense INT NOT NULL,
                            attack_min INT NOT NULL,
                            attack_max INT NOT NULL,
                            FOREIGN KEY (team_id) REFERENCES teams(id),
                            FOREIGN KEY (role_id) REFERENCES roles(id),
                            FOREIGN KEY (weapon) REFERENCES weapons(name),
                            FOREIGN KEY (shield) REFERENCES shields(name)
);

CREATE INDEX idx_characters_team_id ON characters(team_id);
CREATE INDEX idx_characters_role_id ON characters(role_id);

CREATE TABLE abilities (
                           name VARCHAR(50) PRIMARY KEY,
                           display_name VARCHAR(100) NOT NULL,
                           type VARCHAR(50) NOT NULL,
                           description TEXT,
                           range INT NOT NULL,
                           image_url VARCHAR(255)
);

CREATE TABLE rooms (
                       game_session_id VARCHAR(50) PRIMARY KEY,
                       current_turn INT NOT NULL DEFAULT 0,
                       phase VARCHAR(50) NOT NULL,
                       board JSONB NOT NULL,
                       winner INT NOT NULL DEFAULT -1,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE room_teams (
                            game_session_id VARCHAR(50),
                            team_id INT,
                            characters JSONB NOT NULL,
                            PRIMARY KEY (game_session_id, team_id),
                            FOREIGN KEY (game_session_id) REFERENCES rooms(game_session_id) ON DELETE CASCADE,
                            FOREIGN KEY (team_id) REFERENCES teams(id)
);

CREATE TABLE room_players (
                              game_session_id VARCHAR(50),
                              team_id INT,
                              client_id VARCHAR(50),
                              PRIMARY KEY (game_session_id, team_id),
                              FOREIGN KEY (game_session_id) REFERENCES rooms(game_session_id) ON DELETE CASCADE
);

CREATE INDEX idx_room_teams_game_session_id ON room_teams(game_session_id);
CREATE INDEX idx_room_players_game_session_id ON room_players(game_session_id);

-- +goose Down
DROP TABLE room_players;
DROP TABLE room_teams;
DROP TABLE rooms;
DROP TABLE abilities;
DROP TABLE characters;
DROP TABLE roles;
DROP TABLE teams;
DROP TABLE shields;
DROP TABLE weapons;
DROP TABLE refresh_tokens;
DROP TABLE users;