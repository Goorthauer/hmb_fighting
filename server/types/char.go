package types

import (
	"math/rand"
	"time"
)

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
