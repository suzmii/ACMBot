CREATE TABLE IF NOT EXISTS "atcoder_users" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"username" VARCHAR(255) NOT NULL,
	"avatar_url" VARCHAR(255) NOT NULL,
	"rank" VARCHAR(255) NOT NULL,
	"rating" INTEGER NOT NULL DEFAULT 0,
	"highest_rating" INTEGER NOT NULL DEFAULT 0,
	"promotion_message" VARCHAR(255) NOT NULL DEFAULT '',
	"submission_statistics" JSONB NOT NULL DEFAULT '{}'::jsonb,
	PRIMARY KEY("id")
);

CREATE INDEX "idx_atcoder_users__username"
ON "atcoder_users" ("username");


CREATE TABLE IF NOT EXISTS "atcoder_submissions" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"user_id" BIGINT NOT NULL,
	"submission_id" BIGINT NOT NULL,
	"problem_id" VARCHAR(255) NOT NULL,
	"point" DOUBLE PRECISION NOT NULL,
	"status" VARCHAR(255) NOT NULL,
	"at" TIMESTAMPTZ NOT NULL,
	PRIMARY KEY("id")
);

CREATE INDEX "idx_atcoder_submissions__at"
ON "atcoder_submissions" ("at");

ALTER TABLE "atcoder_submissions"
ADD FOREIGN KEY("user_id") REFERENCES "atcoder_users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
