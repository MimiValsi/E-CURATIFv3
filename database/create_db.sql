SELECT 'CREATE DATABASE ecuratif' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ecuratif')\gexec

\c ecuratif

CREATE TABLE IF NOT EXISTS source (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR,
	code_gmao VARCHAR
);

CREATE TABLE IF NOT EXISTS info (
	id SERIAL NOT NULL,
	agent VARCHAR,
	ouvrage VARCHAR,
	priorite INTEGER,
	detail TEXT,
	source_id INTEGER,
	created DATE,
	updated DATE,
	status VARCHAR,
	evenement VARCHAR,
	commentaire VARCHAR,
	echeance VARCHAR,
	entite VARCHAR,
	CONSTRAINT fk_source
	  FOREIGN KEY(source_id) REFERENCES source(id)
);

CREATE TABLE IF NOT EXISTS session (
	token VARCHAR(43) PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  hashed_password VARCHAR(60) NOT NULL,
  created DATE NOT NULL
);

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
