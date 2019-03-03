package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type score struct {
	ID            int
	HoleID        int
	Pairing1ID    int
	Pairing2ID    int
	MatchupID     int
	Pairing1Score int
	Pairing2Score int
	HoleNumber    int
}

func (s *score) saveScore(db *sql.DB) error {

	row := db.QueryRow(`
		SELECT s.id
		FROM score s
  	WHERE
			s.matchup_id = $1
			AND s.hole_id = $2;`,
		s.MatchupID, s.HoleID)

	if err := row.Scan(&s.ID); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if s.ID != 0 {
		//update
		_, err := db.Exec(
			`UPDATE score
			SET pairing1_score = $1, pairing2_score = $2
	    WHERE id = $3`,
			s.Pairing1Score, s.Pairing2Score, s.ID)

		if err != nil {
			return err
		}
	} else {
		//insert
		err := db.QueryRow(
			`INSERT INTO score
	    (hole_id, pairing1_id, pairing1_score, pairing2_id, pairing2_score, matchup_id)
	    VALUES($1, $2, $3, $4, $5, $6)
	    RETURNING id`,
			s.HoleID, s.Pairing1ID, s.Pairing1Score, s.Pairing2ID, s.Pairing2Score, s.MatchupID).Scan(&s.ID)

		if err != nil {
			return err
		}
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
      s.pairing1_score,
			s.pairing2_id,
      s.pairing2_score
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
		if err := rows.Scan(
			&s.ID,
			&s.HoleID,
			&s.HoleNumber,
			&s.Pairing1ID,
			&s.Pairing1Score,
			&s.Pairing1ID,
			&s.Pairing2Score,
		); err != nil {
			return nil, err
		}
		s.MatchupID = matchupId
		scores = append(scores, s)
	}

	return scores, nil
}

func getTotalHolesWonByPairing(pairingId int, scoresForMatchup []score) int {

	holesWon := 0

	for i := 0; i < len(scoresForMatchup); i++ {
		if pairingId == scoresForMatchup[i].Pairing1ID {
			if scoresForMatchup[i].Pairing1Score < scoresForMatchup[i].Pairing2Score {
				holesWon++
			}
		} else if pairingId == scoresForMatchup[i].Pairing2ID {
			if scoresForMatchup[i].Pairing2Score < scoresForMatchup[i].Pairing1Score {
				holesWon++
			}
		}
	}

	return holesWon
}

func getHoleLastPlayedForMatchup(matchupId int, scoresForMatchup []score) int {

	lastHolePlayed := 0

	for i := 0; i < len(scoresForMatchup); i++ {
		if scoresForMatchup[i].MatchupID == matchupId {
			if scoresForMatchup[i].HoleNumber > lastHolePlayed {
				lastHolePlayed = scoresForMatchup[i].HoleNumber
			}
		}
	}

	return lastHolePlayed
}
