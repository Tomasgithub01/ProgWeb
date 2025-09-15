-- name: CreateGame :one
INSERT INTO game (name, description, image, state, rating) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, image, state, rating;

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
SET name = $2, description = $3, image = $4, state = $5, rating = $6
WHERE id = $1;

-- name: Delete :exec
DELETE FROM game
WHERE id = $1;