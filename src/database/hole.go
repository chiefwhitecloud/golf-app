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

func getHoles(db *sql.DB, gameId int) ([]hole, error) {
	rows, err := db.Query(`
		SELECT h.id, h.number, h.par, h.yards, h.course_id
		FROM hole h
		INNER JOIN course c ON
		 	c.id = h.course_id
		INNER JOIN game g ON
			g.id = c.game_id
		WHERE g.id = $1
		ORDER BY h.number ASC`,
		gameId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	holes := []hole{}

	for rows.Next() {
		var h hole
		if err := rows.Scan(&h.ID, &h.Number, &h.Par, &h.Yards, &h.CourseID); err != nil {
			return nil, err
		}
		holes = append(holes, h)
	}

	return holes, nil
}

func getHoleCount(db *sql.DB, gameId int) (int, error) {

	var count int

	row := db.QueryRow(`
		SELECT COUNT(h.id)
		FROM hole h
		INNER JOIN course c
			ON c.id = h.course_id
  	WHERE
			c.game_id = $1;`,
		gameId)

	err := row.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}
