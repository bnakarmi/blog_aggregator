-- name: GetFeedFollows :many
SELECT * FROM FEED_FOLLOWS WHERE USER_ID = $1;

-- name: CreateFeedFollow :one
INSERT INTO FEED_FOLLOWS(ID, CREATED_AT, UPDATED_AT, FEED_ID, USER_ID)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM FEED_FOLLOWS WHERE ID = $1;

