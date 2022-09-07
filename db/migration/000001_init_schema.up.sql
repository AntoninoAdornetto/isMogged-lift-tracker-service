CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "lifter" varchar NOT NULL,
  "age" int NOT NULL,
  "weight" int NOT NULL,
  "start_date" date NOT NULL
);

CREATE TABLE "muscle_groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL 
);

CREATE TABLE "exersise" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL UNIQUE,
  "muscle_group" varchar NOT NULL REFERENCES muscle_groups(name) ON DELETE CASCADE
);

CREATE TABLE "lift" (
  "id" bigserial PRIMARY KEY,
  "exersise" varchar NOT NULL REFERENCES exersise(name) ON DELETE CASCADE,
  "weight" real NOT NULL,
  "reps" int NOT NULL,
  "muscle_group" varchar NOT NULL REFERENCES muscle_groups(name) ON DELETE CASCADE,
  "date_lifted" timestamp NOT NULL DEFAULT NOW(),
  "user_id" uuid  REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX ON "accounts" ("lifter");

CREATE INDEX ON "exersise" ("name");

CREATE INDEX ON "lift" ("weight");

--ALTER TABLE "exersise" ADD FOREIGN KEY ("muscle_group") REFERENCES "muscle_groups" ("name") ON DELETE CASCADE;

--ALTER TABLE "lift" ADD FOREIGN KEY ("exersise") REFERENCES "exersise" ("name") ON DELETE CASCADE;

--ALTER TABLE "lift" ADD FOREIGN KEY ("muscle_group") REFERENCES "muscle_groups" ("name") ON DELETE CASCADE;

--ALTER TABLE "lift" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
