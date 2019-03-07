package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/chiefwhitecloud/golf-app/api"
	_ "github.com/lib/pq"
)

func CreateNewGame(db *sql.DB, importData api.DataImport) error {

	game := game{}
	err := game.createGame(db, "Name")
	if err != nil {
		return err
	}

	captain1 := captain{
		Name:   importData.Match.Captains[0],
		GameID: game.ID,
	}
	err = captain1.createCaptain(db)
	if err != nil {
		return err
	}

	captain2 := captain{
		Name:   importData.Match.Captains[1],
		GameID: game.ID,
	}
	err = captain2.createCaptain(db)
	if err != nil {
		return err
	}

	course := course{
		Name:   importData.Match.Course.Name,
		GameID: game.ID,
	}
	err = course.createCourse(db)
	if err != nil {
		return err
	}

	for i := 0; i < len(importData.Match.Course.Holes); i++ {
		hole := hole{
			Number:   importData.Match.Course.Holes[i].Number,
			Yards:    importData.Match.Course.Holes[i].Yards,
			Par:      importData.Match.Course.Holes[i].Par,
			CourseID: course.ID,
		}

		err = hole.createHole(db)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(importData.Match.Matchups); i++ {

		matchupImport := importData.Match.Matchups[i]

		matchupDB := matchup{
			Name:   fmt.Sprintf("Group %d", i+1),
			GameID: game.ID,
		}

		matchupDB.createMatchup(db)

		for x := 0; x < 2; x++ {
			pair := matchupImport.Pairs[x]

			capt, err := getCaptainByName(captain1, captain2, pair.Captain)
			if err != nil {
				return err
			}
			pa := pairing{
				MatchupID: matchupDB.ID,
				CaptainID: capt.ID,
			}
			for y := 0; y < 2; y++ {
				p := player{
					Name: pair.Players[y],
				}
				err := p.createPlayer(db)
				if err != nil {
					return err
				}
				if y == 0 {
					pa.Player1ID = p.ID
				} else if y == 1 {
					pa.Player2ID = p.ID
				}
			}
			err = pa.createPairing(db)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getCaptainByName(c1 captain, c2 captain, name string) (captain, error) {

	capt := captain{}

	if name == c1.Name {
		return c1, nil
	}

	if name == c2.Name {
		return c2, nil
	}

	return capt, errors.New("No captain found")

}
