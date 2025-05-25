# hmb_fighting — Руководство по проекту для AI-агентов

> Цель этого документа — дать ИИ-агенту (или новому разработчику) исчерпывающую картину проекта:
> архитектура, домен, игровые правила, протокол взаимодействия клиента и сервера, структура кода,
> известные особенности и «подводные камни». Документ описывает поведение кода **на момент написания**.

---

## 1. Обзор проекта

**hmb_fighting** — многопользовательская пошаговая тактическая браузерная игра («Hearthstone meets HoMM3» —
пошаговая «боевка» в духе героических стратегий, но с жанровым уклоном в **исторический средневековый бой (ИСБ)**
и **борьбу**). Два игрока выбирают реальные исторические команды/персонажей и ведут бой отрядом персонажей
на гексагонально-подобной (квадратной) сетке 16×9.

Стек:

| Слой | Технологии |
|---|---|
| Клиент | Чистый **HTML/CSS/JS (ES-модули, Canvas 2D)**, без фреймворков и сборщиков |
| Сервер | **Go 1.23** (net/http, gorilla/websocket), паттерн Handler → Usecase → Database |
| Хранилище | Абстракция `Database`. Две реализации: **in-memory Mock** (используется по умолчанию) и **PostgreSQL + Redis** (кэш) |
| Миграции | **goose** (SQL-файлы в `server/migrations`) |
| Авторизация | JWT (access/refresh), пароли — bcrypt |
| Развёртывание | Docker + docker-compose (app, postgres, redis) |

**Тематика контента:** персонажи и команды названы в честь реальных клубов ИСБ (исторический средневековый бой)
и их участников (реальные люди с именами и характеристиками). Персонажи с `IsTitanArmour=true` — «гиганты в титановом
доспехе», имеют особые модификаторы.

> ⚠️ **Важно для агентов:** фронтенд обращается к `http://localhost:8080` и `ws://localhost:8080` **жёстко зашитыми
> строками**. Без запущенного на `:8080` бэкенда приложение работать не будет. Есть два HTML-входа — см. §10.

---

## 2. Структура репозитория

```
hmb_fighting/
├── index.html               # Экран №1: логин/регистрация → создание/вход в комнату → выбор команды (устаревший/первичный вход)
├── game.html                # Экран №2: игровой бой (Canvas + панели), основной игровой HTML
├── styles.css               # Все стили (auth-экраны, игровой UI, карточки, модалки)
├── main.exe / server/cmd/main.exe  # Скомпилированный сервер (артефакт, ~24 МБ)
├── go.mod / go.sum          # Go-модуль hmb_fighting, Go 1.23
├── Dockerfile               # Multi-stage: сборка Go → alpine + goose (миграции при старте)
├── docker-compose.yml       # app(:8080) + db(postgres:15) + redis(7)
├── .dockerignore
├── server/                  # ——— ВСЯ серверная часть ———
│   ├── cmd/main.go          # Точка входа: CORS, маршруты HTTP, запуск :8080
│   ├── db/
│   │   ├── database.go      # Интерфейс Database (все методы доступа к данным)
│   │   ├── mockdb.go        # Реализация in-memory (weapons/shields/teams/characters/abilities/roles/users/rooms)
│   │   └── realdb.go        # Реализация PostgreSQL + Redis (кэш конфигов, пользователей, комнат)
│   ├── dtos/dtos.go         # Ответы API: RegisterUserResp, SelectTeamResp
│   ├── entities/            # Доменные структуры и «чистая» игровая логика
│   │   ├── entities.go      # User, TeamConfig, Role, Team, Client, Node (для A*)
│   │   ├── character.go     # Character, PrepareToFight, SetAbilities (случайная выдача приёмов)
│   │   ├── fights.go        # Weapon, Shield, Ability, Effect, Action, OpportunityAttack, Battlelog
│   │   └── game.go          # Room, GameState, A*, дальности, урон, броски борьбы, порядок ходов, победа
│   ├── handlers/handlers.go # HTTP+WS хендлеры, upgrader WS, валидация входных данных
│   ├── jwt/jwt.go           # Генерация/проверка access (15 мин) и refresh (7 дней) токенов HS256
│   ├── jwt/jwt_test.go      # Тесты JWT
│   ├── migrations/          # goose-миграции: схема, справочники, персонажи
│   │   ├── 00001_create_tables.sql
│   │   ├── 00002_insert_initial_data.sql   # оружие, щиты, приёмы, роли, команды
│   │   └── 00003_insert_characters.sql     # ~200 персонажей
│   ├── types/types.go       # Константы: размеры доски, фазы, типы действий, роли, атаки вдогонку
│   ├── usecase/             # Бизнес-логика оркестрации
│   │   ├── usecase.go       # WebSocket-сессия, обработка действий игроков, broadcast
│   │   ├── rooms.go         # CreateRoom, RestartRoom, LeaveRoom, initRoom
│   │   ├── teams.go         # SelectTeam, SetTeam, CheckTeams
│   │   ├── users.go         # RegisterUser, LoginUser, RefreshToken, CheckClient
│   │   ├── usecase_test.go  # Тесты
│   ├── utils/utils.go       # bcrypt hash/check
│   └── validators/validators.go  # Простые проверки входных полей
├── static/                  # Статические ассеты (отдаются не сервером, а из корня!)
│   ├── characters/*.png     # Портреты персонажей (по имени файла/id)
│   ├── teams/*.png          # Гербы команд
│   ├── weapons|shields|abilities|icons|*.png, default.jpg
└── js/                      # ES-модули клиента (импортируются из game.html)
    ├── constants.js         # Canvas, cellWidth/Height=60, адаптивные размеры
    ├── state.js             # Глобальное состояние клиента (module-level переменные + сеттеры)
    ├── utils.js             # findCharacter, addLogEntry
    ├── websocket.js         # connectWebSocket, sendMessage, refreshToken (WS://localhost:8080/ws)
    ├── gameLogic.js         # Клиентский A*, canMove/canAttack, дистанции
    ├── renderCanvas.js      # Вся отрисовка Canvas: сетка, зоны, персонажи, стрелки, анимация, частицы
    ├── renderUI.js          # DOM: карточки персонажей, способности, фаза, лог, шапка хода
    ├── eventHandlers.js     # Мышь/перетаскивание (drag&drop), клики, end turn, start
    └── main.js              # «Старый» неиспользуемый вход (ожидал HTML-разметку roomSelect и пр.)
```

