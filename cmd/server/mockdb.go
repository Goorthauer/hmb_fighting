package main

// MockDatabase имитирует базу данных с предопределёнными данными
type MockDatabase struct{}

// GetWeapons возвращает конфигурацию оружия
func (m *MockDatabase) GetWeapons() (map[string]Weapon, error) {
	return map[string]Weapon{
		"falchion":           {Name: "falchion", DisplayName: "Фальшион", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"axe":                {Name: "axe", DisplayName: "Топор", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 8},
		"two_handed_sword":   {Name: "two_handed_sword", DisplayName: "Двуручный меч", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 0},
		"two_handed_halberd": {Name: "two_handed_halberd", DisplayName: "Алебарда", Range: 2, IsTwoHanded: true, ImageURL: "./static/weapons/default.png", AttackBonus: 2, GrappleBonus: 8},
		"sword":              {Name: "sword", DisplayName: "Меч", Range: 1, IsTwoHanded: false, ImageURL: "./static/weapons/default.png", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetShields() (map[string]Shield, error) {
	return map[string]Shield{
		"buckler": {Name: "buckler", DisplayName: "Баклер", DefenseBonus: 1, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 1},
		"shield":  {Name: "shield", DisplayName: "Тарч", DefenseBonus: 2, ImageURL: "./static/shields/default.png", AttackBonus: 1, GrappleBonus: 0},
		"tower":   {Name: "tower", DisplayName: "Ростовой щит", DefenseBonus: 3, ImageURL: "./static/shields/default.png", AttackBonus: 0, GrappleBonus: -1},
		"":        {Name: "none", DefenseBonus: 0, ImageURL: "", AttackBonus: 0, GrappleBonus: 0},
	}, nil
}

func (m *MockDatabase) GetAbilities() (map[string]Ability, error) {
	abilities := map[string]Ability{
		"yama_arashi": {
			Name:        "yama_arashi",
			DisplayName: "Подхват",
			Description: "Мощный бросок через бедро, использующий силу и импульс противника для стремительного падения.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"trip": {
			Name:        "trip",
			DisplayName: "Зацеп",
			Description: "Точный удар по ноге, нарушающий баланс противника и заставляющий его упасть.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"hip_toss": {
			Name:        "hip_toss",
			DisplayName: "Высед",
			Description: "Бросок через бедро с использованием вращения и силы противника.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"front_sweep": {
			Name:        "front_sweep",
			DisplayName: "Передняя подножка",
			Description: "Быстрый удар по передней ноге, опрокидывающий противника вперёд.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"back_sweep": {
			Name:        "back_sweep",
			DisplayName: "Задняя подножка",
			Description: "Подсечка задней ноги, использующая вес противника для опрокидывания.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"outer_hook": {
			Name:        "outer_hook",
			DisplayName: "Внешний зацеп",
			Description: "Подсечка внешней стороны ноги, нарушающая равновесие противника.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"inner_hook": {
			Name:        "inner_hook",
			DisplayName: "Внутренний зацеп",
			Description: "Подсечка внутренней стороны ноги, выводящая противника из равновесия.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"shoulder_throw": {
			Name:        "shoulder_throw",
			DisplayName: "Бросок через плечо",
			Description: "Рывок противника через плечо с использованием его импульса.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"spinning_throw": {
			Name:        "spinning_throw",
			DisplayName: "Вращательный бросок",
			Description: "Бросок противника через спину с использованием вращения.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"double_leg": {
			Name:        "double_leg",
			DisplayName: "Двойной захват ног",
			Description: "Захват обеих ног противника с последующим опрокидыванием.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"single_leg": {
			Name:        "single_leg",
			DisplayName: "Одинарный захват ноги",
			Description: "Захват одной ноги противника с последующим броском.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"overhook_throw": {
			Name:        "overhook_throw",
			DisplayName: "Бросок через захват",
			Description: "Бросок противника через верхний захват руки.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
		"underhook_throw": {
			Name:        "underhook_throw",
			DisplayName: "Бросок через подхват",
			Description: "Бросок противника через нижний захват руки.",
			Type:        "wrestle",
			Range:       1,
			ImageURL:    "./static/abilities/default.jpg",
		},
	}
	return abilities, nil
}

// GetCharacters возвращает список персонажей

func (m *MockDatabase) GetCharacters() ([]Character, error) {
	return []Character{
		// TeamID 1: Партизан Два
		{ID: 15, Name: "Тюляков Алексей", TeamID: 1, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 177, Weight: 84, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},         // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 14, Name: "Чуклов Григорий", TeamID: 1, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 12, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "shield", Height: 171, Weight: 77, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},           // Боец (было Stamina: 10, Initiative: 10)
		{ID: 22, Name: "Корчагин Максим", TeamID: 1, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "sword", Shield: "buckler", Height: 175, Weight: 80, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},              // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 23, Name: "Шеварков Иван", TeamID: 1, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "shield", Height: 170, Weight: 78, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                              // Боец (было Stamina: 8, Initiative: 8)
		{ID: 24, Name: "Остраумов Антон", TeamID: 1, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 8, Initiative: 8, Wrestling: 14, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 182, Weight: 88, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},        // Борец (было Stamina: 10, Initiative: 10)
		{ID: 28, Name: "Наумов Александр", TeamID: 1, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "two_handed_sword", Shield: "buckler", Height: 181, Weight: 87, CountOfAbility: 4, ImageURL: "./static/characters/default.png"}, // Поддержка (было Stamina: 12, Initiative: 15)
		{ID: 40, Name: "Сотов Николай", TeamID: 1, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "buckler", Height: 181, Weight: 89, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},             // Боец (было Stamina: 8, Initiative: 8)
		{ID: 21, Name: "Ивлев Владимир", TeamID: 1, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 167, Weight: 72, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},           // Борец (было Stamina: 10, Initiative: 10)
		{ID: 75, Name: "Коврижин Дмитрий", TeamID: 1, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "two_handed_sword", Shield: "", Height: 177, Weight: 82, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},          // Боец (было Stamina: 8, Initiative: 8)
		{ID: 112, Name: "Янышевский Максим", TeamID: 1, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "buckler", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                         // Боец (было Stamina: 8, Initiative: 8)
		{ID: 113, Name: "Флерко Игорь", TeamID: 1, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "", Height: 172, Weight: 78, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                                         // Борец (было Stamina: 10, Initiative: 10)

		// TeamID 2: Юг
		{ID: 51, Name: "Рябцев Сергей", TeamID: 2, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "axe", Shield: "", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                        // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 25, Name: "Свиридов Даниил", TeamID: 2, IsActive: true, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 15, AttackMax: 25, Defense: 10, Initiative: 11, Wrestling: 10, Attack: 15, Weapon: "axe", Shield: "buckler", Height: 179, Weight: 85, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},            // Убийца (было Stamina: 14, Initiative: 14)
		{ID: 50, Name: "Леванин Артем", TeamID: 2, IsActive: true, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_halberd", Shield: "shield", Height: 180, Weight: 86, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Танк (было Stamina: 14, Initiative: 10)
		{ID: 114, Name: "Добряков Иван", TeamID: 2, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "shield", Height: 180, Weight: 85, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 115, Name: "Ладыгин Александр", TeamID: 2, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 178, Weight: 82, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 116, Name: "Глинка Вадим", TeamID: 2, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},            // Борец (было Stamina: 10, Initiative: 10)
		{ID: 117, Name: "Костин Сергей", TeamID: 2, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "shield", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 3: НСК
		{ID: 38, Name: "Голодяев Иван", TeamID: 3, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "axe", Shield: "shield", Height: 170, Weight: 76, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},            // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 31, Name: "Комержеев Михаил", TeamID: 3, IsActive: true, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 170, Weight: 76, CountOfAbility: 4, ImageURL: "./static/characters/default.png"}, // Танк (было Stamina: 14, Initiative: 10)
		{ID: 205, Name: "Пальцев Алексей", TeamID: 3, IsActive: true, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 170, Weight: 76, CountOfAbility: 4, ImageURL: "./static/characters/default.png"}, // Убийца (было Stamina: 14, Initiative: 10)
		{ID: 61, Name: "Алексеев Евгений", TeamID: 3, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 181, Weight: 89, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},  // Боец (было Stamina: 10, Initiative: 10)
		{ID: 69, Name: "Епифанов Роман", TeamID: 3, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "sword", Shield: "", Height: 179, Weight: 85, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},               // Борец (было Stamina: 10, Initiative: 10)
		{ID: 118, Name: "Пашнин Артём", TeamID: 3, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},    // Боец (было Stamina: 8, Initiative: 8)
		{ID: 119, Name: "Фадеев Сергей", TeamID: 3, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "buckler", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},       // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 120, Name: "Гореванов Владимир", TeamID: 3, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "shield", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 4: Школа ИСБ Байард
		{ID: 54, Name: "Штрейс Николай", TeamID: 4, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 177, Weight: 82, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},           // Борец (было Stamina: 10, Initiative: 10)
		{ID: 62, Name: "Черняев Павел", TeamID: 4, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 63, Name: "Тягунов Антон", TeamID: 4, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "", Height: 174, Weight: 79, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                                         // Борец (было Stamina: 10, Initiative: 10)
		{ID: 68, Name: "Власов Всеволод", TeamID: 4, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 171, Weight: 77, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},   // Боец (было Stamina: 10, Initiative: 10)
		{ID: 73, Name: "Савчук Сергей", TeamID: 4, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "sword", Shield: "buckler", Height: 183, Weight: 90, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 76, Name: "Киселев Михаил", TeamID: 4, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 180, Weight: 86, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                // Боец (было Stamina: 10, Initiative: 10)
		{ID: 82, Name: "Шошин Максим", TeamID: 4, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "buckler", Height: 181, Weight: 89, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},             // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 121, Name: "Крепков Сергей", TeamID: 4, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                   // Боец (было Stamina: 8, Initiative: 8)
		{ID: 122, Name: "Котти Константин", TeamID: 4, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "buckler", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Борец (было Stamina: 10, Initiative: 10)

		// TeamID 5: НРБ
		{ID: 123, Name: "Пак Александр", TeamID: 5, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "shield", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Боец (было Stamina: 10, Initiative: 10)
		{ID: 124, Name: "Сажин Сергей", TeamID: 5, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "axe", Shield: "", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                        // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 125, Name: "Сабанин Александр", TeamID: 5, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},      // Боец (было Stamina: 10, Initiative: 10)
		{ID: 126, Name: "Евсеев Михаил", TeamID: 5, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "shield", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 127, Name: "Борисов Андрей", TeamID: 5, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "sword", Shield: "", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                    // Борец (было Stamina: 10, Initiative: 10)
		{ID: 128, Name: "Коваленков Владимир", TeamID: 5, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},         // Боец (было Stamina: 10, Initiative: 10)
		{ID: 129, Name: "Гордеев Александр", TeamID: 5, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "shield", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},         // Боец (было Stamina: 8, Initiative: 8)

		// TeamID 6: Старые Друзья
		{ID: 3, Name: "Опарин Егор", TeamID: 6, IsActive: true, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 15, AttackMax: 25, Defense: 10, Initiative: 11, Wrestling: 12, Attack: 15, Weapon: "falchion", Shield: "buckler", Height: 198, Weight: 100, CountOfAbility: 5, ImageURL: "./static/characters/default.png", IsTitanArmour: true},        // Убийца (было Stamina: 14, Initiative: 15)
		{ID: 5, Name: "Русанов Ярослав", TeamID: 6, IsActive: true, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 196, Weight: 105, CountOfAbility: 5, ImageURL: "./static/characters/default.png", IsTitanArmour: true}, // Танк (было Stamina: 14, Initiative: 15)
		{ID: 7, Name: "Соколов Савва", TeamID: 6, IsActive: true, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 15, AttackMax: 25, Defense: 10, Initiative: 11, Wrestling: 12, Attack: 15, Weapon: "falchion", Shield: "shield", Height: 172, Weight: 80, CountOfAbility: 4, ImageURL: "./static/characters/default.png", IsTitanArmour: true},        // Убийца (было Stamina: 14, Initiative: 15)
		{ID: 9, Name: "Панченко Михаил", TeamID: 6, IsActive: true, RoleID: 2, HP: 100, Stamina: 9, AttackMin: 12, AttackMax: 18, Defense: 12, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "sword", Shield: "buckler", Height: 174, Weight: 79, CountOfAbility: 4, ImageURL: "./static/characters/default.png", IsTitanArmour: true},         // Боец (было Stamina: 12, Initiative: 15)
		{ID: 18, Name: "Литвиненко Алексей", TeamID: 6, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "", Height: 173, Weight: 81, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                                  // Боец (было Stamina: 10, Initiative: 10)
		{ID: 20, Name: "Дусявичус", TeamID: 6, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "", Height: 172, Weight: 76, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                                                              // Борец (было Stamina: 10, Initiative: 10)
		{ID: 43, Name: "Томилин Антон", TeamID: 6, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "buckler", Height: 180, Weight: 86, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                     // Боец (было Stamina: 8, Initiative: 8)
		{ID: 46, Name: "Беляев Павел", TeamID: 6, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "buckler", Height: 183, Weight: 90, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                                     // Боец (было Stamina: 8, Initiative: 8)
		{ID: 45, Name: "Туктарев Максим", TeamID: 6, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 176, Weight: 81, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                                  // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 130, Name: "Franco Strydom", TeamID: 6, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                                           // Борец (было Stamina: 10, Initiative: 10)

		// TeamID 7: Черная Земля
		{ID: 131, Name: "Тулинов Алексей", TeamID: 7, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 132, Name: "Иванов Юрий", TeamID: 7, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "sword", Shield: "shield", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},     // Борец (было Stamina: 10, Initiative: 10)
		{ID: 133, Name: "Мусихин Сергей", TeamID: 7, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},    // Боец (было Stamina: 10, Initiative: 10)
		{ID: 201, Name: "Глянцев Валерий", TeamID: 7, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 8, Initiative: 8)
		{ID: 202, Name: "Гольтянин Владимир", TeamID: 7, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},  // Боец (было Stamina: 8, Initiative: 8)
		{ID: 203, Name: "Тимофеев Андрей", TeamID: 7, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},   // Поддержка (было Stamina: 12, Initiative: 14)

		// TeamID 8: Vivus Ferro
		{ID: 92, Name: "Копытенков Анатолий", TeamID: 8, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "shield", Height: 180, Weight: 86, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 98, Name: "Илюхин Александр", TeamID: 8, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "shield", Height: 166, Weight: 73, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                   // Боец (было Stamina: 8, Initiative: 8)
		{ID: 100, Name: "Огородний Родион", TeamID: 8, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_halberd", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},   // Борец (было Stamina: 10, Initiative: 10)
		{ID: 134, Name: "Бугрий Михаил", TeamID: 8, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "buckler", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},       // Боец (было Stamina: 8, Initiative: 8)
		{ID: 135, Name: "Ершов Дмитрий", TeamID: 8, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "shield", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                      // Борец (было Stamina: 10, Initiative: 10)
		{ID: 136, Name: "Ивашко Александр", TeamID: 8, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                      // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 9: Молодые львы
		{ID: 137, Name: "Демченков Антон", TeamID: 9, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "buckler", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Боец (было Stamina: 8, Initiative: 8)
		{ID: 138, Name: "Владимиров Александр", TeamID: 9, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Борец (было Stamina: 10, Initiative: 10)
		{ID: 139, Name: "Карайман Александр", TeamID: 9, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 140, Name: "Панин Евгений", TeamID: 9, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "buckler", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                 // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 141, Name: "Пугачев Георгий", TeamID: 9, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "shield", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Борец (было Stamina: 10, Initiative: 10)
		{ID: 142, Name: "Высоцкий Роман", TeamID: 9, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},          // Боец (было Stamina: 10, Initiative: 10)
		{ID: 143, Name: "Папян Иван", TeamID: 9, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "buckler", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                        // Боец (было Stamina: 8, Initiative: 8)

		// TeamID 10: Ганза
		{ID: 144, Name: "Киселев Дмитрий", IsActive: true, TeamID: 10, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "shield", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                // Боец (было Stamina: 10, Initiative: 10)
		{ID: 145, Name: "Шмидт Андрей", IsActive: true, TeamID: 10, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                        // Боец (было Stamina: 8, Initiative: 8)
		{ID: 146, Name: "Байрамгулов Руслан", IsActive: true, TeamID: 10, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "buckler", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 147, Name: "Чудинов Иван", IsActive: true, TeamID: 10, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "shield", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                       // Боец (было Stamina: 8, Initiative: 8)
		{ID: 148, Name: "Политковский Лев", IsActive: true, TeamID: 10, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                   // Борец (было Stamina: 10, Initiative: 10)
		{ID: 149, Name: "Укрюков Тимофей", IsActive: true, TeamID: 10, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "buckler", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},  // Поддержка (было Stamina: 12, Initiative: 14)

		// TeamID 11: Партизан
		{ID: 1, Name: "Баксанов Бенедикт", IsActive: true, TeamID: 11, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 198, Weight: 110, CountOfAbility: 5, ImageURL: "./static/characters/benya.png", IsTitanArmour: true},              // Танк (было Stamina: 15, Initiative: 15)
		{ID: 2, Name: "Клыков Александр", IsActive: true, TeamID: 11, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 15, AttackMax: 25, Defense: 10, Initiative: 11, Wrestling: 12, Attack: 15, Weapon: "two_handed_halberd", Shield: "", Height: 180, Weight: 90, CountOfAbility: 5, ImageURL: "./static/characters/sasha_klykov.png", IsTitanArmour: true},      // Убийца (было Stamina: 14, Initiative: 15)
		{ID: 4, Name: "Баранас Роман", IsActive: true, TeamID: 11, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 192, Weight: 90, CountOfAbility: 5, ImageURL: "./static/characters/baranas_roma.png", IsTitanArmour: true},                 // Танк (было Stamina: 15, Initiative: 15)
		{ID: 8, Name: "Астошенок Александр", IsActive: true, TeamID: 11, RoleID: 4, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 22, Defense: 10, Initiative: 11, Wrestling: 15, Attack: 13, Weapon: "falchion", Shield: "buckler", Height: 187, Weight: 100, CountOfAbility: 4, ImageURL: "./static/characters/astoshenok_sasha.png", IsTitanArmour: true}, // Борец (было Stamina: 14, Initiative: 15)
		{ID: 10, Name: "Голованов Николай", IsActive: true, TeamID: 11, RoleID: 2, HP: 100, Stamina: 9, AttackMin: 12, AttackMax: 18, Defense: 12, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 200, Weight: 100, CountOfAbility: 4, ImageURL: "./static/characters/golovanov_kolya.png", IsTitanArmour: true},    // Боец (было Stamina: 12, Initiative: 15)
		{ID: 13, Name: "Кунченко Дмитрий", IsActive: true, TeamID: 11, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 11, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 198, Weight: 110, CountOfAbility: 4, ImageURL: "./static/characters/kunchenko_dima.png", IsTitanArmour: true},   // Боец (было Stamina: 10, Initiative: 15)
		{ID: 16, Name: "Ткачук Никита", IsActive: true, TeamID: 11, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 180, Weight: 80, CountOfAbility: 4, ImageURL: "./static/characters/tkachuk_nikita.png", IsTitanArmour: true},           // Борец (было Stamina: 10, Initiative: 15)
		{ID: 19, Name: "Надеждин Александр", IsActive: true, TeamID: 11, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 11, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "shield", Height: 175, Weight: 80, CountOfAbility: 4, ImageURL: "./static/characters/nadejdin_sasha.png", IsTitanArmour: true},         // Боец (было Stamina: 10, Initiative: 15)
		{ID: 26, Name: "Калугин Григорий", TeamID: 11, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "shield", Height: 185, Weight: 105, CountOfAbility: 4, ImageURL: "./static/characters/default.png", IsTitanArmour: true},                               // Поддержка (было Stamina: 12, Initiative: 15)
		{ID: 150, Name: "Кравченко Игорь", TeamID: 11, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 11, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 176, Weight: 95, CountOfAbility: 2, ImageURL: "./static/characters/default.png", IsTitanArmour: true},                              // Боец (было Stamina: 10, Initiative: 15)
		{ID: 6, Name: "Топоев Никита", IsActive: true, TeamID: 12, RoleID: 1, HP: 100, Stamina: 11, AttackMin: 15, AttackMax: 25, Defense: 10, Initiative: 11, Wrestling: 12, Attack: 15, Weapon: "falchion", Shield: "shield", Height: 170, Weight: 78, CountOfAbility: 5, ImageURL: "./static/characters/default.png", IsTitanArmour: true},                  // Убийца (было Stamina: 14, Initiative: 15)

		// TeamID 13: ЗДЯ
		{ID: 35, Name: "Вахрамеев Олег", IsActive: true, TeamID: 12, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "shield", Height: 166, Weight: 73, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 8, Initiative: 8)
		{ID: 36, Name: "Сазанов Никита", IsActive: true, TeamID: 12, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 175, Weight: 80, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},        // Боец (было Stamina: 8, Initiative: 8)
		{ID: 58, Name: "Сметрин Кирилл", TeamID: 12, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 67, Name: "Маштаков Александр", IsActive: true, TeamID: 12, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 183, Weight: 90, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Борец (было Stamina: 10, Initiative: 10)
		{ID: 72, Name: "Неудачин Павел", IsActive: true, TeamID: 12, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},             // Боец (было Stamina: 8, Initiative: 8)
		{ID: 81, Name: "Веретенников Игорь", IsActive: true, TeamID: 12, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "", Height: 178, Weight: 84, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 90, Name: "Зыкин Александр", TeamID: 12, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "", Height: 179, Weight: 85, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                          // Боец (было Stamina: 8, Initiative: 8)
		{ID: 93, Name: "Семёнов Евгений", TeamID: 12, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                            // Боец (было Stamina: 8, Initiative: 8)
		{ID: 63, Name: "Тягунов Антон", IsActive: true, TeamID: 12, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "", Height: 174, Weight: 79, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},             // Борец (было Stamina: 10, Initiative: 10)
		{ID: 151, Name: "Берестнев Кирилл", IsActive: true, TeamID: 12, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},        // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 13: Мальтийский крест
		{ID: 87, Name: "Емалетдинов Евгений", IsActive: true, TeamID: 13, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "", Height: 176, Weight: 81, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},           // Боец (было Stamina: 10, Initiative: 10)
		{ID: 99, Name: "Шарыкин Владимир", IsActive: true, TeamID: 13, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 175, Weight: 80, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},              // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 105, Name: "Нуриев Руслан", IsActive: true, TeamID: 13, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 174, Weight: 79, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                       // Боец (было Stamina: 8, Initiative: 8)
		{ID: 152, Name: "Соболев Алексей", IsActive: true, TeamID: 13, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},       // Боец (было Stamina: 10, Initiative: 10)
		{ID: 153, Name: "Шаихов Григорий", IsActive: true, TeamID: 13, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Борец (было Stamina: 10, Initiative: 10)
		{ID: 154, Name: "Булгаков Александр", IsActive: true, TeamID: 13, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 14: Урфин Джус
		{ID: 37, Name: "Овчинников Михаил", IsActive: true, TeamID: 14, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 42, Name: "Маркелов Алексей", IsActive: true, TeamID: 14, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 174, Weight: 79, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                          // Боец (было Stamina: 8, Initiative: 8)
		{ID: 48, Name: "Голованов Илья", IsActive: true, TeamID: 14, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "", Height: 179, Weight: 85, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                        // Боец (было Stamina: 10, Initiative: 10)
		{ID: 96, Name: "Лебедев Василий", IsActive: true, TeamID: 14, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "", Height: 177, Weight: 82, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},             // Борец (было Stamina: 10, Initiative: 10)
		{ID: 155, Name: "Сорокин Максим", IsActive: true, TeamID: 14, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "buckler", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 156, Name: "Стуров Владимир", IsActive: true, TeamID: 14, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "shield", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Борец (было Stamina: 10, Initiative: 10)
		{ID: 157, Name: "Миркович Георгий", IsActive: true, TeamID: 14, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},        // Боец (было Stamina: 10, Initiative: 10)
		{ID: 158, Name: "Митрофанов Сергей", IsActive: true, TeamID: 14, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "buckler", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                 // Боец (было Stamina: 8, Initiative: 8)

		// TeamID 15: Высшая Школа ИСБ Санкт-Петербург
		{ID: 66, Name: "Коробков Михаил", IsActive: true, TeamID: 15, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "", Height: 176, Weight: 81, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                      // Боец (было Stamina: 10, Initiative: 10)
		{ID: 159, Name: "Быстров Роман", IsActive: true, TeamID: 15, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "shield", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                    // Боец (было Stamina: 10, Initiative: 10)
		{ID: 160, Name: "Аржанцев Александр", IsActive: true, TeamID: 15, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "falchion", Shield: "", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                   // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 161, Name: "Ильичев Андрей", IsActive: true, TeamID: 15, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "buckler", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},       // Боец (было Stamina: 10, Initiative: 10)
		{ID: 162, Name: "Паутов Олег", IsActive: true, TeamID: 15, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "shield", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                          // Боец (было Stamina: 8, Initiative: 8)
		{ID: 163, Name: "Тюляндин Иван", IsActive: true, TeamID: 15, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                        // Борец (было Stamina: 10, Initiative: 10)
		{ID: 164, Name: "Плешанов Владислав", IsActive: true, TeamID: 15, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "buckler", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 165, Name: "Степанов Михаил", IsActive: true, TeamID: 15, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "shield", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                  // Поддержка (было Stamina: 12, Initiative: 14)

		// TeamID 16: Медвежья пядь
		{ID: 11, Name: "Каменев", TeamID: 16, IsActive: true, RoleID: 0, HP: 100, Stamina: 11, AttackMin: 12, AttackMax: 18, Defense: 15, Initiative: 11, Wrestling: 12, Attack: 12, Weapon: "two_handed_halberd", Shield: "shield", Height: 180, Weight: 87, CountOfAbility: 4, ImageURL: "./static/characters/default.png", IsTitanArmour: true}, // Танк (было Stamina: 14, Initiative: 15)
		{ID: 12, Name: "Грызлов Виталий", TeamID: 16, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "axe", Shield: "", Height: 169, Weight: 76, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                                     // Боец (было Stamina: 10, Initiative: 10)
		{ID: 17, Name: "Курицын Сергей", TeamID: 16, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "shield", Height: 181, Weight: 89, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                 // Боец (было Stamina: 10, Initiative: 10)
		{ID: 27, Name: "Намазов Рафаэль", TeamID: 16, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "", Height: 176, Weight: 81, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                                 // Борец (было Stamina: 10, Initiative: 10)
		{ID: 29, Name: "Никитин Александр", TeamID: 16, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "shield", Height: 174, Weight: 79, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                               // Боец (было Stamina: 8, Initiative: 8)
		{ID: 33, Name: "Марычев Михаил", TeamID: 16, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 177, Weight: 82, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                                        // Боец (было Stamina: 8, Initiative: 8)
		{ID: 34, Name: "Балясников Антон", TeamID: 16, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "buckler", Height: 180, Weight: 86, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},                // Боец (было Stamina: 10, Initiative: 10)
		{ID: 166, Name: "Цырулик Владимир", TeamID: 16, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                                    // Борец (было Stamina: 10, Initiative: 10)

		// TeamID 17: Байард
		{ID: 30, Name: "Шавлакадзе Эдуард", TeamID: 17, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "", Height: 169, Weight: 75, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},         // Боец (было Stamina: 10, Initiative: 10)
		{ID: 53, Name: "Нойманн Кирилл", TeamID: 17, IsActive: true, RoleID: 3, HP: 100, Stamina: 9, AttackMin: 10, AttackMax: 15, Defense: 10, Initiative: 11, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "shield", Height: 171, Weight: 77, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                // Поддержка (было Stamina: 12, Initiative: 14)
		{ID: 64, Name: "Шостаковский Антон", TeamID: 17, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_halberd", Shield: "buckler", Height: 180, Weight: 86, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 78, Name: "Вершинин Святослав", TeamID: 17, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 175, Weight: 80, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                    // Боец (было Stamina: 8, Initiative: 8)
		{ID: 49, Name: "Васильев Иван", TeamID: 17, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 174, Weight: 79, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                                // Боец (было Stamina: 10, Initiative: 10)
		{ID: 167, Name: "Галкин Алексей", TeamID: 17, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "falchion", Shield: "buckler", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Боец (было Stamina: 10, Initiative: 10)
		{ID: 168, Name: "Цзинь Фэнхао", TeamID: 17, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},          // Борец (было Stamina: 10, Initiative: 10)
		{ID: 169, Name: "Зуев Роман", TeamID: 17, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "axe", Shield: "", Height: 170, Weight: 76, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                                // Боец (было Stamina: 8, Initiative: 8)
		{ID: 170, Name: "Францкевич Эдуард", IsActive: true, TeamID: 17, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "buckler", Height: 175, Weight: 80, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},              // Боец (было Stamina: 10, Initiative: 10)

		// TeamID 18: RaubRitter
		{ID: 32, Name: "Найдеров Алексей", TeamID: 18, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "shield", Height: 171, Weight: 77, CountOfAbility: 4, ImageURL: "./static/characters/default.png"},            // Борец (было Stamina: 10, Initiative: 10)
		{ID: 41, Name: "Козлов Михаил", TeamID: 18, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},      // Боец (было Stamina: 10, Initiative: 10)
		{ID: 56, Name: "Козырев Александр", TeamID: 18, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "sword", Shield: "shield", Height: 166, Weight: 73, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},               // Боец (было Stamina: 8, Initiative: 8)
		{ID: 83, Name: "Поташный Климентий", TeamID: 18, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 169, Weight: 75, CountOfAbility: 3, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
		{ID: 88, Name: "Коваленко Павел", TeamID: 18, IsActive: true, RoleID: 1, HP: 100, Stamina: 9, AttackMin: 15, AttackMax: 22, Defense: 8, Initiative: 9, Wrestling: 10, Attack: 15, Weapon: "axe", Shield: "buckler", Height: 183, Weight: 90, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},                 // Убийца (было Stamina: 12, Initiative: 12)
		{ID: 171, Name: "Воронов Ричард", TeamID: 18, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "falchion", Shield: "shield", Height: 178, Weight: 83, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},             // Борец (было Stamina: 10, Initiative: 10)

		// TeamID 19: Межевой рыцарь
		{ID: 89, Name: "Басов Захар", TeamID: 19, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "shield", Height: 171, Weight: 77, CountOfAbility: 3, ImageURL: "./static/characters/default.png"},      // Боец (было Stamina: 10, Initiative: 10)
		{ID: 172, Name: "Василенко Михаил", TeamID: 19, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "two_handed_halberd", Shield: "", Height: 174, Weight: 79, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},      // Боец (было Stamina: 8, Initiative: 8)
		{ID: 173, Name: "Безуглый Георгий", TeamID: 19, IsActive: true, RoleID: 4, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 6, Initiative: 8, Wrestling: 12, Attack: 12, Weapon: "axe", Shield: "buckler", Height: 181, Weight: 87, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},             // Борец (было Stamina: 10, Initiative: 10)
		{ID: 174, Name: "Платонов Игорь", TeamID: 19, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "sword", Shield: "shield", Height: 176, Weight: 81, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},             // Боец (было Stamina: 10, Initiative: 10)
		{ID: 175, Name: "Тищенко Алексей", TeamID: 19, IsActive: true, RoleID: 2, HP: 100, Stamina: 6, AttackMin: 10, AttackMax: 15, Defense: 8, Initiative: 6, Wrestling: 8, Attack: 10, Weapon: "falchion", Shield: "", Height: 179, Weight: 84, CountOfAbility: 2, ImageURL: "./static/characters/default.png"},                 // Боец (было Stamina: 8, Initiative: 8)
		{ID: 176, Name: "Никитин Богдан", TeamID: 19, IsActive: true, RoleID: 2, HP: 100, Stamina: 8, AttackMin: 12, AttackMax: 18, Defense: 10, Initiative: 8, Wrestling: 10, Attack: 12, Weapon: "two_handed_sword", Shield: "buckler", Height: 182, Weight: 88, CountOfAbility: 2, ImageURL: "./static/characters/default.png"}, // Боец (было Stamina: 10, Initiative: 10)
	}, nil
}

func (m *MockDatabase) GetRoleConfig() (map[string]Role, error) {
	return map[string]Role{
		"0": {
			Name: "Танк",
			ID:   "0",
		},
		"1": {
			Name: "Убийца",
			ID:   "1",
		},
		"2": {
			Name: "Боец",
			ID:   "2",
		},
		"3": {
			Name: "Поддержка",
			ID:   "3",
		},
		"4": {
			Name: "Борец",
			ID:   "4",
		},
	}, nil
}

func (m *MockDatabase) GetTeamsConfig() (map[int]TeamConfig, error) {
	return map[int]TeamConfig{
		1:  {ID: 1, Name: "Партизан Два", IconURL: "./static/teams/partizan_dva.png", Description: "Вторая команда партизан, стойкие и выносливые бойцы."},
		2:  {ID: 2, Name: "Юг", IconURL: "./static/teams/south.png", Description: "Команда южных земель, известная своей тактикой."},
		3:  {ID: 3, Name: "НСК", IconURL: "./static/teams/nsk.png", Description: "Новосибирские бойцы, сильные и решительные."},
		6:  {ID: 6, Name: "Старые Друзья", IconURL: "./static/teams/old_friends.png", Description: "Ветераны, проверенные временем."},
		7:  {ID: 7, Name: "Черная Земля", IconURL: "./static/teams/black_land.png", Description: "Таинственные воины темных земель."},
		11: {ID: 11, Name: "Партизан", IconURL: "./static/teams/partizan.png", Description: "Первая команда партизан, мастера скрытности."},
		12: {ID: 12, Name: "Злой дух Ямбуя", IconURL: "./static/teams/yambuya_spirit.png", Description: "Мистические воины севера."},
		16: {ID: 16, Name: "Медвежья пядь", IconURL: "./static/teams/bear_span.png", Description: "Сильные, как медведи, бойцы."},

		4:  {ID: 4, Name: "Школа ИСБ Байард", IconURL: "./static/teams/isb_bayard.png", Description: "Ученики школы боевых искусств Байард."},
		5:  {ID: 5, Name: "НРБ", IconURL: "./static/teams/nrb.png", Description: "Непреклонные рыцари битвы."},
		8:  {ID: 8, Name: "Vivus Ferro", IconURL: "./static/teams/vivus_ferro.png", Description: "Живые клинки, мастера оружия."},
		9:  {ID: 9, Name: "Молодые львы", IconURL: "./static/teams/young_lions.png", Description: "Юные и амбициозные бойцы."},
		10: {ID: 10, Name: "Ганза", IconURL: "./static/teams/hansa.png", Description: "Союз торговцев и воинов."},
		13: {ID: 13, Name: "Мальтийский крест", IconURL: "./static/teams/maltese_cross.png", Description: "Рыцари ордена, верные клятве."},
		14: {ID: 14, Name: "Урфин Джус", IconURL: "./static/teams/urfin_jus.png", Description: "Команда загадочных мастеров."},
		15: {ID: 15, Name: "Высшая Школа ИСБ Санкт-Петербург", IconURL: "./static/teams/isb_spb.png", Description: "Элита школы боевых искусств СПб."},
		17: {ID: 17, Name: "Байард", IconURL: "./static/teams/bayard.png", Description: "Рыцари чести и доблести."},
		18: {ID: 18, Name: "RaubRitter", IconURL: "./static/teams/raubritter.png", Description: "Разбойные рыцари, мастера боя."},
		19: {ID: 19, Name: "Межевой рыцарь", IconURL: "./static/teams/border_knight.png", Description: "Стражи границ и традиций."},
	}, nil
}

// NewMockDatabase создаёт новый экземпляр MockDatabase
func NewMockDatabase() Database {
	return &MockDatabase{}
}

func (m *MockDatabase) SetUser(refreshToken string, user User) error {
	mutex.Lock()
	users[user.Email] = user
	if refreshToken != "" {
		usersWithRefresh[refreshToken] = user
	}
	mutex.Unlock()
	return nil
}

func (m *MockDatabase) GetUserByEmail(email string) (User, error) {
	mutex.Lock()
	user, exists := users[email]
	mutex.Unlock()
	if exists {
		return user, nil
	}
	return User{}, nil
}

func (m *MockDatabase) GetUserByRefresh(token string) (User, error) {
	mutex.Lock()
	user, exists := usersWithRefresh[token]
	mutex.Unlock()
	if exists {
		return user, nil
	}
	return User{}, nil
}

func (m *MockDatabase) GetRoom(roomID string) (*Game, error) {
	mutex.Lock()
	game, exists := rooms[roomID]
	mutex.Unlock()
	if exists {
		return game, nil
	}
	return nil, nil
}

func (m *MockDatabase) SetRoom(game *Game) error {
	mutex.Lock()
	rooms[game.GameSessionId] = game
	mutex.Unlock()
	return nil
}

var rooms = make(map[string]*Game)
var users = make(map[string]User)
var usersWithRefresh = make(map[string]User)
