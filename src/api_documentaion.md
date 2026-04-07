# TicTacToe API Documentation

Это руководство описывает эндпоинты для аутентификации, игровых механик и управления профилями пользователей.

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
Ответ:
```
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
Возвращает главную страницу или базовую информацию об API.

**Эндпоинт:** `GET /api/`

Аутентификация: Не требуется

### Создать игру с ИИ
Создает новую игру против компьютерного соперника.

**Эндпоинт:** `POST /api/game/ai`

Аутентификация: Требуется

Ответ:

```json
{
    "game_id": "7bfc36fe-498a-49be-9ff4-3ada92f5ea0e",
    "message": "The game was created",
    "created_at": "2026-04-07T14:11:21.652615246Z"
}
```