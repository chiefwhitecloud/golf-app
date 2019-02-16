package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type captain struct {
	ID     int
	Name   string
	GameID int
}

func (c *captain) createCaptain(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO captain(name, game_id) VALUES($1, $2) RETURNING id",
		c.Name, c.GameID).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getCaptains(db *sql.DB, gameId int) ([]captain, error) {
	rows, err := db.Query(
    "SELECT id, name, game_id FROM captain WHERE game_id = $1",
    gameId)

  if err != nil {
      return nil, err
  }

  defer rows.Close()

  captains := []captain{}

  for rows.Next() {
      var c captain
      if err := rows.Scan(&c.ID, &c.Name, &c.GameID); err != nil {
          return nil, err
      }
      captains = append(captains, c)
  }

  return captains, nil
}
