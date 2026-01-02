-include .env

export

.PHONY: migrateup migratedown

migrateup:
	migrate -path database/migration -database "$(DSN)" -verbose up 1

migratedown:
	migrate -path database/migration -database "$(DSN)" -verbose down 1
