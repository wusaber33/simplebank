-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-10-16T09:31:41.289Z

CREATE TABLE "Student" (
  "id" bigsserial PRIMARY KEY
  "name" varchar,
  "password" varchar NOT NULL,
  "sex" varchar NOT NULL,
  "age" bigint NOT NULL,
);

CREATE TABLE "CourseStudent" (
  "course_id" bigint NOT NULL,
  "student_id" bigint NOT NULL,
  "grade" varchar,
)

CREATE TABLE "Course" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "score" bigint NOT NULL,

)

CREATE TABLE "Teacher" (
  "id" bigsserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "title" varchar NOT NULL,
  "salary" bigint NOT NULL,
);

CREATE TABLE "CourseTeacher" (
  "course_id" bigint NOT NULL,
  "teacher_id" bigint NOT NULL,
);



CREATE INDEX ON "accounts" ("owner");

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
