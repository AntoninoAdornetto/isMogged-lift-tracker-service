CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "lifter" varchar NOT NULL,
  "birth_date" date NOT NULL,
  "weight" int NOT NULL,
  "start_date" date NOT NULL
);

CREATE TABLE "muscle_groups" (
  "id" bigserial PRIMARY KEY,
  "group_name" varchar UNIQUE NOT NULL 
);

CREATE TABLE "exersise" (
  "id" bigserial PRIMARY KEY,
  "exersise_name" varchar NOT NULL UNIQUE,
  "muscle_group" varchar NOT NULL REFERENCES muscle_groups(group_name) ON DELETE CASCADE
);

CREATE TABLE "lift" (
  "id" bigserial PRIMARY KEY,
  "exersise_name" varchar NOT NULL REFERENCES exersise(exersise_name) ON DELETE CASCADE,
  "weight" real NOT NULL,
  "reps" int NOT NULL,
  "date_lifted" timestamp NOT NULL DEFAULT NOW(),
  "user_id" uuid  REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX ON "accounts" ("lifter");

CREATE INDEX ON "exersise" ("exersise_name");

CREATE INDEX ON "lift" ("weight");
