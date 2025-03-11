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
                       board JSONB NOT NULL, -- [16][9]int как JSON
                       winner INT NOT NULL DEFAULT -1,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE room_teams (
                            game_session_id VARCHAR(50),
                            team_id INT,
                            characters JSONB NOT NULL, -- []Character как JSON
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



--------- inserts ---------

INSERT INTO weapons (name, display_name, range, is_two_handed, image_url, attack_bonus, grapple_bonus)
VALUES
    ('falchion', 'Фальшион', 1, FALSE, './static/weapons/default.png', 2, 0),
    ('axe', 'Топор', 1, FALSE, './static/weapons/default.png', 0, 8),
    ('two_handed_sword', 'Двуручный меч', 2, TRUE, './static/weapons/default.png', 2, 0),
    ('two_handed_halberd', 'Алебарда', 2, TRUE, './static/weapons/default.png', 2, 8),
    ('sword', 'Меч', 1, FALSE, './static/weapons/default.png', 0, 0)
    ON CONFLICT (name) DO NOTHING;

INSERT INTO shields (name, display_name, defense_bonus, image_url, attack_bonus, grapple_bonus)
VALUES
    ('buckler', 'Баклер', 1, './static/shields/default.png', 1, 1),
    ('shield', 'Тарч', 2, './static/shields/default.png', 1, 0),
    ('tower', 'Ростовой щит', 3, './static/shields/default.png', 0, -1)
    ON CONFLICT (name) DO NOTHING;

INSERT INTO abilities (name, display_name, type, description, range, image_url)
VALUES
    ('yama_arashi', 'Подхват', 'wrestle', 'Мощный бросок через бедро, использующий силу и импульс противника для стремительного падения.', 1, './static/abilities/default.jpg'),
    ('trip', 'Зацеп', 'wrestle', 'Точный удар по ноге, нарушающий баланс противника и заставляющий его упасть.', 1, './static/abilities/default.jpg'),
    ('hip_toss', 'Высед', 'wrestle', 'Бросок через бедро с использованием вращения и силы противника.', 1, './static/abilities/default.jpg'),
    ('front_sweep', 'Передняя подножка', 'wrestle', 'Быстрый удар по передней ноге, опрокидывающий противника вперёд.', 1, './static/abilities/default.jpg'),
    ('back_sweep', 'Задняя подножка', 'wrestle', 'Подсечка задней ноги, использующая вес противника для опрокидывания.', 1, './static/abilities/default.jpg'),
    ('outer_hook', 'Внешний зацеп', 'wrestle', 'Подсечка внешней стороны ноги, нарушающая равновесие противника.', 1, './static/abilities/default.jpg'),
    ('inner_hook', 'Внутренний зацеп', 'wrestle', 'Подсечка внутренней стороны ноги, выводящая противника из равновесия.', 1, './static/abilities/default.jpg'),
    ('shoulder_throw', 'Бросок через плечо', 'wrestle', 'Рывок противника через плечо с использованием его импульса.', 1, './static/abilities/default.jpg'),
    ('spinning_throw', 'Вращательный бросок', 'wrestle', 'Бросок противника через спину с использованием вращения.', 1, './static/abilities/default.jpg'),
    ('double_leg', 'Двойной захват ног', 'wrestle', 'Захват обеих ног противника с последующим опрокидыванием.', 1, './static/abilities/default.jpg'),
    ('single_leg', 'Одинарный захват ноги', 'wrestle', 'Захват одной ноги противника с последующим броском.', 1, './static/abilities/default.jpg'),
    ('overhook_throw', 'Бросок через захват', 'wrestle', 'Бросок противника через верхний захват руки.', 1, './static/abilities/default.jpg'),
    ('underhook_throw', 'Бросок через подхват', 'wrestle', 'Бросок противника через нижний захват руки.', 1, './static/abilities/default.jpg')
    ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (id, name)
VALUES
    ('0', 'Танк'),
    ('1', 'Убийца'),
    ('2', 'Боец'),
    ('3', 'Поддержка'),
    ('4', 'Борец')
    ON CONFLICT (id) DO NOTHING;


INSERT INTO teams (id, name, icon_url, description)
VALUES
    (1, 'Партизан Два', './static/teams/partizan_dva.png', 'Вторая команда партизан, стойкие и выносливые бойцы.'),
    (2, 'Юг', './static/teams/south.png', 'Команда южных земель, известная своей тактикой.'),
    (3, 'НСК', './static/teams/nsk.png', 'Новосибирские бойцы, сильные и решительные.'),
    (6, 'Старые Друзья', './static/teams/old_friends.png', 'Ветераны, проверенные временем.'),
    (7, 'Черная Земля', './static/teams/black_land.png', 'Таинственные воины темных земель.'),
    (11, 'Партизан', './static/teams/partizan.png', 'Первая команда партизан, мастера скрытности.'),
    (12, 'Злой дух Ямбуя', './static/teams/yambuya_spirit.png', 'Мистические воины севера.'),
    (16, 'Медвежья пядь', './static/teams/bear_span.png', 'Сильные, как медведи, бойцы.'),
    (4, 'Школа ИСБ Байард', './static/teams/isb_bayard.png', 'Ученики школы боевых искусств Байард.'),
    (5, 'НРБ', './static/teams/nrb.png', 'Непреклонные рыцари битвы.'),
    (8, 'Vivus Ferro', './static/teams/vivus_ferro.png', 'Живые клинки, мастера оружия.'),
    (9, 'Молодые львы', './static/teams/young_lions.png', 'Юные и амбициозные бойцы.'),
    (10, 'Ганза', './static/teams/hansa.png', 'Союз торговцев и воинов.'),
    (13, 'Мальтийский крест', './static/teams/maltese_cross.png', 'Рыцари ордена, верные клятве.'),
    (14, 'Урфин Джус', './static/teams/urfin_jus.png', 'Команда загадочных мастеров.'),
    (15, 'Высшая Школа ИСБ Санкт-Петербург', './static/teams/isb_spb.png', 'Элита школы боевых искусств СПб.'),
    (17, 'Байард', './static/teams/bayard.png', 'Рыцари чести и доблести.'),
    (18, 'RaubRitter', './static/teams/raubritter.png', 'Разбойные рыцари, мастера боя.'),
    (19, 'Межевой рыцарь', './static/teams/border_knight.png', 'Стражи границ и традиций.')
    ON CONFLICT (id) DO NOTHING;

INSERT INTO characters (id, name, team_id, role_id, count_of_ability, image_url, is_active, weapon, shield, is_titan_armour, height, weight, hp, stamina, initiative, wrestling, attack, defense, attack_min, attack_max)
VALUES
    (15, 'Тюляков Алексей', 1, '3', 4, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 177, 84, 100, 9, 11, 8, 12, 10, 12, 18),
    (14, 'Чуклов Григорий', 1, '2', 4, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 171, 77, 100, 8, 8, 10, 12, 12, 12, 18),
    (22, 'Корчагин Максим', 1, '1', 4, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 175, 80, 100, 9, 9, 10, 15, 8, 15, 22),
    (23, 'Шеварков Иван', 1, '2', 4, './static/characters/default.png', FALSE, 'falchion', 'shield', FALSE, 170, 78, 100, 6, 6, 8, 10, 10, 10, 15),
    (24, 'Остраумов Антон', 1, '4', 4, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 182, 88, 100, 8, 8, 14, 12, 8, 12, 18),
    (28, 'Наумов Александр', 1, '3', 4, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 181, 87, 100, 9, 11, 8, 10, 10, 10, 15),
    (40, 'Сотов Николай', 1, '2', 3, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 181, 89, 100, 6, 6, 8, 10, 10, 10, 15),
    (21, 'Ивлев Владимир', 1, '4', 4, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 167, 72, 100, 8, 8, 12, 12, 6, 12, 18),
    (75, 'Коврижин Дмитрий', 1, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 177, 82, 100, 6, 6, 8, 10, 8, 10, 15),
    (112, 'Янышевский Максим', 1, '2', 2, './static/characters/default.png', FALSE, 'falchion', 'buckler', FALSE, 175, 80, 100, 6, 6, 8, 10, 8, 10, 15),
    (113, 'Флерко Игорь', 1, '4', 2, './static/characters/default.png', FALSE, 'axe', NULL, FALSE, 172, 78, 100, 8, 8, 12, 12, 6, 12, 18),

    (51, 'Рябцев Сергей', 2, '1', 3, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 169, 75, 100, 9, 9, 10, 15, 8, 15, 22),
    (25, 'Свиридов Даниил', 2, '1', 4, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 179, 85, 100, 11, 11, 10, 15, 10, 15, 25),
    (50, 'Леванин Артем', 2, '0', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'shield', FALSE, 180, 86, 100, 11, 8, 12, 12, 15, 12, 18),
    (114, 'Добряков Иван', 2, '3', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 180, 85, 100, 9, 11, 8, 10, 10, 10, 15),
    (115, 'Ладыгин Александр', 2, '2', 2, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 178, 82, 100, 8, 8, 10, 12, 10, 12, 18),
    (116, 'Глинка Вадим', 2, '4', 2, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (117, 'Костин Сергей', 2, '2', 2, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 181, 87, 100, 8, 8, 10, 12, 10, 12, 18),

    (38, 'Голодяев Иван', 3, '1', 3, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 170, 76, 100, 9, 9, 10, 15, 8, 15, 22),
    (31, 'Комержеев Михаил', 3, '0', 4, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 170, 76, 100, 11, 8, 12, 12, 15, 12, 18),
    (205, 'Пальцев Алексей', 3, '1', 4, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 170, 76, 100, 11, 8, 12, 12, 15, 12, 18),
    (61, 'Алексеев Евгений', 3, '2', 3, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 181, 89, 100, 8, 8, 10, 12, 10, 12, 18),
    (69, 'Епифанов Роман', 3, '4', 3, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 179, 85, 100, 8, 8, 12, 12, 6, 12, 18),
    (118, 'Пашнин Артём', 3, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 176, 81, 100, 6, 6, 8, 10, 8, 10, 15),
    (119, 'Фадеев Сергей', 3, '3', 2, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 179, 84, 100, 9, 11, 8, 10, 10, 10, 15),
    (120, 'Гореванов Владимир', 3, '2', 2, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),

    (54, 'Штрейс Николай', 4, '4', 3, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 177, 82, 100, 8, 8, 12, 12, 6, 12, 18),
    (62, 'Черняев Павел', 4, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 169, 75, 100, 8, 8, 10, 12, 10, 12, 18),
    (63, 'Тягунов Антон', 4, '4', 3, './static/characters/default.png', FALSE, 'axe', NULL, FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (68, 'Власов Всеволод', 4, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 171, 77, 100, 8, 8, 10, 12, 10, 12, 18),
    (73, 'Савчук Сергей', 4, '1', 3, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 183, 90, 100, 9, 9, 10, 15, 8, 15, 22),
    (76, 'Киселев Михаил', 4, '2', 3, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 180, 86, 100, 8, 8, 10, 12, 10, 12, 18),
    (82, 'Шошин Максим', 4, '3', 3, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 181, 89, 100, 9, 11, 8, 10, 10, 10, 15),
    (121, 'Крепков Сергей', 4, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 170, 76, 100, 6, 6, 8, 10, 8, 10, 15),
    (122, 'Котти Константин', 4, '4', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 175, 80, 100, 8, 8, 12, 12, 6, 12, 18),

    (123, 'Пак Александр', 5, '2', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 178, 83, 100, 8, 8, 10, 12, 10, 12, 18),
    (124, 'Сажин Сергей', 5, '1', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 174, 79, 100, 9, 9, 10, 15, 8, 15, 22),
    (125, 'Сабанин Александр', 5, '2', 2, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 181, 87, 100, 8, 8, 10, 12, 10, 12, 18),
    (126, 'Евсеев Михаил', 5, '3', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', 'shield', FALSE, 176, 81, 100, 9, 11, 8, 10, 10, 10, 15),
    (127, 'Борисов Андрей', 5, '4', 2, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 179, 84, 100, 8, 8, 12, 12, 6, 12, 18),
    (128, 'Коваленков Владимир', 5, '2', 2, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),
    (129, 'Гордеев Александр', 5, '2', 2, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 170, 76, 100, 6, 6, 8, 10, 8, 10, 15),

    (3, 'Опарин Егор', 6, '1', 5, './static/characters/default.png', TRUE, 'falchion', 'buckler', TRUE, 198, 100, 100, 11, 11, 12, 15, 10, 15, 25),
    (5, 'Русанов Ярослав', 6, '0', 5, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, TRUE, 196, 105, 100, 11, 11, 12, 12, 15, 12, 18),
    (7, 'Соколов Савва', 6, '1', 4, './static/characters/default.png', TRUE, 'falchion', 'shield', TRUE, 172, 80, 100, 11, 11, 12, 15, 10, 15, 25),
    (9, 'Панченко Михаил', 6, '2', 4, './static/characters/default.png', TRUE, 'sword', 'buckler', TRUE, 174, 79, 100, 9, 11, 12, 12, 12, 12, 18),
    (18, 'Литвиненко Алексей', 6, '2', 4, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 173, 81, 100, 8, 8, 10, 12, 10, 12, 18),
    (20, 'Дусявичус', 6, '4', 4, './static/characters/default.png', FALSE, 'axe', NULL, FALSE, 172, 76, 100, 8, 8, 12, 12, 6, 12, 18),
    (43, 'Томилин Антон', 6, '2', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 180, 86, 100, 6, 6, 8, 10, 8, 10, 15),
    (46, 'Беляев Павел', 6, '2', 3, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 183, 90, 100, 6, 6, 8, 10, 8, 10, 15),
    (45, 'Туктарев Максим', 6, '3', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 176, 81, 100, 9, 11, 8, 10, 10, 10, 15),
    (130, 'Franco Strydom', 6, '4', 2, './static/characters/default.png', FALSE, 'two_handed_sword', NULL, FALSE, 175, 80, 100, 8, 8, 12, 12, 6, 12, 18),

    (131, 'Тулинов Алексей', 7, '2', 2, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 178, 83, 100, 8, 8, 10, 12, 10, 12, 18),
    (132, 'Иванов Юрий', 7, '4', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (133, 'Мусихин Сергей', 7, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 181, 87, 100, 8, 8, 10, 12, 10, 12, 18),
    (201, 'Глянцев Валерий', 7, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 181, 87, 100, 6, 6, 8, 10, 8, 10, 15),
    (202, 'Гольтянин Владимир', 7, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 181, 87, 100, 6, 6, 8, 10, 8, 10, 15),
    (203, 'Тимофеев Андрей', 7, '3', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 181, 87, 100, 9, 11, 8, 10, 10, 10, 15),

    (92, 'Копытенков Анатолий', 8, '2', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'shield', FALSE, 180, 86, 100, 8, 8, 10, 12, 10, 12, 18),
    (98, 'Илюхин Александр', 8, '2', 3, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 166, 73, 100, 6, 6, 8, 10, 8, 10, 15),
    (100, 'Огородний Родион', 8, '4', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 182, 88, 100, 8, 8, 12, 12, 6, 12, 18),
    (134, 'Бугрий Михаил', 8, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 176, 81, 100, 6, 6, 8, 10, 8, 10, 15),
    (135, 'Ершов Дмитрий', 8, '4', 2, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 179, 84, 100, 8, 8, 12, 12, 6, 12, 18),
    (136, 'Ивашко Александр', 8, '2', 2, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),

    (137, 'Демченков Антон', 9, '2', 2, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 170, 76, 100, 6, 6, 8, 10, 8, 10, 15),
    (138, 'Владимиров Александр', 9, '4', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 175, 80, 100, 8, 8, 12, 12, 6, 12, 18),
    (139, 'Карайман Александр', 9, '2', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 178, 83, 100, 8, 8, 10, 12, 10, 12, 18),
    (140, 'Панин Евгений', 9, '3', 2, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 174, 79, 100, 9, 11, 8, 10, 10, 10, 15),
    (141, 'Пугачев Георгий', 9, '4', 2, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 181, 87, 100, 8, 8, 12, 12, 6, 12, 18),
    (142, 'Высоцкий Роман', 9, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 176, 81, 100, 8, 8, 10, 12, 10, 12, 18),
    (143, 'Папян Иван', 9, '2', 2, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 179, 84, 100, 6, 6, 8, 10, 8, 10, 15),

    (144, 'Киселев Дмитрий', 10, '2', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),
    (145, 'Шмидт Андрей', 10, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 170, 76, 100, 6, 6, 8, 10, 8, 10, 15),
    (146, 'Байрамгулов Руслан', 10, '2', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 175, 80, 100, 8, 8, 10, 12, 10, 12, 18),
    (147, 'Чудинов Иван', 10, '2', 2, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 178, 83, 100, 6, 6, 8, 10, 8, 10, 15),
    (148, 'Политковский Лев', 10, '4', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (149, 'Укрюков Тимофей', 10, '3', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 181, 87, 100, 9, 11, 8, 10, 10, 10, 15),

    (1, 'Баксанов Бенедикт', 11, '0', 5, './static/characters/benya.png', TRUE, 'falchion', 'buckler', TRUE, 198, 110, 100, 11, 11, 12, 12, 15, 12, 18),
    (2, 'Клыков Александр', 11, '1', 5, './static/characters/sasha_klykov.png', TRUE, 'two_handed_halberd', NULL, TRUE, 180, 90, 100, 11, 11, 12, 15, 10, 15, 25),
    (4, 'Баранас Роман', 11, '0', 5, './static/characters/baranas_roma.png', TRUE, 'axe', 'buckler', TRUE, 192, 90, 100, 11, 11, 12, 12, 15, 12, 18),
    (8, 'Астошенок Александр', 11, '4', 4, './static/characters/astoshenok_sasha.png', TRUE, 'falchion', 'buckler', TRUE, 187, 100, 100, 11, 11, 15, 13, 10, 12, 22),
    (10, 'Голованов Николай', 11, '2', 4, './static/characters/golovanov_kolya.png', TRUE, 'falchion', 'buckler', TRUE, 200, 100, 100, 9, 11, 12, 12, 12, 12, 18),
    (13, 'Кунченко Дмитрий', 11, '2', 4, './static/characters/kunchenko_dima.png', TRUE, 'two_handed_halberd', NULL, TRUE, 198, 110, 100, 8, 11, 10, 12, 10, 12, 18),
    (16, 'Ткачук Никита', 11, '4', 4, './static/characters/tkachuk_nikita.png', TRUE, 'falchion', 'buckler', TRUE, 180, 80, 100, 8, 11, 12, 12, 6, 12, 18),
    (19, 'Надеждин Александр', 11, '2', 4, './static/characters/nadejdin_sasha.png', TRUE, 'sword', 'shield', TRUE, 175, 80, 100, 8, 8, 10, 12, 10, 12, 18),
    (26, 'Кравченко Игорь', 11, '2', 4, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 180, 85, 100, 8, 8, 10, 12, 10, 12, 18),
    ON CONFLICT (id) DO NOTHING;

INSERT INTO characters (id, name, team_id, role_id, count_of_ability, image_url, is_active, weapon, shield, is_titan_armour, height, weight, hp, stamina, initiative, wrestling, attack, defense, attack_min, attack_max)
VALUES
    -- TeamID 12: Злой дух Ямбуя
    (35, 'Вахрамеев Олег', 12, '2', 3, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 166, 73, 100, 6, 6, 8, 10, 8, 10, 15),
    (36, 'Сазанов Никита', 12, '2', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 175, 80, 100, 6, 6, 8, 10, 8, 10, 15),
    (58, 'Сметрин Кирилл', 12, '2', 3, './static/characters/default.png', FALSE, 'two_handed_halberd', 'buckler', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),
    (67, 'Маштаков Александр', 12, '4', 3, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 183, 90, 100, 8, 8, 12, 12, 6, 12, 18),
    (72, 'Неудачин Павел', 12, '2', 3, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 169, 75, 100, 6, 6, 8, 10, 8, 10, 15),
    (81, 'Веретенников Игорь', 12, '2', 3, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 178, 84, 100, 8, 8, 10, 12, 10, 12, 18),
    (90, 'Зыкин Александр', 12, '2', 3, './static/characters/default.png', FALSE, 'sword', NULL, FALSE, 179, 85, 100, 6, 6, 8, 10, 8, 10, 15),
    (93, 'Семёнов Евгений', 12, '2', 3, './static/characters/default.png', FALSE, 'axe', NULL, FALSE, 169, 75, 100, 6, 6, 8, 10, 8, 10, 15),
    (63, 'Тягунов Антон', 12, '4', 3, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (151, 'Берестнев Кирилл', 12, '2', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 179, 84, 100, 8, 8, 10, 12, 10, 12, 18),
    (6, 'Топоев Никита', 12, '1', 5, './static/characters/default.png', TRUE, 'falchion', 'shield', TRUE, 170, 78, 100, 11, 11, 12, 15, 10, 15, 25),

    -- TeamID 13: Мальтийский крест
    (87, 'Емалетдинов Евгений', 13, '2', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 176, 81, 100, 8, 8, 10, 12, 10, 12, 18),
    (99, 'Шарыкин Владимир', 13, '3', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 175, 80, 100, 9, 11, 8, 10, 10, 10, 15),
    (105, 'Нуриев Руслан', 13, '2', 3, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 174, 79, 100, 6, 6, 8, 10, 8, 10, 15),
    (152, 'Соболев Алексей', 13, '2', 2, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),
    (153, 'Шаихов Григорий', 13, '4', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 170, 76, 100, 8, 8, 12, 12, 6, 12, 18),
    (154, 'Булгаков Александр', 13, '2', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 175, 80, 100, 8, 8, 10, 12, 10, 12, 18),

    -- TeamID 14: Урфин Джус
    (37, 'Овчинников Михаил', 14, '2', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18),
    (42, 'Маркелов Алексей', 14, '2', 3, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 174, 79, 100, 6, 6, 8, 10, 8, 10, 15),
    (48, 'Голованов Илья', 14, '2', 3, './static/characters/default.png', TRUE, 'sword', NULL, FALSE, 179, 85, 100, 8, 8, 10, 12, 10, 12, 18),
    (96, 'Лебедев Василий', 14, '4', 3, './static/characters/default.png', TRUE, 'two_handed_sword', NULL, FALSE, 177, 82, 100, 8, 8, 12, 12, 6, 12, 18),
    (155, 'Сорокин Максим', 14, '3', 2, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 178, 83, 100, 9, 11, 8, 10, 10, 10, 15),
    (156, 'Стуров Владимир', 14, '4', 2, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 174, 79, 100, 8, 8, 12, 12, 6, 12, 18),
    (157, 'Миркович Георгий', 14, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 181, 87, 100, 8, 8, 10, 12, 10, 12, 18),
    (158, 'Митрофанов Сергей', 14, '2', 2, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 176, 81, 100, 6, 6, 8, 10, 8, 10, 15),

    -- TeamID 15: Высшая Школа ИСБ Санкт-Петербург
    (66, 'Коробков Михаил', 15, '2', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 176, 81, 100, 8, 8, 10, 12, 10, 12, 18),
    (159, 'Быстров Роман', 15, '2', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 179, 84, 100, 8, 8, 10, 12, 10, 12, 18),
    (160, 'Аржанцев Александр', 15, '1', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 182, 88, 100, 9, 9, 10, 15, 8, 15, 22),
    (161, 'Ильичев Андрей', 15, '2', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 170, 76, 100, 8, 8, 10, 12, 10, 12, 18),
    (162, 'Паутов Олег', 15, '2', 2, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 175, 80, 100, 6, 6, 8, 10, 8, 10, 15),
    (163, 'Тюляндин Иван', 15, '4', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 178, 83, 100, 8, 8, 12, 12, 6, 12, 18),
    (164, 'Плешанов Владислав', 15, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 174, 79, 100, 8, 8, 10, 12, 10, 12, 18),
    (165, 'Степанов Михаил', 15, '3', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 181, 87, 100, 9, 11, 8, 10, 10, 10, 15),

    -- TeamID 16: Медвежья пядь
    (11, 'Каменев', 16, '0', 4, './static/characters/default.png', TRUE, 'two_handed_halberd', 'shield', TRUE, 180, 87, 100, 11, 11, 12, 12, 15, 12, 18),
    (12, 'Грызлов Виталий', 16, '2', 4, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 169, 76, 100, 8, 8, 10, 12, 10, 12, 18),
    (17, 'Курицын Сергей', 16, '2', 4, './static/characters/default.png', TRUE, 'two_handed_halberd', 'shield', FALSE, 181, 89, 100, 8, 8, 10, 12, 10, 12, 18),
    (27, 'Намазов Рафаэль', 16, '4', 4, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 176, 81, 100, 8, 8, 12, 12, 6, 12, 18),
    (29, 'Никитин Александр', 16, '2', 4, './static/characters/default.png', TRUE, 'axe', 'shield', FALSE, 174, 79, 100, 6, 6, 8, 10, 8, 10, 15),
    (33, 'Марычев Михаил', 16, '2', 4, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 177, 82, 100, 6, 6, 8, 10, 8, 10, 15),
    (34, 'Балясников Антон', 16, '2', 4, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 180, 86, 100, 8, 8, 10, 12, 10, 12, 18),
    (166, 'Цырулик Владимир', 16, '4', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 176, 81, 100, 8, 8, 12, 12, 6, 12, 18),

    -- TeamID 17: Байард
    (30, 'Шавлакадзе Эдуард', 17, '2', 4, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 169, 75, 100, 8, 8, 10, 12, 10, 12, 18),
    (53, 'Нойманн Кирилл', 17, '3', 3, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 171, 77, 100, 9, 11, 8, 10, 10, 10, 15),
    (64, 'Шостаковский Антон', 17, '2', 3, './static/characters/default.png', TRUE, 'two_handed_halberd', 'buckler', FALSE, 180, 86, 100, 8, 8, 10, 12, 10, 12, 18),
    (78, 'Вершинин Святослав', 17, '2', 3, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 175, 80, 100, 6, 6, 8, 10, 8, 10, 15),
    (49, 'Васильев Иван', 17, '2', 3, './static/characters/default.png', FALSE, 'falchion', 'buckler', FALSE, 174, 79, 100, 8, 8, 10, 12, 10, 12, 18),
    (167, 'Галкин Алексей', 17, '2', 2, './static/characters/default.png', TRUE, 'falchion', 'buckler', FALSE, 179, 84, 100, 8, 8, 10, 12, 10, 12, 18),
    (168, 'Цзинь Фэнхао', 17, '4', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 182, 88, 100, 8, 8, 12, 12, 6, 12, 18),
    (169, 'Зуев Роман', 17, '2', 2, './static/characters/default.png', TRUE, 'axe', NULL, FALSE, 170, 76, 100, 6, 6, 8, 10, 8, 10, 15),
    (170, 'Францкевич Эдуард', 17, '2', 2, './static/characters/default.png', TRUE, 'sword', 'buckler', FALSE, 175, 80, 100, 8, 8, 10, 12, 10, 12, 18),

    -- TeamID 18: RaubRitter
    (32, 'Найдеров Алексей', 18, '4', 4, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 171, 77, 100, 8, 8, 12, 12, 6, 12, 18),
    (41, 'Козлов Михаил', 18, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 169, 75, 100, 8, 8, 10, 12, 10, 12, 18),
    (56, 'Козырев Александр', 18, '2', 3, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 166, 73, 100, 6, 6, 8, 10, 8, 10, 15),
    (83, 'Поташный Климентий', 18, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 169, 75, 100, 8, 8, 10, 12, 10, 12, 18),
    (88, 'Коваленко Павел', 18, '1', 3, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 183, 90, 100, 9, 9, 10, 15, 8, 15, 22),
    (171, 'Воронов Ричард', 18, '4', 2, './static/characters/default.png', TRUE, 'falchion', 'shield', FALSE, 178, 83, 100, 8, 8, 12, 12, 6, 12, 18),

    -- TeamID 19: Межевой рыцарь
    (89, 'Басов Захар', 19, '2', 3, './static/characters/default.png', TRUE, 'two_handed_sword', 'shield', FALSE, 171, 77, 100, 8, 8, 10, 12, 10, 12, 18),
    (172, 'Василенко Михаил', 19, '2', 2, './static/characters/default.png', TRUE, 'two_handed_halberd', NULL, FALSE, 174, 79, 100, 6, 6, 8, 10, 8, 10, 15),
    (173, 'Безуглый Георгий', 19, '4', 2, './static/characters/default.png', TRUE, 'axe', 'buckler', FALSE, 181, 87, 100, 8, 8, 12, 12, 6, 12, 18),
    (174, 'Платонов Игорь', 19, '2', 2, './static/characters/default.png', TRUE, 'sword', 'shield', FALSE, 176, 81, 100, 8, 8, 10, 12, 10, 12, 18),
    (175, 'Тищенко Алексей', 19, '2', 2, './static/characters/default.png', TRUE, 'falchion', NULL, FALSE, 179, 84, 100, 6, 6, 8, 10, 8, 10, 15),
    (176, 'Никитин Богдан', 19, '2', 2, './static/characters/default.png', TRUE, 'two_handed_sword', 'buckler', FALSE, 182, 88, 100, 8, 8, 10, 12, 10, 12, 18)
    ON CONFLICT (id) DO NOTHING;