> ⚠️ **`js/main.js`** — рудимент: ссылается на элементы (`roomSelect`, `joinRoomBtn` и т.д.), которых нет ни в
> `index.html`, ни в `game.html`. Ни один HTML его не импортирует. Не использовать как образец входа в игру;
> актуальный вход — `game.html` (модуль-скрипт внутри него).

---

## 3. Запуск и конфигурация

### 3.1 Локальный запуск без Docker (текущий режим разработки)

В `server/cmd/main.go` используется **MockDatabase** (строки закомментированы — создание Postgres):

```go
//database, err := db.NewPostgresDatabase() // ЗАКОММЕНТИРОВАНО
database := db.NewMockDatabase()
```

```bash
# из корня репозитория
go run ./server/cmd/main.go        # сервер на :8080
```

- Все данные (оружие/персонажи/команды) захардкожены в `mockdb.go` и идентичны миграциям.
- Пользователи и комнаты хранятся в package-level map'ах `mockdb.go`: `users`, `usersWithRefresh`, `rooms`.
- `index.html`, `game.html`, `styles.css`, `js/`, `static/` открываются **файлом** (`file://`) либо простым
  статик-сервером. **Сервер не раздаёт статику** (нет http.FileServer) — клиент живёт отдельно от Go-бэкенда.

### 3.2 Запуск через Docker (полноценный стек)

```bash
docker-compose up --build
```

- `app` (порт 8080): выполняет `goose -dir /app/migrations postgres "$DATABASE_URL" up && ./hmb_fighting`.
  Миграции применяются **при каждом старте** контейнера.
- `db`: PostgreSQL 15 (`user/password/hmb_fighting`, порт 5432).
- `redis`: Redis 7 с паролем `redispass` (опционально; в коде `REDIS_URL` без пароля = `redis://redis:6379/0`,
  тогда `requirepass` в контейнере помешает подключению — см. §11.6).

Ожидаемые переменные окружения для `realdb.go`: `DATABASE_URL`, `REDIS_URL`.

### 3.3 Сборка бинарника

```bash
go build -o hmb_fighting ./server/cmd/main.go
```

---

## 4. Игровой домен и правила

### 4.1 Доска и координаты

- Доска: `[16][9]` (VerticalSize=16, HorizontalSize=9). **X — вертикаль (0..15), Y — горизонталь (0..8)**.
- Значение ячейки: `-1` = пусто, иначе `charID`.
- Начало хода/«своя половина»: **TeamID 0 → X < 8**, **TeamID 1 → X >= 8**. В `game.html`/Canvas это зоны
  зелёная/красная (заливка клеток).
- Соседи — ортогональные 4 направления (не диагональ) для A*/перемещения.
- Расстояния для атак/приёмов — **Chebyshev** (максимум из |dx|,|dy|).

### 4.2 Фазы игры (`server/types/types.go`)

```go
GamePhasePickTeam = "pick_team"   // игроки выбирают команды (до подключения обоих в WS)
GamePhaseSetup    = "setup"       // расстановка персонажей (минимум 5 на команду)
GamePhaseMove     = "move"        // ход персонажа: перемещение на Stamina клеток (или атака)
GamePhaseAction   = "action"      // после перемещения: одна атака/приём
GamePhaseFinished = "finished"    // игра окончена (Winner != -1)
```

Типы действий клиента:

```go
ActionPlace    = "place"     // расстановка {characterID, position}
ActionStart    = "start"     // начать бой (оба игрока расставили ≥5)
ActionMove     = "move"      // перемещение {characterID, position}
ActionAttack   = "attack"    // {characterID, targetID}
ActionAbility  = "ability"   // {characterID, targetID, ability}
ActionEndTurn  = "end_turn"  // {clientID}
```

Каждое `Action` в WS-сообщении дополнительно содержит `clientID`; сервер **сверяет** его с `claims.ClientID`
(см. `HandleWebSocket` в `usecase.go`) и отбрасывает чужие действия.

### 4.3 Характеристики персонажа (`entities/character.go`)

