/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE users
(
    id          SERIAL PRIMARY KEY,
    phoneNumber VARCHAR(35)                           NOT NULL,
    fullName    VARCHAR(60)                           NOT NULL,
    password    CHAR(60)                              NOT NULL,
    saltKey     CHAR(36)                              NOT NULL,
    createdAt   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedAt   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT idx_user_phone_number UNIQUE (phoneNumber)
);
