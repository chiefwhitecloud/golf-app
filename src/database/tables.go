package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func CreateTables(db *sql.DB) error {

	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}

	return nil
}

const tableCreationQuery = `

DROP TABLE IF EXISTS hole;
DROP TABLE IF EXISTS course;
DROP TABLE IF EXISTS matchup;
DROP TABLE IF EXISTS pairing;
DROP TABLE IF EXISTS player;
DROP TABLE IF EXISTS captain;
DROP TABLE IF EXISTS game;

CREATE TABLE game
(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL
);

CREATE TABLE captain
(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
game_id INT REFERENCES game(id)
);

CREATE TABLE course (
 id serial PRIMARY KEY,
 name VARCHAR (50) UNIQUE NOT NULL,
 game_id INT REFERENCES game(id)
);

CREATE TABLE hole (
  id serial PRIMARY KEY,
	number INT NOT NULL,
  par INT NOT NULL,
	yards INT NOT NULL,
  course_id INT REFERENCES course(id)
);

CREATE TABLE player (
 	id serial PRIMARY KEY,
 	name VARCHAR (50) NOT NULL
);

CREATE TABLE pairing (
 	id serial PRIMARY KEY,
 	player1_id INT REFERENCES player(id),
  player2_id INT REFERENCES player(id),
	captain_id INT REFERENCES captain(id)
);

CREATE TABLE matchup (
 	id serial PRIMARY KEY,
 	name VARCHAR (50) NOT NULL,
 	pairing1_id INT REFERENCES pairing(id),
  pairing2_id INT REFERENCES pairing(id),
 	game_id INT REFERENCES game(id)
);
`
