CREATE TABLE "users" 
("id" bigserial,"created_at" timestamptz,"updated_at"
	timestamptz,"deleted_at" timestamptz,"email" text,
	"full_name" text,PRIMARY KEY ("id"));
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");