| Поле | Смысл |
|---|---|
| `Height`, `Weight` | Влияют на броски борьбы (модификатор `(Δрост+Δвес)/10*5`) |
| `HP` | Здоровье (базово 100) |
| `Stamina` | Очки движения за ход |
| `Initiative` | Порядок ходов (по убыванию) |
| `Wrestling` | «Борьба»: чем выше, тем выше шанс успеха броска/приёма и ниже шанс «уронить себя» |
| `Attack` / `Defense` | «Нивелируют» друг друга при расчёте урона |
| `AttackMin`/`AttackMax` | Диапазон базового урона |
| `Weapon`, `Shield` | Именованные конфиги (бонусы к атаке/борьбе/защите) |
| `IsTitanArmour` | Персонаж в «титановом» доспехе — особые правки (см. ниже) |
| `CountOfAbility` | Сколько случайных приёмов выдаётся персонажу на бой |
| `RoleID` | 0 Танк, 1 Убийца, 2 Боец, 3 Поддержка, 4 Борец |
| `IsActive` | Доступен ли персонаж для выбора команды (в mock почти у всех true; у части false) |
| `Abilities` | Набор выданных приёмов (случайный) |
| `Effects` | Активные эффекты (сейчас не выдаются, но учтены в расчётах защиты) |
| `Position` | `[-1,-1]` = не размещён/мёртв-не-на-доске |

### 4.4 Экипировка и конфиги

Оружие (`GetWeapons`): `falchion`(Фальшион, дист.1), `axe`(Топор, 1), `two_handed_sword`(Двуручный меч, 2),
`two_handed_halberd`(Алебарда, 2), `sword`(Меч, 1). Бонусы: `attackBonus`, `grappleBonus`.

Щиты (`GetShields`): `buckler`(Баклер: def+1, atk+1, gr+1), `shield`(Тарч: def+2, atk+1), `tower`(Ростовой щит: def+3, gr−1).

Приёмы/способности (`GetAbilities`) — 13 борцовских техник (Подхват, Зацеп, Высед, Передняя/Задняя подножка,
Внешний/Внутренний зацеп, Бросок через плечо/спину, Двойной/Одинарный захват ног, броски через захват руки…).
Тип у всех `wrestle`, Range=1. Выдаются **случайно без повторов** в `SetAbilities`.

### 4.5 «Титановый доспех» (`PrepareToFight`)

Если `IsTitanArmour`:
- `Wrestling +1`, `Stamina +1`, `Initiative +1`
- `Defense -2` (не ниже 0), `HP -5` (не ниже 1)

### 4.6 Команды

Команда — «настоящий» клуб ИСБ с уникальным `ID` (в БД это ключ `teams.id`, в игре появляется как `realTeamID`).
Список команд см. в `00002_insert_initial_data.sql`/`GetTeams` (19 команд). Уникальные ID: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19.

Каждая команда имеет ростер персонажей (обычно 7–10, активных обычно меньше). **В бою игрок использует только
персонажей своей команды** (все активные персонажи выбранной команды попадают в его отряд, но на доску выставляется
минимум 5).

---

## 5. Жизненный цикл игры (end-to-end)

### Шаг 1. Авторизация (`index.html`)

1. `POST /register` `{name,email,password}` или `POST /login` `{email,password}`.
2. Сервер: bcrypt-хэш пароля, `ClientID = generateClientID()` = `"<unixnano>-<rand>"`.
3. Ответ `{accessToken, refreshToken, clientID}` сохраняется в `localStorage`.
4. JWT: access TTL 15 мин (`tokenType=access`), refresh TTL 7 дней (`tokenType=refresh`), роль при выдаче всегда
   `spectator` (см. §11.3).

### Шаг 2. Комната

1. `POST /create-room` `{accessToken}` → `{roomID}` (uuid, он же `GameSessionId`). Создатель становится `Players[0]`.
2. Второй игрок вводит чужой `roomID` в `index.html` (join) — комната не «ищется», сервер просто добавляет участника
   при `SetTeam`/WS-подключении.

### Шаг 3. Выбор команды

`POST /select-team` (без токена!) → `{availableTeams: map[realTeamID]TeamConfig, characters: map[realTeamID][]Character}`
(возвращаются только `IsActive=true` персонажи). Затем `POST /set-team` `{roomID, realTeamID, accessToken}`.

В `SetTeam`:
- Проверяется токен, наличие комнаты/команды/персонажей.
- Запрещено выбирать команду, уже выбранную одним из игроков (`game.TeamsConfig`).
- Игроку присваивается **боевой** `teamID` 0 или 1 (`Players[0]→0`, `Players[1]→1`; если второй слот свободен и роль
  `spectator` — он становится игроком `Players[1]=clientID`, `teamID=1`).
- `TeamID` персонажей заменяется на боевой (0/1), применяется бонус щита к `Defense`, `PrepareToFight`.
- Когда оба игрока в `Players` — инициализируется порядок ходов и фаза → `setup`.

### Шаг 4. Подключение к бою (`game.html`)

`connectWebSocket(room, false, cb)` → `ws://localhost:8080/ws?room=<roomID>&accessToken=<accessToken>`.

Логика распределения слотов в `HandleWebSocket`:
- `Players[0]==clientID` → `TeamID=0`;
- иначе если `len(Players)<2 && role==spectator` → становится игроком слот 1;
- `Players[1]==clientID` → `TeamID=1`;
- иначе → **спектатор** (`TeamID=-1`).

Сервер сразу рассылает полное `GameState` всем подключённым (см. §6).

### Шаг 5. Расстановка (`setup`)

