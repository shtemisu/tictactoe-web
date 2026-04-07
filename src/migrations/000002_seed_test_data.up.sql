INSERT INTO users (id, login, password_hash, created_at, updated_at) VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'bob', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'charlie', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444444', 'admin', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- WAITING
INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', ARRAY[0,0,0,0,0,0,0,0,0], '11111111-1111-1111-1111-111111111111', NULL, '11111111-1111-1111-1111-111111111111', 'waiting', NULL, NOW(), NOW());

-- PLAYING
INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', ARRAY[1,0,0,0,2,0,0,0,0], '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'playing', NULL, NOW(), NOW());

-- GAME OVER
INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('eeeeeeee-eeee-eeee-eeee-eeeeee2eeeeee', ARRAY[1,1,1,2,2,0,0,0,0], '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'game_over', '11111111-1111-1111-1111-111111111111',
    NOW(), NOW());

INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', ARRAY[1,0,0,0,2,0,0,0,0], '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'game_over', '22222222-2222-2222-2222-222222222222', NOW(), NOW());

INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', ARRAY[1,0,0,0,2,0,0,0,0], '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'game_over', '11111111-1111-1111-1111-111111111111', NOW(), NOW());
-- DRAW
INSERT INTO games (id, board, firstPlayer_id, secondPlayer_id, current_turn, status, winner, created_at, updated_at) VALUES
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', ARRAY[1,2,1,1,2,2,2,1,1], '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', NULL, 'draw', NULL, NOW(), NOW());