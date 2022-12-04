CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL UNIQUE,
  "password" VARCHAR NOT NULL,
  "weight" REAL NOT NULL DEFAULT 0.0,
  "body_fat" REAL NOT NULL DEFAULT 0.0,
  "start_date" date NOT NULL DEFAULT NOW() 
);

CREATE TABLE "muscle_group" (
  "id" SMALLSERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "category" (
  "id" SMALLSERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "exercise" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL UNIQUE,
  "muscle_group" VARCHAR NOT NULL REFERENCES muscle_group(name) ON DELETE CASCADE,
  "category" VARCHAR NOT NULL REFERENCES category(name)
);

CREATE TABLE "workout" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "start_time" TIMESTAMP NOT NULL,
  "finish_time" TIMESTAMP NOT NULL
);

CREATE TABLE "lift" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "exercise_name" VARCHAR NOT NULL REFERENCES exercise(name) ON UPDATE CASCADE,
  "weight_lifted" REAL NOT NULL,
  "reps" SMALLINT NOT NULL,
  "user_id" uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  "workout_id" uuid NOT NULL REFERENCES workout(id) ON DELETE CASCADE
);

CREATE TABLE "workout_template" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR NOT NULL,
  "rating" REAL
);

CREATE INDEX ON "accounts" ("id");
CREATE INDEX ON "exercise" ("name");
CREATE INDEX ON "lift" ("weight_lifted");
CREATE INDEX ON "lift" ("reps");
CREATE INDEX ON "workout" ("start_time");
CREATE INDEX ON "workout_template" ("rating");
