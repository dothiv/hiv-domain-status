DROP TABLE IF EXISTS domain;

CREATE TABLE domain (
	id serial primary key not null,
	name varchar(128),
	created timestamp DEFAULT current_timestamp
);

CREATE INDEX domain__dn_idx ON domain ( name );
