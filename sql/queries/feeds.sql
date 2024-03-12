-- name: CreateFeed :one
INSERT INTO FEEDS(ID, CREATED_AT, UPDATED_AT, NAME, URL, USER_ID, LAST_FETCHED_AT)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM FEEDS;

-- name: GetNextFeedsToFetch :many
SELECT * FROM FEEDS ORDER BY LAST_FETCHED_AT ASC NULLS FIRST;

-- name: MarkFeedFetched :exec
UPDATE FEEDS
 SET 
 LAST_FETCHED_AT = CURRENT_TIMESTAMP
, UPDATED_AT = CURRENT_TIMESTAMP
 WHERE ID = $1;
