-- name: CreateFeedFollow :one
INSERT INTO feedfollows (id,created_at,update_at,user_id,feed_id)
VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feedfollows WHERE user_id=$1;

-- name: DeleteFeedFollow :exec
DELETE FROM feedfollows WHERE id =$1 AND user_id=$2;