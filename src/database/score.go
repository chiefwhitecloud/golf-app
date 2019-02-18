package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type score struct {
	ID   int
  HoleID int
  Pairing1ID int
  Pairing2ID int
  MatchupID int
  Pairing1Strokes int
  Pairing2Strokes int
	Name string
}

func (s *score) createScore(db *sql.DB) error {
	err := db.QueryRow(
		`INSERT INTO score
    (hole_id, pairing1_id, pairing1_strokes, pairing2_id, pairing2_strokes, matchup_id)
    VALUES($1, $2, $3, $4, $5, $6)
    RETURNING id`,
		s.HoleID, s.Pairing1ID, s.Pairing1Strokes, s.Pairing2ID, s.Pairing2Strokes, s.MatchupID).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func getScores(db *sql.DB, gameId int) ([]score, error) {
	rows, err := db.Query(`
		SELECT
      s.id score_id,
			s.pairing1_id,
      s.pairing1_strokes,
			s.pairing2_id,
      s.pairing2_strokes,
			s.matchup_id
		FROM score s
 			inner join matchup m
    		on s.matchup_id = m.id
    	WHERE
    		m.game_id = $1;`,
    gameId)

  if err != nil {
      return nil, err
  }

  defer rows.Close()

  scores := []score{}

  for rows.Next() {
      var s score
      if err := rows.Scan(&s.ID,
				&s.Pairing1ID,
				&s.Pairing1Strokes,
				&s.Pairing1ID,
				&s.Pairing2Strokes,
				&s.MatchupID); err != nil {
          return nil, err
      }
      scores = append(scores, s)
  }

  return scores, nil
}
