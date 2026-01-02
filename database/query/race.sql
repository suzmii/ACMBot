-- name: CreateRaceRaw :exec
insert into races (races) values ($1);

-- name: GetLastRaceRaw :one
select races from races order by created_at desc limit 1;

-- name: GetLastRaceCreatedAt :one
select created_at from races order by created_at desc limit 1;