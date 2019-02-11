package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type pairing struct {
	ID        int
	Player1ID int
	Player2ID int
	CaptainID int
}

func (p *pairing) createPairing(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO pairing(player1_id, player2_id, captain_id) VALUES($1, $2, $3) RETURNING id",
		p.Player1ID, p.Player2ID, p.CaptainID).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}
