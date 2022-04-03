CREATE TABLE "users" 
("id" BIGSERIAL,
"created_at" TIMESTAMPTZ,
"updated_at" TIMESTAMPTZ,
"deleted_at" TIMESTAMPTZ,
"email" VARCHAR (50) UNIQUE NOT NULL,
"password" VARCHAR (256) NOT NULL,
"first_name" VARCHAR (60) NOT NULL,
"last_name" VARCHAR (60) NOT NULL,
"devices" JSON ,
"verified_email" BOOLEAN NOT NULL DEFAULT FALSE ,
"last_verify_email_date" TIMESTAMPTZ,
"verify_email_token" VARCHAR (40),
"forgot_password_token" VARCHAR (40),
"last_forgot_email_date" TIMESTAMPTZ,
PRIMARY KEY ("id")
);
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
CREATE INDEX "idx_users_created_at" ON "users" ("created_at");