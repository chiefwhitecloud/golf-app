package database

import (
	"database/sql"
	"github.com/chiefwhitecloud/golf-app/api"
	_ "github.com/lib/pq"
  "strconv"
  "fmt"
)

func PopulateScoresheet(db *sql.DB, gameId int) (api.Scoresheet, error) {

  scoresheet := api.Scoresheet{}

  captains, err := getCaptains(db, gameId);
  if err != nil {
      return scoresheet, err
  }

  captainsList := map[string]api.CaptainIndent{}
	for i := range captains {
		captainsList[strconv.Itoa(captains[i].ID)] = api.CaptainIndent{Name: captains[i].Name}
	}

  scoresheet.CaptainsList = captainsList

  pairings, err := getPairings(db, gameId);
  if err != nil {
    return scoresheet, err
  }

  pairingList := map[string]api.PairingIndent{}
	for i := range pairings {
		pairingList[strconv.Itoa(pairings[i].ID)] = api.PairingIndent{
      Name: fmt.Sprintf("%s / %s", pairings[i].Player1Name, pairings[i].Player2Name),
      CaptainID: strconv.Itoa(pairings[i].CaptainID),
    }
	}

	scoresheet.PairingsList = pairingList

	matchups, err := getMatchups(db, gameId);
  if err != nil {
    return scoresheet, err
  }

	matchupScoreInfos := make([]api.MatchupScoreInfo, len(matchups))

	for i := range matchups {
		matchupScoreInfos[i].Name = matchups[i].Name
	}

	scoreInfo := api.ScoreInfo{}

	scoreInfo.Matchups = matchupScoreInfos

	scoresheet.Score = scoreInfo

  return scoresheet, nil;

}
