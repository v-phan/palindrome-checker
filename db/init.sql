
-- Creation of users table
CREATE TABLE IF NOT EXISTS public."user"
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 ),
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    PRIMARY KEY (id)
);