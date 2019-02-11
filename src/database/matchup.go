package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type matchup struct {
	ID         int
	Name       string
	Pairing1ID int
	Pairing2ID int
	GameID     int
}

func (m *matchup) createMatchup(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO matchup(name, pairing1_id, pairing2_id, game_id) VALUES($1, $2, $3) RETURNING id",
		m.Name, m.Pairing1ID, m.Pairing2ID, m.GameID).Scan(&m.ID)

	if err != nil {
		return err
	}

	return nil
}
