package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type player struct {
	ID   int
	Name string
}

func (p *player) createPlayer(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO player(name) VALUES($1) RETURNING id",
		p.Name).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}
