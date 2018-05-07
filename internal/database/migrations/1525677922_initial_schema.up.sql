BEGIN;

CREATE TABLE users (
	id uuid NOT NULL PRIMARY KEY,
	username text NOT NULL,
	email text NOT NULL,
	address text NOT NULL
);


COMMIT;
