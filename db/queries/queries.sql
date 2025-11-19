-- name: CreateGame :one
INSERT INTO game (name, description, image, link) 
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, image, link;

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
SET name = $2, description = $3, image = $4 ,  link = $5
WHERE id = $1;

-- name: Delete :exec
DELETE FROM game
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES ($1, $2)
RETURNING id, name, password;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByName :one
SELECT *
FROM users
WHERE name = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, password = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY name;

-- name: CreateUserPlaysGame :one
INSERT INTO plays (id_game, id_user, state, rating)
VALUES ($1, $2, $3, $4)
RETURNING id_game, id_user, state, rating;

-- name: GetUserPlaysGame :one
SELECT *
FROM plays
WHERE id_game = $1 AND id_user = $2;

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

-- name: UpdateUserPlaysGame :exec
UPDATE plays
SET state = $3, rating = $4
WHERE id_game = $1 AND id_user = $2;

-- name: ListUserPlaysGames :many
SELECT *
FROM plays
ORDER BY id_game, id_user;