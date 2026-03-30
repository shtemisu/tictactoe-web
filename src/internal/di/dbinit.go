package di

const initGameTable = `
CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    board INT[],
    firstPlayer_id TEXT,
    secondPlayer_id TEXT,
    current_turn TEXT,
    status TEXT,
    winner TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)`

const initUserTable = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    login TEXT,
    passwordHash TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);`
