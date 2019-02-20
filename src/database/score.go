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
	HoleNumber int
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

func getScoresForMatchup(db *sql.DB, matchupId int) ([]score, error) {
	rows, err := db.Query(`
		SELECT
      s.id score_id,
			s.hole_id,
			h.number,
			s.pairing1_id,
      s.pairing1_strokes,
			s.pairing2_id,
      s.pairing2_strokes
		FROM score s
			inner join hole h
				on s.hole_id = h.id
    WHERE
    	s.matchup_id = $1;`,
    matchupId)

  if err != nil {
      return nil, err
  }

  defer rows.Close()

  scores := []score{}

  for rows.Next() {
      var s score
      if err := rows.Scan(&s.ID,
				&s.HoleID,
				&s.HoleNumber,
				&s.Pairing1ID,
				&s.Pairing1Strokes,
				&s.Pairing1ID,
				&s.Pairing2Strokes,
			); err != nil {
          return nil, err
      }
      scores = append(scores, s)
  }

  return scores, nil
}


func getTotalHolesWonByPairing(pairingId int, scoresForMatchup []score) int {

	holesWon := 0

	for i := 0; i < len(scoresForMatchup); i++ {
		if (pairingId == scoresForMatchup[i].Pairing1ID){
			if (scoresForMatchup[i].Pairing1Strokes < scoresForMatchup[i].Pairing2Strokes){
				holesWon++
			}
		} else if (pairingId == scoresForMatchup[i].Pairing2ID) {
			if (scoresForMatchup[i].Pairing2Strokes < scoresForMatchup[i].Pairing1Strokes){
				holesWon++
			}
		}
	}

	return holesWon
}

func getHoleLastPlayedForMatchup(matchupId int, scoresForMatchup []score) int {

	lastHolePlayed := 0

	for i := 0; i < len(scoresForMatchup); i++ {
		if (scoresForMatchup[i].MatchupID == matchupId){
			if (scoresForMatchup[i].HoleNumber > lastHolePlayed){
				lastHolePlayed = scoresForMatchup[i].HoleNumber
			}
		}
	}

	return lastHolePlayed
}
