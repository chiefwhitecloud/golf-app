package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type matchup struct {
	ID     int
	Name   string
	GameID int
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

func (m *matchup) getMatchup(db *sql.DB) error {

	row := db.QueryRow(`
		SELECT m.name, m.game_id
		FROM matchup m
  	WHERE
			m.id = $1;`,
		m.ID)

	if err := row.Scan(&m.Name, &m.GameID); err != nil {
		return err
	}

	return nil
}

func getMatchups(db *sql.DB, gameId int) ([]matchup, error) {
	rows, err := db.Query(`
		SELECT m.id matchup_id,
			m.name
		FROM matchup m
  	WHERE
  		m.game_id = $1;`,
		gameId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matchups := []matchup{}

	for rows.Next() {
		var m matchup
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		matchups = append(matchups, m)
	}

	return matchups, nil
}
