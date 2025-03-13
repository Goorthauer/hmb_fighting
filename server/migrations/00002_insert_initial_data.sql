-- +goose Up
INSERT INTO weapons (name, display_name, range, is_two_handed, image_url, attack_bonus, grapple_bonus)
VALUES
    ('falchion', 'Фальшион', 1, FALSE, './static/weapons/default.png', 2, 0),
    ('axe', 'Топор', 1, FALSE, './static/weapons/default.png', 0, 8),
    ('two_handed_sword', 'Двуручный меч', 2, TRUE, './static/weapons/default.png', 2, 0),
    ('two_handed_halberd', 'Алебарда', 2, TRUE, './static/weapons/default.png', 2, 8),
    ('sword', 'Меч', 1, FALSE, './static/weapons/default.png', 0, 0);

INSERT INTO shields (name, display_name, defense_bonus, image_url, attack_bonus, grapple_bonus)
VALUES
    ('buckler', 'Баклер', 1, './static/shields/default.png', 1, 1),
    ('shield', 'Тарч', 2, './static/shields/default.png', 1, 0),
    ('tower', 'Ростовой щит', 3, './static/shields/default.png', 0, -1);

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
    ('underhook_throw', 'Бросок через подхват', 'wrestle', 'Бросок противника через нижний захват руки.', 1, './static/abilities/default.jpg');

INSERT INTO roles (id, name)
VALUES
    ('0', 'Танк'),
    ('1', 'Убийца'),
    ('2', 'Боец'),
    ('3', 'Поддержка'),
    ('4', 'Борец');

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
    (19, 'Межевой рыцарь', './static/teams/border_knight.png', 'Стражи границ и традиций.');

-- +goose Down
TRUNCATE TABLE teams;
TRUNCATE TABLE roles;
TRUNCATE TABLE abilities;
TRUNCATE TABLE shields;
TRUNCATE TABLE weapons;