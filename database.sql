/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( CYCLE INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 99999 CACHE 1 ),
	name character varying(60) NOT NULL,
  phone_number character varying(16) NOT NULL,
  password character varying(64) NOT NULL,
  successful_login bigint NOT NULL DEFAULT 0,
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT users_name_key UNIQUE (name),
  CONSTRAINT users_phone_number UNIQUE (phone_number)
);
