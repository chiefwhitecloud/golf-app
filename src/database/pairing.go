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
	MatchupID int
	Player1Name string
	Player2Name string
}

func (p *pairing) createPairing(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO pairing(player1_id, player2_id, captain_id, matchup_id) VALUES($1, $2, $3, $4) RETURNING id",
		p.Player1ID, p.Player2ID, p.CaptainID, p.MatchupID).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getPairings(db *sql.DB, gameId int) ([]pairing, error) {
	rows, err := db.Query(`
		SELECT p.id pairing_id,
			pl1.id player1_id,
			pl1.name player1_name,
			pl2.id player2_id,
    	pl2.name player2_name,
    	p.captain_id
		FROM pairing p
 			inner join matchup m
    		on p.matchup_id = m.id
    	inner join player pl1
    		on pl1.id = p.player1_id
    	inner join player pl2
    		on pl2.id = p.player2_id
    	WHERE
    		m.game_id = $1;`,
    gameId)

  if err != nil {
      return nil, err
  }

  defer rows.Close()

  pairings := []pairing{}

  for rows.Next() {
      var p pairing
      if err := rows.Scan(&p.ID, &p.Player1ID, &p.Player1Name, &p.Player2ID, &p.Player2Name, &p.CaptainID); err != nil {
          return nil, err
      }
      pairings = append(pairings, p)
  }

  return pairings, nil
}
