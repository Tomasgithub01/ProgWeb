-- name: CreateGame :one
INSERT INTO game (name, description, image) 
VALUES ($1, $2, $3)
RETURNING id, name, description, image;

-- name: GetGame :one
SELECT * 
FROM game
WHERE id = $1;

-- name: ListGames :many
SELECT *
FROM game
ORDER BY name;

-- name: UpdateGame :exec
UPDATE game 
SET name = $2, description = $3, image = $4
WHERE id = $1;

-- name: Delete :exec
DELETE FROM game
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO user (name, password)
VALUES ($1, $2)
RETURNING id, name, password;

-- name: GetUser :one
SELECT id
FROM user
WHERE name = $1 AND password = $2;

-- name: CreateUserPlaysGame :one
INSERT INTO plays (id_game, id_user)
VALUES ($1, $2)
RETURNING id_game, id_user;

-- name: UpdateState :exec
UPDATE plays
SET state = $3
WHERE id_game = $1 AND id_user = $2;

-- name: UpdateRating :exec
UPDATE plays 
SET rating = $3
WHERE id_game = $1 AND id_user = $2;

-- name: AverageGameRating :one
SELECT avg(COALESCE(rating, 0))
FROM plays
WHERE id_game = $1;

-- name: DeleteUserPlaysGame :exec
DELETE FROM plays
WHERE id_game = $1 AND id_user = $2;