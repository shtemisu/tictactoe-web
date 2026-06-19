# TicTacToe API Documentation

Это руководство описывает эндпоинты для аутентификации, игровых механик и управления профилями пользователей.

## Запуск в Docker-контейнерах
Для запуска в директории src введите в терминнал следующую команду:

`docker-compose up --build`

или

`docker-compose up -d`

## Запуск фронта без сборки

в директории web/
`python3 -m http.Server 5500`

## Base URL
http://localhost:8080

## JWT Аутентификация

Большинство endpoints нуждаются JWT аутентификации. Необходимо добавить accessToken в заголовок "Authorization: Bearer $TOKEN"

## API endpoints

### Authentication

#### Вход в систему
Аутентификация пользователя и получение токенов.

**Эндпоинт:** `POST /api/auth/signin`

**Тело запроса:**
```json
{
    "login": "string",
    "password": "string"
}
```
**Ответ:**
```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```
#### Регистрация
Регистрация нового пользователя.

**Эндпоинт:** `POST /api/auth/signin`

Тело запроса:

```json
{
    "login": "string",
    "password": "string"
}
```
Ответ:
```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```

#### Обновление Access токена
Генерация нового access токена по refresh токену.

**Эндпоинт:** `POST /api/auth/refresh/access`

Тело запроса:

```json
{
    "refresh_token": "string"
}
```
Ответ:

```json
{
    "access_token": "string"
}
```
#### Обновление Refresh токена
Генерация новой пары access и refresh токенов.

**Эндпоинт:** `POST /api/auth/refresh/refresh`

Тело запроса:
```json
{
    "refresh_token": "string"
}
```
Ответ:

```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```

## Управление играми
### Главная страница
Возвращает главную страницу.

**Эндпоинт:** `GET /api/`

Аутентификация: Не требуется

### Создать игру с ИИ
Создает новую игру против компьютерного соперника.

**Эндпоинт:** `POST /api/game/ai`

Аутентификация: Требуется

Ответ:

```json
{
    "game_id": "gameUUID",
    "message": "The game was created",
    "created_at": "2026-04-07T14:11:21.652615246Z"
}
```

### Создать мультиплеерную игру

**Эндпоинт:** `POST /api/game/multiplayer`

Аутентификация: Требуется

Ответ:

```json
{
    "game_id": "gameUUID",
    "message": "The game was created",
    "created_at": "2026-04-07T14:11:21.652615246Z"
}
```

### Подключиться к существующей игре

**Эндпоинт:** `POST /api/game/multiplayer/{gameID}`
UUID игры необходимо вставить в {gameID}

Ответ:
```json
{
    "id": "66ad7b67-3242-4aaf-8abe-28fccfe9edb6",
    "board": [
        [
            0,
            0,
            0
        ],
        [
            0,
            0,
            0
        ],
        [
            0,
            0,
            0
        ]
    ],
    "firstPlayer_ID": "44444444-4444-4444-4444-444444444444",
    "secondPlayer_ID": "22222222-2222-2222-2222-222222222222",
    "status": "playing",
    "winner": "nothing",
    "current_turn": "X"
}
```

### Cделать ход
**Эндпоинт:** `POST /api/game/{gameID}`


Аутентификация: Требуется

Тело запроса: 
```json
{
    "row" : 0,
    "col" : 0,
}
```
После отправки запроса, возвращается уже с ходом сделанным алгоритмом Минимакс

Ответ:
```json
{
    "id": "gameUUID",
    "board": [
        [
            1,
            0,
            0
        ],
        [
            0,
            2,
            0
        ],
        [
            0,
            0,
            0
        ]
    ],
    "firstPlayer_ID": "playerUUID",
    "secondPlayer_ID": "anotherPlayerUUID",
    "status": "playing",
    "winner": "nothing",
    "current_turn": "X"
}
```
### Получение игры по UUID

**Эндпоинт:** `GET /api/game/{gameID}`

Аутентификация: Требуется

Ответ:
```json
{
    "id": "gameUUID",
    "board": [
        [
            1,
            0,
            0
        ],
        [
            0,
            2,
            0
        ],
        [
            0,
            0,
            0
        ]
    ],
    "firstPlayer_ID": "playerUUID",
    "secondPlayer_ID": "anotherPlayerUUID",
    "status": "playing",
    "winner": "nothing",
    "current_turn": "X"
}
```

### Получение доступных для подключения игр

**Эндпоинт:** `GET /api/game/{available}`

Аутентификация: не требуется

Ответ:
```json
{
    "games": [
        {
            "game_id": "88504133-2f76-4404-9a82-b0398fa12564",
            "first_player_id": "44444444-4444-4444-4444-444444444444",
            "status": "waiting"
        },
        {
            "game_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
            "first_player_id": "11111111-1111-1111-1111-111111111111",
            "status": "waiting"
        }
    ],
    "count": 2
}
```

### Получение списка лидеров (leaderboard)

**Эндпоинт:** `GET /api/user/leaderboard/{limit}`

Ответ: 
```json
[
    {
        "id": "11111111-1111-1111-1111-111111111111",
        "winrate": 50
    },
    {
        "id": "22222222-2222-2222-2222-222222222222",
        "winrate": 25
    }
]
```

### Получение списка игр по accessToken

**Эндпоинт:** `GET /api/user/me`

Аутентификация: требуется

Ответ:
```json
{
    "id": "22222222-2222-2222-2222-222222222222",
    "login": "bob",
    "games_played": 6,
    "wins": 1
}
```

### Получение пользователя по UUID

**Эндпоинт:** `GET /api/user/{userID}`

Аутентификация: требуется

Ответ: 
```json
{
    "id": "22222222-2222-2222-2222-222222222222",
    "login": "bob",
    "games_played": 6,
    "wins": 1
}
```