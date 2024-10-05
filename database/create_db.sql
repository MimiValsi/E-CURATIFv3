SELECT 'CREATE DATABASE ecuratif' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ecuratif')\gexec

\c ecuratif

CREATE TABLE IF NOT EXISTS source (
	id serial PRIMARY KEY NOT NULL,
	name varchar,
	code_gmao varchar
);

CREATE TABLE IF NOT EXISTS info (
	id serial not null,
	agent varchar,
	ouvrage varchar,
	priorite integer,
	detail text,
	source_id integer,
	created date,
	updated date,
	status varchar,
	evenement varchar,
	commentaire varchar,
	echeance varchar,
	entite varchar,
	CONSTRAINT fk_source
	  FOREIGN KEY(source_id) REFERENCES source(id)
);

CREATE TABLE sessions (
	token VARCHAR(43) PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

INSERT INTO source (name, code_gmao)
VALUES ('ALSACE', 'ALSAC');

INSERT INTO source (name, code_gmao)
VALUES ('AMPERE', 'AMPER');

INSERT INTO source (name, code_gmao)
VALUES ('ARGENTEUIL', 'ARGE6');

INSERT INTO source (name, code_gmao)
VALUES ('BILLANCOURT', 'BILLA');

INSERT INTO source (name, code_gmao)
VALUES ('BOULE', 'BOULE');

INSERT INTO source (name, code_gmao)
VALUES ('COURBEVOIE', 'CZBEV');

INSERT INTO source (name, code_gmao)
VALUES ('DANTON', 'DANT5');

INSERT INTO source (name, code_gmao)
VALUES ('LA BRICHE', 'BRICH');

INSERT INTO source (name, code_gmao)
VALUES ('LEVALLOIS', 'LEVAL');

INSERT INTO source (name, code_gmao)
VALUES ('MENUS', 'MENUS');

INSERT INTO source (name, code_gmao)
VALUES ('NANTERRE', 'NANTE');

INSERT INTO source (name, code_gmao)
VALUES ('NOVION', 'NOVIO');

INSERT INTO source (name, code_gmao)
VALUES ('PUTEAUX', 'PUTEA');

INSERT INTO source (name, code_gmao)
VALUES ('RUEIL', 'RUEIL');

INSERT INTO source (name, code_gmao)
VALUES ('TILLIERS', 'TILLI');

INSERT INTO source (name, code_gmao)
VALUES ('SAINT OUEN', 'SSOU5');
