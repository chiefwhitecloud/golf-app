package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type matchup struct {
	ID         int
	Name       string
	GameID     int
}

func (m *matchup) createMatchup(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO matchup(name, game_id) VALUES($1, $2) RETURNING id",
		m.Name, m.GameID).Scan(&m.ID)

	if err != nil {
		return err
	}

	return nil
}
