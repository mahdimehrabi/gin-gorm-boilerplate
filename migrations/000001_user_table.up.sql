CREATE TABLE "users" 
("id" BIGSERIAL,
"created_at" TIMESTAMPTZ,
"updated_at" TIMESTAMPTZ,
"deleted_at" TIMESTAMPTZ,
"email" VARCHAR (50) UNIQUE NOT NULL,
"password" VARCHAR (256) NOT NULL,
"first_name" VARCHAR (60) NOT NULL,
"last_name" VARCHAR (60) NOT NULL,
"refresh_token_secret" VARCHAR (50) NOT NULL,
"devices" JSON (50) NOT NULL,
PRIMARY KEY ("id")
);
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
CREATE INDEX "idx_users_created_at" ON "users" ("created_at");