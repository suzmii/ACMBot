CREATE TABLE IF NOT EXISTS "races" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"races" JSONB NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY("id")
);

CREATE INDEX "idx_races__created_at"
ON "races" ("created_at");


CREATE TABLE IF NOT EXISTS "codeforces_users" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"username" VARCHAR(255) NOT NULL,
	"avatar_url" VARCHAR(255) NOT NULL,
	"friend_num" INTEGER NOT NULL DEFAULT 0,
	"rating_records" JSONB NOT NULL DEFAULT '[]'::jsonb,
	"submission_statistics" JSONB NOT NULL DEFAULT '{}'::jsonb,
	PRIMARY KEY("id")
);

CREATE INDEX "idx_codeforces_users__username"
ON "codeforces_users" ("username");


CREATE TABLE IF NOT EXISTS "codeforces_submissions" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"user_id" BIGINT NOT NULL,
	"problem" JSONB NOT NULL,
	"status" VARCHAR(255) NOT NULL,
	"at" TIMESTAMPTZ NOT NULL,
	PRIMARY KEY("id")
);

CREATE INDEX "idx_codeforces_submissions__at"
ON "codeforces_submissions" ("at");

ALTER TABLE "codeforces_submissions"
ADD FOREIGN KEY("user_id") REFERENCES "codeforces_users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
