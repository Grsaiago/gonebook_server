CREATE TABLE contact (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL DEFAULT(gen_random_uuid()),
  "first_name" VARCHAR(20) NOT NULL,
  "last_name" VARCHAR(20) NOT NULL,
  "age" INTEGER NOT NULL,
  "phone" VARCHAR(30) UNIQUE NOT NULL,
  "created_at" DATE NOT NULL DEFAULT (NOW()),
  "updated_at" DATE NOT NULL DEFAULT (NOW())
);
