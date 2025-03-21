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

-- name: GetFeedFollowsForUser :many
SELECT 
    feed.id,
    feed.url,
    feed.created_at,
    feed.updated_at,
    feed.name
FROM feed_follows
JOIN feed ON feed_follows.feed_id = feed.id
WHERE feed_follows.user_id = $1;

-- name: CreateFeedFollow :one
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

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows 
WHERE feed_follows.user_id = $1 AND feed_id = (SELECT id FROM feed WHERE url = $2);

-- name: MarkFeedFetched :exec
UPDATE feed
    SET last_fetched_at = $2, updated_at = $3
    WHERE feed.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feed
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