- Клиент тянет карточки персонажей (`#characterCards`, только неразмещённые своей команды) на доску.
- Сервер валидирует в `handleSetupPhase` (ActionPlace): персонаж жив, принадлежит команде клиента, ячейка пуста,
  клетка в своей половине (`team0: x<8`, `team1: x>=8`). Размещение повторно передвигает персонажа.
- «Начать бой» (`ActionStart`): сервер требует **2 игроков** в комнате и ≥5 размещённых живых у каждой команды.
  Неразмещённые при этом «погибают» (`HP=0`, лог «killed due to not being placed»).
- Переход в `move`; `CurrentTurn` — персонаж с максимальной инициативой среди живых (первый в `InitialOrder`).

### Шаг 6. Боевой цикл (`move` → `action`)

Ход персонажа:
1. **`move`**: можно переместиться (одно нажатие/перетаскивание) в любую клетку в пределах `Stamina` шагов A*-путём,
   **либо атаковать/применить приём без движения** (клиент рисует стрелку и «выбегает-возвращается»).
   Если персонаж встал в новую клетку → `Phase=action` (у него остаётся право атаки/приёма).
2. **`action`**: одна атака `attack` ИЛИ один приём `ability`, ИЛИ завершение хода `end_turn`.
3. `end_turn` / любое успешное боевое действие → `NextTurn()`: следующий живой персонаж по кругу `InitialOrder`
   (порядок по инициативе, установлен один раз), при этом тикают эффекты. Проверка победителя: если у команды
   не осталось живых — победа другой команды (`Winner`, `Phase=finished`).

Правила валидации действий на сервере (usecase.go):
- `move`: фаза `move`, клетка на доске и пустая, путь A* ≤ stamina; попутно проверяются **атаки вдогонку**.
- `attack`: фаза `move`/`action`, цель на доске, жива, вражеская, дистанция (Chebyshev) ≤ `weapon.Range`.
- `ability`: фаза `action`, цель валидна (см. `canTarget`), приём `range` достижим, приём **есть в списке**
  `currentChar.Abilities`; после применения приём удаляется из списка (одноразовые).
- Все хендлеры сначала проверяют `currentChar.TeamID == client.TeamID` (иначе действие игнорируется).

### Шаг 7. Конец игры и рестарт

- Клиент при `phase==finished` показывает оверлей («Победа!»/«Поражение!» по `winner == teamID`) и эффект частиц.
- `POST /restart` `{accessToken, roomID}`: сбрасывает `Phase=setup`, `Winner=-1`, доску, HP/эффекты/приёмы
  персонажей, `InitialOrder`, лог (доступно только участникам комнаты).
- `POST /leave-room` `{accessToken, roomID}`: удаляет игрока из `Players`/`Connections`; если комната опустела —
  `Phase=pick_team`.

---

## 6. Протокол WebSocket / GameState

### 6.1 Сообщение от сервера (broadcast `GameState`)

`usecase.broadcastRoomState` шлёт каждому клиенту **персонализированный** объект (`teamID`, `clientID` свои):

```jsonc
{
  "teams": [ /* [0]=Team0 {characters:[...]}, [1]=Team1 {characters:[...]} */ ],
  "winner": -1,
  "currentTurn": 123,
  "phase": "move",               // "pick_team" | "setup" | "move" | "action" | "finished"
  "board": [[-1, ...], ...],     // 16x9, charID или -1
  "teamID": 0,                   // для ПЕРСОНАЛИЗАЦИИ ответа
  "clientID": "1736...-123",
  "gameSessionId": "uuid",
  "weaponsConfig": {...},        // справочник оружия (id → Weapon)
  "abilitiesConfig": {...},      // справочник приёмов
  ".shieldsConfig": {...},
  "teamsConfig": [ /* [0] и [1]: реальные TeamConfig {iconURL, name, ID, description} */ ],
  "battlelog": [ {"Time":"15:04:05","Action":"..."} ]
}
```

> **Внимание:** `teams` в `GameState` — **массив из 2 элементов** (`[2]Team`), при этом если команда не выбрана,
> элемент нулевой (`characters: null`). Клиентская `findCharacter` умеет перебирать и массив, и объект.

JSON-ключи соответствуют Go-tag'ам: `"hp"`, `"stamina"`, `"position"`, `"abilities"`, `"team"` (атрибут TeamID у
Character тегирован `json:"team"`!), `"weapon"`, `"shield"`, `"attackMin"`, `"attackMax"`, `"wrestling"`, `"initiative"` и т.д.

### 6.2 Сообщение от клиента (Action)

```jsonc
{ "type":"place",  "clientID":"...", "characterID": 4,  "position": [5, 2] }
{ "type":"start",  "clientID":"..." }
{ "type":"move",   "clientID":"...", "characterID": 4,  "position": [6, 2] }
{ "type":"attack", "clientID":"...", "characterID": 4,  "targetID": 71 }
{ "type":"ability","clientID":"...", "characterID": 4,  "targetID": 71, "ability":"yama_arashi" }
{ "type":"end_turn","clientID":"..." }
```

Сервер отвечает **новым полным GameState** (broadcast) после каждого принятого действия и после подключения/отключения.

---

## 7. HTTP API (маршруты `server/cmd/main.go`)

Все CORS-методы: `GET, POST, OPTIONS`. CORS заголовки: `Access-Control-Allow-Origin: *`.

