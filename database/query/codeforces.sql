-- name: CreateCodeforcesUserRaw :one
INSERT INTO codeforces_users (username, avatar_url, friend_num, rating_records, submission_statistics) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetCodeforcesUserByIDRaw :one
SELECT * FROM codeforces_users where id = $1;

-- name: GetCodeforcesUserByUsernameRaw :one
SELECT * FROM codeforces_users where username = $1;

-- name: GetCodeforcesUserLastSubmission :one
SELECT * FROM codeforces_submissions WHERE user_id = $1 ORDER BY id DESC LIMIT 1;

-- name: GetCodeforcesSubmissionsAfter :many
SELECT * FROM codeforces_submissions WHERE user_id = $1 AND at > $2 ORDER BY id ASC;

-- name: UpdateCodeforcesRatingRecordsRaw :exec
UPDATE codeforces_users SET rating_records = $2 WHERE id = $1 returning *;

-- name: UpdateCodeforcesSubmissionStatisticsRaw :exec
UPDATE codeforces_users SET submission_statistics = $2 WHERE id = $1 returning *;

-- name: CreateCodeforcesSubmissionsRaw :exec
INSERT INTO codeforces_submissions (user_id, problem, status, at)
SELECT unnest($1::bigint[]), unnest($2::jsonb[]), unnest($3::varchar[]), unnest($4::timestamptz[]);

