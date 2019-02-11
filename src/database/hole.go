package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type hole struct {
	ID       int
	Number   int
	Par      int
	Yards    int
	CourseID int
}

func (h *hole) createHole(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO hole(number, par, yards, course_id) VALUES($1, $2, $3, $4) RETURNING id",
		h.Number, h.Par, h.Yards, h.CourseID).Scan(&h.ID)

	if err != nil {
		return err
	}

	return nil
}
