package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type course struct {
	ID     int
	Name   string
	GameID int
}

func (c *course) createCourse(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO course(name, game_id) VALUES($1, $2) RETURNING id",
		c.Name, c.GameID).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}
