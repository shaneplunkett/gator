-- name: CreateFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feed;

-- name: GetFeedByUrl :one
SELECT * FROM feed WHERE url = $1 LIMIT 1;

-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feed.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN feed
    ON  inserted_feed_follow.feed_id = feed.id 
INNER JOIN users
    ON inserted_feed_follow.user_id = users.id;






