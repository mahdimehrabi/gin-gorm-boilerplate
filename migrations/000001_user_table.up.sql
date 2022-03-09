CREATE TABLE "users" 
("id" BIGSERIAL,
"created_at" TIMESTAMPTZ,
"updated_at" TIMESTAMPTZ,
"deleted_at" TIMESTAMPTZ,
"email" VARCHAR (50) UNIQUE NOT NULL,
"password" VARCHAR (300) NOT NULL,
"full_name" VARCHAR (60) NOT NULL,
"must_logout" BOOLEAN NOT NULL DEFAULT FALSE,
PRIMARY KEY ("id")
);
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
CREATE INDEX "idx_users_created_at" ON "users" ("created_at");