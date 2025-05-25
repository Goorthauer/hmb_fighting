package entities

import (
	"math/rand"
	"time"
)

type Character struct {
	//base
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TeamID         int    `json:"team"`
	RoleID         int    `json:"role"`
	CountOfAbility int    `json:"-"`
	ImageURL       string `json:"imageURL"`
	IsActive       bool   `json:"isActive"`
	//заполняются в бою или перед инициализацией
	Abilities []string `json:"abilities"`
	Effects   []Effect `json:"effects"`
	Position  [2]int   `json:"position"`
	// Снаряжение
	Weapon        string `json:"weapon"`
	Shield        string `json:"shield"`
	IsTitanArmour bool   `json:"IsTitanArmour"` //новая, показывает из титана комплект доспехов у человека или нет.
	// Основные характеристики
	Height     int `json:"height"`
	Weight     int `json:"weight"`
	HP         int `json:"hp"`
	Stamina    int `json:"stamina"`
	Initiative int `json:"initiative"`
	Wrestling  int `json:"wrestling"` //новая, показатель борьбы. Если показатель у персонажа показатель высокий, а у противника низкий, то шанс успеха приема большой, если наоборот - шанс неуспеха приема большой.
	Attack     int `json:"attack"`    //новая, показатель атаки. Нужен для невилирования защиты противника.
	Defense    int `json:"defense"`
	// остальное
	AttackMin int `json:"attackMin"`
	AttackMax int `json:"attackMax"`
}

func (c *Character) PrepareToFight(abilities map[string]Ability) {
	c.Position = [2]int{-1, -1}
	c.SetAbilities(abilities)
	if c.IsTitanArmour { // выглядит тупо, но пока пойдет.
		c.Wrestling += 1
		c.Stamina += 1
		c.Initiative += 1
		c.Defense -= 2
		c.HP -= 5
		if c.HP < 1 {
			c.HP = 1
		}
		if c.Defense < 0 {
			c.Defense = 0
		}
	}
}

func (c *Character) SetAbilities(abilitiesConfig map[string]Ability) {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Преобразуем ключи карты в слайс
	keys := make([]string, 0, len(abilitiesConfig))
	for key := range abilitiesConfig {
		keys = append(keys, key)
	}

	// Очищаем текущие способности персонажа
	c.Abilities = make([]string, 0)

	// Выбираем случайные способности
	for i := 0; i < c.CountOfAbility; i++ {
		if len(keys) == 0 {
			break // Если пул способностей пуст, выходим из цикла
		}

		// Выбираем случайный индекс
		randomIndex := rand.Intn(len(keys))
		// Добавляем выбранную способность в слайс персонажа
		c.Abilities = append(c.Abilities, keys[randomIndex])
		// Удаляем выбранную способность из пула, чтобы избежать дублирования
		keys = append(keys[:randomIndex], keys[randomIndex+1:]...)
	}
}
