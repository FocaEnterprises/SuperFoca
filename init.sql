DROP TABLE IF EXISTS ranks;
DROP TABLE IF EXISTS titles;

CREATE TABLE titles (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE,
	iq FLOAT NOT NULL UNIQUE
);

CREATE TABLE ranks (
	id SERIAL PRIMARY KEY,
	user_id VARCHAR(255) NOT NULL,
	guild_id VARCHAR(255) NOT NULL,
	iq FLOAT NOT NULL DEFAULT 1,
	title_id INT DEFAULT 1,
	CONSTRAINT user_per_guild UNIQUE (user_id, guild_id),
	CONSTRAINT valid_iq CHECK (iq >= 1),
	FOREIGN KEY (title_id) REFERENCES titles(id)
);

INSERT INTO titles (title, iq) VALUES
	('Ameba', 1.0),
	('Neandertal', 10.0),
	('Mula do PT', 13.0),
	('Gado Bolsonarista', 22.0),
	('Amante do Bolsa Família', 30.0),
	('Filhos do Olavo', 40.0),
	('Estudantes do MOBRAL', 50.0),
	('Investidor de Tigrinho', 60.0),
	('Seguidor do Marçal', 70.0),
	('Aluno do Primo Rico', 80.0),
	('Funcionário Público', 90.0),
	('Humano', 100.0),
	('Bem Nutrido', 110.0),
	('Asiático', 120.0);

DROP FUNCTION IF EXISTS update_ranks;

CREATE OR REPLACE FUNCTION update_ranks()
RETURNS trigger AS
$$
BEGIN
	NEW.title_id = (
		SELECT id
		FROM titles
		WHERE NEW.iq + 1 >= titles.iq
		ORDER BY titles.iq DESC
		LIMIT 1
	);
	return NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS check_update ON ranks;

CREATE TRIGGER check_update
	BEFORE INSERT OR UPDATE ON ranks
	FOR EACH ROW
	EXECUTE FUNCTION update_ranks();