| Метод | Путь | Тело (JSON) | Ответ (успех) | Комментарий |
|---|---|---|---|---|
| POST | `/register` | `{name,email,password}` | `{accessToken,refreshToken,clientID}` | 400 при ошибке валидации |
| POST | `/login` | `{email,password}` | `{accessToken,refreshToken,clientID}` | пароль bcrypt |
| POST | `/refresh` | `{refreshToken}` | `{accessToken,refreshToken,clientID}` | 401 если просрочен/не refresh |
| POST | `/check-client` | `{clientID,accessToken}` | `{success,valid}` | сверка clientID и claims |
| POST | `/create-room` | `{accessToken}` | `{success,roomID}` | uuid комнаты |
| POST | `/select-team` | — (пусто) | `{availableTeams,characters}` | **без токена** |
| POST | `/set-team` | `{roomID,realTeamID,accessToken}` | `{success,message}` | realTeamID — ID реальной команды |
| POST | `/check-teams` | `{roomID,accessToken}` | `{allTeamsSelected}` | true, когда фаза уже `setup` |
| POST | `/restart` | `{accessToken,roomID}` | `{success,message}` | 403 посторонним |
| POST | `/leave-room` | `{accessToken,roomID}` | `{success,message}` | 403 если не игрок |
| GET/WSS | `/ws?room=..&accessToken=..` | — | — | WebSocket-апгрейд |
| GET | `/swagger/` | — | Swagger UI | swaggo/http-swagger (без сгенерированных аннотаций — пустой/минимальный) |

Хендлеры всегда декодируют JSON в локальную struct и вызывают `validators.*`, затем `usecase.*`. Статусы ошибок:
`400 Bad Request` (невалидное тело/поля), `500` (внутренние), `401/403` (авторизация/права), `404` (комната/команда).

---

## 8. Ключевые алгоритмы и формулы

### 8.1 Расчёт урона обычной атаки `CalculateDamage` (server)

```
base = rand(AttackMin..AttackMax)
totalDefense = target.Defense + Σ effects.DefenseMod
base += (attacker.Attack − totalDefense)
base += weapon.AttackBonus + shield.AttackBonus
base += CountSurroundingEnemies(target) * 2     // «окружение»
damage = max(0, base)
```

### 8.2 Атака вдогонку (opportunity attacks) — при перемещении мимо врага

В `CheckOpportunityAttacks` для каждого врага, чья «зона угрозы» (соседние 8 клеток вокруг) пересекается с путём
(выход из зоны, либо вход в зону и выход), бросается:

```
tripChance   = clamp(15 + (WrestlingDiff)*2 + pathLen*3 + enemies*5, 5, 90)
attackChance = min(45 + pathLen*2 + enemies*3, 90 - tripChance)
roll = rand(100)
roll < tripChance            → «подсечка»: атакующий валит, цель падает (урон = всё HP цели)
roll < tripChance+attack     → обычная атака CalculateDamage
иначе                        → ничего
```

### 8.3 Бросок/приём (борьба) `ApplyWrestlingMove`

Пять исходов, вероятности зависят от `WrestlingDiff = attacker.Wrestling − target.Wrestling`,
модификаторов роста/веса `mod = (ΔHeight+ΔWeight)/10*5`, окружения `surrounding*5`, бонусов weapon/shield `GrappleBonus`:

```
successChance       = 25 + diff*5 + mod + env + wBonus + sBonus        (clamp 5..90)
partialSuccessChance= 25 + mod + env + wBonus + sBonus                 (clamp ≥5)
nothingChance       = 25
failureChance       = 15 − diff*2 − (mod+env)/2                        (clamp ≥5)  → оба падают (оба HP=0)
totalFailureChance  = 10 − diff*3 − (mod+env)/2                        (clamp ≥5)  → атакующий падает сам
нормализация процентов к сумме 100; бросок rand(100)
```

Исходы:
1. `success` — цель мгновенно повержена (`HP=0`).
2. `partial` — урон `CalculateDamageAfterWrestle` (без бонусов оружия!).
3. `nothing` — промах.
4. `failure` — **оба** падают.
5. `totalFailure` — падает атакующий.

### 8.4 Поиск пути — A*

- Сервер: `FindPath` (openList + closed map, эвристика Манхэттен, соседи 4, вес шага 1, обрезка по stamina,
  возврат пути и списка opportunity-атак).
- Клиент: дублирующая реализация в `gameLogic.js` (используется для подсветки зоны движения и анимации).
- `Node` определён в `entities.go` (сервер) и в `gameLogic.js` (клиент).

### 8.5 Порядок ходов

- `InitTurnOrder()` — вызывается один раз при готовности обеих команд (в `SetTeam`) или при первом
  действии/`NextTurn`, если пуст: сортировка всех живых персонажей по `Initiative` по убыванию (пузырьковая).
- `NextTurn()`: победа при вымирании команды; далее циклический поиск следующего живого по `InitialOrder`;
  фаза всегда `move`; эффекты уменьшаются на 1.

---

## 9. Авторизация и безопасность (JWT)

`server/jwt/jwt.go`:

- `Claims`: `ClientID`, `Email`, `Role` (`"player"`/`"spectator"`), `TokenType`, RegisteredClaims.
- Секрет: **хардкод** `[]byte("your-secret-key")` — для продакшена менять (см. §11.7).
- `GenerateTokenPair(user, role)` — access (15 мин), refresh (7 дней), HS256.
- `ValidateToken` — принимает только `TokenType=="access"` и не истёкшие.
- `RefreshToken` — принимает только `TokenType=="refresh"`, выдаёт новую пару, **сохраняя роль из старого токена**.

