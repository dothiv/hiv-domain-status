DROP TABLE IF EXISTS domain;

CREATE TABLE domain (
	id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	name varchar(128),
	created timestamp DEFAULT current_timestamp
);

CREATE INDEX domain__dn_idx ON domain ( name );
