-- name: CreatePost :one
INSERT INTO POSTS(ID, CREATED_AT, UPDATED_AT, TITLE, URL, DESCRIPTION, PUBLISHED_AT, FEED_ID)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT P.ID AS ID
    , P.CREATED_AT AS CREATED_AT
    , P.UPDATED_AT AS UPDATED_AT
    , P.TITLE AS TITLE
    , P.URL AS URL
    , P.DESCRIPTION AS DESCRIPTION
    , P.PUBLISHED_AT AS PUBLISHED_AT
    , P.FEED_ID AS FEED_ID
    FROM POSTS P
LEFT JOIN FEEDS F
    ON F.ID = P.FEED_ID
    AND F.USER_ID = $1
ORDER BY P.UPDATED_AT 
LIMIT $2;