Интересный момент: регистрация/логин всегда выдают роль `spectator`, а «повышение» до игрока происходит при
занятии слота комнаты (в `SetTeam`/WS) — проверяется `claims.Role == UserRoleSpectator`, и роль в **самом токене
клиента** при этом не обновляется (используется только для разовой проверки при первом подключении).

---

## 10. Клиентская часть — подробнее

### 10.1 Два входа

1. **`index.html`** (без модулей, обычные `<script>`) — полный флоу: auth → create/join room → select team →
   переход на `game.html`. Функции: `register`, `login`, `createRoom`, `loadAvailableTeams`, `selectTeam`,
   `joinRoom`, `copyRoomId`, `showErrorModal`, `switchAuthTab`, `checkAuth`.
   **Баг/недоделка:** в `selectTeam` fetch идёт на `/set-team`, но `selectTeam(teamID)` вызывается только из
   `confirmTeamSelection`; кнопка `confirm-team-btn` находится в `#team-form`, которая показывается через
   `copyRoomId`/`joinRoom`. UI цепочки — рабочая, но хрупкая.
2. **`game.html`** — модульный вход в бой (import из `js/*`). Хранит данные в замыкании `latestData`, вызывает
   `renderGame(data)` на каждое WS-сообщение. Также содержит обработчики: RAGE QUIT → `/leave-room`,
   «Начать бой» → POST `/ws` (странность — см. §11.5), рестарт через модалку `/restart`.

### 10.2 Состояние и модули

- `state.js`: единственный источник состояния — переменные на уровне модуля; геттеры экспортируются как
  «живые» ссылки (re-export), сеттеры — функции. Импорты в других модулях читают эти ссылки **по значению на момент
  импорта** из-за ре-экспорта — фактически все модули работают с одними и теми же объектами через геттеры-объекты.
- `websocket.js`: переподключение через 1 с при закрытии (и полная очистка localStorage), auto-refresh токена при
  `data.error === 'Invalid token'`.
- `gameLogic.js`: A*, `getGridPosition(event)` (клик→клетка через `getBoundingClientRect`), зоны действия.
- `renderCanvas.js`: рисует всё. Сетка рисуется **всегда поверх заливок** — порядок: clear → grid → зоны/дальности →
  персонажи → драг/анимация.
- `renderUI.js`: карточки, способности, фаза, лог, шапка. `turnOrder`/`roundStarted` сохраняются в
  `localStorage["turnState_<gameSessionId>"]` — клиент пытается «переживать» обновления страницы, но при новой сессии
  сбрасывает.
- `eventHandlers.js`: drag&drop (`dragstart` на карточках, `drop` на canvas) + мышь (`mousedown`…`mouseup`) +
  клики. Обработчики вешаются один раз (флаг `eventListenersSet` в `game.html` / `data-listeners-set` в `main.js`).

### 10.3 Canvas (константы)

- `cellWidth = cellHeight = 60`; доска 16×9 = 960×540 (full-HD) / 840×472 / 720×405 при адаптации.
- Координата X — вертикаль (строки), Y — горизонталь (столбцы) — соответствует серверу `board[x][y]`.

### 10.4 Цветовые и DOM-детали UI

- `.team0` фиолетовый (`rgba(128,0,128,…)`), `.team1` тёмно-красный — на канвасе; в DOM-карточках классы
  `team0`/`team1` для стилей.
- Текущий персонаж обводится жёлтой рамкой на канвасе и имеет класс `.current` в карточке; мёртвые — `.dead`;
  размещённые в setup — `.placed`.
- `#turnText`: «ВАШ ХОД» (синий) / «ХОД ПРОТИВНИКА» (красный) / «РАССТАНОВКА» / «Противник еще выбирает команду».
- Лог боя строится из `data.battlelog` полностью каждый раз (`.innerHTML=''`, затем `addLogEntry`).

---

## 11. Известные баги, несостыковки и «странности» (важно для агентов!)

При правках сверяйтесь с этим списком; вероятно, это точки роста/рефакторинга:

1. **Фронтенд жёстко зашит на `localhost:8080`** (fetch в `index.html`/`game.html`, WS URL в `websocket.js`) —
   для деплоя нужно выносить в конфиг.
2. **Статика не отдаётся сервером**: в `main.go` нет `http.FileServer`. HTML/JS/CSS/static должны раздаваться
   отдельно (или это планировалось — но сейчас так).
3. **`js/main.js` не подключён ни к одному HTML** и ссылается на несуществующие элементы (`roomSelect`,
   `joinRoomBtn`, `roomSelection`). Мёртвый код-кандидат на удаление или реанимацию.
4. **`POST /ws` в `game.html`** (кнопка «Начать бой») отправляет JSON на `http://localhost:8080/ws`, где хендлер —
   WebSocket-upgrade. Запрос, скорее всего, падает/игнорируется; фактически «начать бой» шлётся и по WS как
   `{type:'start'}` (в `eventHandlers.handleStartGame`). Проверить, не дублируется ли логика.
