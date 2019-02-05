package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (g *game) createGame(db *sql.DB, name string) error {
	err := db.QueryRow(
		"INSERT INTO game(name) VALUES($1) RETURNING id",
		g.Name).Scan(&g.ID)

	if err != nil {
		return err
	}

	return nil
}
