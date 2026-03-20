-- name: CreateAtcoderUserRaw :one
INSERT INTO atcoder_users (username, avatar_url, rank, rating, highest_rating, promotion_message, submission_statistics)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAtcoderUserByIDRaw :one
SELECT * FROM atcoder_users WHERE id = $1;

-- name: GetAtcoderUserByUsernameRaw :one
SELECT * FROM atcoder_users WHERE username = $1;

-- name: GetAtcoderSubmissionsAfter :many
SELECT * FROM atcoder_submissions WHERE user_id = $1 AND at > $2 ORDER BY id ASC;

-- name: UpdateAtcoderSubmissionStatisticsRaw :exec
UPDATE atcoder_users SET submission_statistics = $2 WHERE id = $1 returning *;

-- name: CreateAtcoderSubmissionsRaw :exec
INSERT INTO atcoder_submissions (user_id, submission_id, problem_id, point, status, at)
SELECT
	unnest($1::bigint[]),
	unnest($2::bigint[]),
	unnest($3::varchar[]),
	unnest($4::double precision[]),
	unnest($5::varchar[]),
	unnest($6::timestamptz[]);