5. **Уровень `Connections` + `Players` в WS-хендлере**: когда подключается spectator с ролью из токена `spectator`,
   он может «занять» слот 1 при `len(Players)<2`. Если это был настоящий второй игрок, уже выбравший команду через
   HTTP (`set-team`), его `Players[1]` уже занят — конфликтов нет, но логика ролей неочевидна и завязана на
   содержимое JWT.
6. **Redis с паролем**: в `docker-compose` redis требует `requirepass redispass`, а `REDIS_URL` задан без пароля —
   при `realdb.go` подключение упадёт. Нужно либо `redis://:redispass@redis:6379/0`, либо убрать `requirepass`.
7. **JWT-секрет захардкожен**, в `usecase.go` есть неиспользуемые импорты/переменные? — нет, но `rand.Seed` в
   `SetAbilities` устарел (Go 1.20+): используйте `rand.New(rand.NewSource(...))`. Аналогично логика `SetAbilities`
   перебирает map (порядок недетерминирован) — это нормально для случайной выдачи.
8. **`canTarget`** требует, чтобы цель была на доске (`isPositionOnBoard(target.Position)`), но НЕ проверяет
   `attackMin/attackMax` и не учитывает `Effects` (эффекты вообще нигде не накладываются — задел на будущее).
9. **Сервер использует глобальный `var mutex sync.Mutex`** в `mockdb.go` для map — это ок для mock, но с
   `PostgresDatabase` комнаты не кэшируются в памяти и полностью перечитываются из БД на каждый ход (см. §11.10).
10. **`SetRoom` (Postgres)** перезаписывает комнату целиком транзакцией (DELETE+INSERT комнат/игроков) на **каждый**
    broadcast-цикл? Нет — SetRoom вызывается при действиях, изменяющих комнату. Но стоит помнить: `GetRoom` не
    восстанавливает `WeaponsConfig/ShieldsConfig/…` из колонок комнаты, а подтягивает их **заново** из БД/кэша
    (поэтому `TeamsConfig` в восстановленной комнате — это полный список всех команд, а не только выбранных;
    в `broadcastRoomState` берутся только `[0]` и `[1]`).
11. **Победа и текущий ход после kill неразмещённых**: в `handleSetupPhase` при `start` неправильно размещённые
    персонажи получают `HP=0` **до** InitTurnOrder — но уже после того как `SetTeam` вызвал `InitTurnOrder` при
    заполнении обоих игроков. `InitTurnOrder` фильтрует только `HP>0`, поэтому фактический порядок корректен;
    однако в `NextTurn` есть ветка сброса фазы в `finished`, если список пуст.
12. **Клиентский `killUnplacedCharacters`** (в `gameLogic.js`) помечает неразмещённых `hp=0` **локально**, а сервер
    делает то же самое при `start` — расхождений нет, но функция на клиенте избыточна и может конфликтовать с
    серверным состоянием при анимации.
13. **`updateAbilityCards`** показывает способности только персонажа с `currentTurn`, т.е. карточки способностей
    «прыгают» при смене хода; в фазе setup текущего персонажа нет — показывается заглушка-стек.
14. **Хрупкий флоу `index.html`:** после «copyRoomId»/«joinRoom» `#room-selection` прячется, а `loadAvailableTeams`
    грузит команды. Кнопка «Copy Room ID» спрятана до создания комнаты. При 401 в `createRoom` — `localStorage.clear()`
    и редирект, но токен не обновляется через `/refresh`.
15. **Два «End Turn»-а**: кнопка `#endTurnBtn` вешается и в `setupEventListeners`, и дублируется в `game.html`
    (только там без `setup`); фактически один обработчик.
16. **`GET /swagger/`** зарегистрирован, но swagger-аннотации не сгенерированы (нет `docs/`) — UI будет пустым.

---

## 12. Таблицы БД (PostgreSQL)

`00001_create_tables.sql` (goose-формат):

- `users(id VARCHAR(50) PK, name, email UNIQUE, password)` + idx email.
- `refresh_tokens(token PK, user_id FK users ON DELETE CASCADE)` + idx.
- `weapons(name PK, display_name, range, is_two_handed, image_url, attack_bonus, grapple_bonus)`.
- `shields(name PK, display_name, defense_bonus, image_url, attack_bonus, grapple_bonus)`.
- `teams(id INT PK, name, icon_url, description)`.
- `roles(id VARCHAR(10) PK, name)`.
- `characters(id INT PK, name, team_id FK, role_id FK, count_of_ability, image_url, is_active, weapon FK,
  shield FK, is_titan_armour, height, weight, hp, stamina, initiative, wrestling, attack, defense, attack_min,
  attack_max)` + idx по team_id/role_id.
- `abilities(name PK, display_name, type, description, range, image_url)`.
- `rooms(game_session_id PK, current_turn, phase, board JSONB, winner, created_at)`.
- `room_teams(game_session_id PK+team_id PK, characters JSONB)` — FK каскад по комнате.
- `room_players(game_session_id PK+team_id PK, client_id)` — FK каскад.

`00002_insert_initial_data.sql`: 5 единиц оружия, 3 щита, 13 приёмов, 5 ролей, 19 команд.
`00003_insert_characters.sql`: ~200 персонажей всех команд (в т.ч. с персональными PNG в `static/characters/` для
команды 11 «Партизан» и нескольких одиночек).

> ⚠️ В `00002` команды добавлены в двух наборах ID: {1,2,3,6,7,11,12,16} и {4,5,8,9,10,13,14,15,17,18,19} —
> порядок вставки отличается от числового; это не ошибка. В `00003` «Злой дух Ямбуя» фигурирует с `team_id=12` и
> комментарием `// TeamID 13: ЗДЯ` в mockdb (комментарий неверен, ID команды 12).

