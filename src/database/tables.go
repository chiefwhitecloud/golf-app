package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func createTables(db *sql.DB) error {

	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}

	return nil
}

const tableCreationQuery = `
DROP TABLE IF EXISTS game

CREATE TABLE game
(
id SERIAL,
name TEXT NOT NULL,
CONSTRAINT game_pkey PRIMARY KEY (id)
)`
