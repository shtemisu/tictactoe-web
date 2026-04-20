CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login TEXT NOT NULL unique,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY,
    board INT[],
    firstPlayer_id UUID,
    secondPlayer_id UUID,
    current_turn UUID,
    status TEXT,
    winner UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_first_player
        FOREIGN KEY (firstPlayer_id) REFERENCES users(id),
    CONSTRAINT fk_second_player
        FOREIGN KEY (secondPlayer_id) REFERENCES users(id),
    CONSTRAINT fk_winner
        FOREIGN KEY (winner) REFERENCES users(id),       

    CONSTRAINT check_players_different
        CHECK (secondPlayer_id IS NULL OR firstPlayer_id != secondPlayer_id)
);