---

## 13. Хранение состояния: Mock vs Postgres (ключевое различие)

**MockDatabase** (активна в dev):
- Все Get-конфиги захардкожены (weapons, shields, abilities, roles, teams, characters).
- `SetUser/GetUserByEmail/GetUserByRefresh/GetRoom/SetRoom` работают с глобальными map.
- `GetRoom` возвращает **ту же ссылку** на `*Room`, что и `SetRoom` сохранил. Именно поэтому мутации комнаты
  (Board, Phase, Teams и пр.) «видны» без повторного чтения. Вся конкурентность — через `room.Mutex`.

**PostgresDatabase** (боевой режим):
- Все Get-конфиги кэшируются в Redis на 24 ч с инвалидацией ключей (`weapons_`, `shields_`, `teams_`,
  `characters_`, `abilities_`, `role_config_`, `user:email:*`, `user:refresh:*`, `room:*`).
- `SetRoom` — полная транзакционная перезапись (upsert `rooms` + delete/insert `room_teams` + delete/insert
  `room_players`), затем инвалидация кэша комнаты. Комната после `GetRoom` — **новый объект**: `Connections`
  всегда пустой (карта WS-соединений не хранится в БД), `Teams`/`Players` подгружаются, конфиги подтягиваются.
  Следовательно при реальной БД **WS-сессии живут только в памяти процесса**, и после перезапуска сервера комнаты
  существуют в БД, но никто к ним не подключён.

---

## 14. Как расширять проект (типовые задачи для агента)

### 14.1 Добавить новую команду/персонажа/оружие
- В mock: `mockdb.go` → `GetTeams()`, `GetCharacters()`, при необходимости `GetWeapons()/GetShields()/GetAbilities()`.
- В БД: новые миграции `00004_*.sql` (INSERT) — не редактировать старые (goose по имени/порядку).
- Добавить картинку в `static/teams/` или `static/characters/` и прописать `imageURL`.
- Проверить, что персонажи `IsActive=true` для доступности в `/select-team`.

### 14.2 Поменять баланс
- Вероятности борьбы — `ApplyWrestlingMove` (game.go); урон — `CalculateDamage*`;
  подсечка/атака вдогонку — `CheckOpportunityAttacks`; расстояния — `DistanceToAttack/Ability`.
- Числа шагов — `Stamina` персонажей (mock + миграции).

### 14.3 Переключить на PostgreSQL
1. `server/cmd/main.go`: раскомментировать `db.NewPostgresDatabase()`, убрать mock.
2. Поднять `docker-compose up -d db redis`, задать `DATABASE_URL`, `REDIS_URL` (см. §3.2 и баг №6 о пароле).
3. Убедиться, что миграции применились (в docker — goose автоматически).

### 14.4 Рефакторинг безопасных мест
- Разнести логику `usecase.go` (WS-цикл, `broadcastRoomState`, `processAction`) по файлам/структурам.
- Вынести `initRoom`-данные из mockdb в отдельные файлы.
- Привести `SetAbilities` к современному rand API.
- Добавить swagger-комментарии и сгенерировать `docs/`.

---

## 15. Глоссарий доменных терминов

| Термин | Значение |
|---|---|
| `realTeamID` | ID «настоящей» команды из справочника (1..19) |
| боевой `teamID` | 0 или 1 — слот игрока в текущем бою (X<8 зона — команда 0) |
| `Players[teamID]` | clientID игрока за слот |
| `ClientID` | уникальный id пользователя, выдаётся при регистрации |
| `GameSessionId`/`roomID` | uuid комнаты (одно и то же) |
| расстановка | фаза setup: 5+ персонажей на свою половину |
| приём/бросок | `ability` типа wrestle: случайная техника из 13, одноразовая |
| подсечка (trip) | атака вдогонку, валящая цель |
| атака вдогонку | бесплатная реакция врага при перемещении мимо него |
| InitialOrder | кольцевой порядок ходов по инициативе (задаётся один раз) |

---

## 16. Чек-лист «что проверить после изменений»

- [ ] `go build ./...` и `go vet ./...` без ошибок; `go test ./...` (есть тесты usecase и jwt).
- [ ] Сервер стартует на :8080 (MockDatabase).
- [ ] `index.html` открывается, регистрация/логин создают токены в localStorage.
- [ ] Создана комната, второй игрок вошёл, оба выбрали разные команды (иначе `set-team` вернёт «Team already selected»).
- [ ] После выбора обеих команд фаза перешла в `setup`; оба игрока расставили ≥5 персонажей в свои зоны.
- [ ] «Начать бой» перевёл в `move`; ходят по очереди, лог пишется, доска синхронна у обоих клиентов.
- [ ] Победа команды корректно фиксируется (`finished`, `winner`), рестарт возвращает в `setup`.
- [ ] При выходе (`RAGE QUIT` / leave-room) комната освобождается корректно.
- [ ] (Если включена БД) — Redis-пароль в `REDIS_URL`, миграции применены, кэш инвалидируется.

---

*Документ подготовлен по состоянию кода в репозитории. При внесении изменений в архитектуру — обновляйте этот файл,
чтобы AI-агенты не опирались на устаревшие данные.*
