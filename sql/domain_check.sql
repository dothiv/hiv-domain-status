DROP TABLE IF EXISTS domain_check;

CREATE TABLE domain_check (
	id SERIAL PRIMARY KEY NOT NULL UNIQUE,
	domain varchar(128) NOT NULL,
	url text NOT NULL,
	status_code integer NOT NULL,
	script_present boolean NOT NULL DEFAULT false,
	iframe_target text DEFAULT NULL,
	iframe_target_ok boolean DEFAULT NULL,
	valid boolean NOT NULL DEFAULT false,
	created timestamp DEFAULT current_timestamp
);

CREATE INDEX domain_check__dn_idx ON domain_check ( domain );
