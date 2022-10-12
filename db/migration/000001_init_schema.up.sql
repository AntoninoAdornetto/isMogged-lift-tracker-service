CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "lifter" varchar NOT NULL,
  "birth_date" date NOT NULL,
  "weight" int NOT NULL,
  "start_date" date NOT NULL DEFAULT NOW() 
);

CREATE TABLE "muscle_groups" (
  "id" bigserial PRIMARY KEY,
  "group_name" varchar UNIQUE NOT NULL 
);

CREATE TABLE "exercise" (
  "id" bigserial PRIMARY KEY,
  "exercise_name" varchar NOT NULL UNIQUE,
  "muscle_group" varchar NOT NULL REFERENCES muscle_groups(group_name) ON DELETE CASCADE
);

CREATE TABLE "set" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TABLE "lift" (
  "id" bigserial PRIMARY KEY,
  "exercise_name" varchar NOT NULL REFERENCES exercise(exercise_name) ON DELETE CASCADE,
  "weight" real NOT NULL,
  "reps" int NOT NULL,
  "date_lifted" timestamp NOT NULL DEFAULT NOW(),
  "user_id" uuid NOT NULL  REFERENCES accounts(id) ON DELETE CASCADE,
  "set_id" uuid NOT NULL REFERENCES set(id) ON DELETE CASCADE
);

CREATE INDEX ON "accounts" ("lifter");

CREATE INDEX ON "exercise" ("exercise_name");

CREATE INDEX ON "lift" ("weight");
