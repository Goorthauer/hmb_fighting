-- +goose Up
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
-- Добавь остальные записи из init.sql сюда по аналогии

-- +goose Down
TRUNCATE TABLE characters;