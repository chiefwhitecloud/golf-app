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

	matchups, err := getMatchups(db, gameId);
  if err != nil {
    return scoresheet, err
  }

	matchupScoreInfos := make([]api.MatchupScoreInfo, len(matchups))

	var totalHolesPlayed int;

	for i := range matchups {

		pairingsForMatchup, err := getPairingsForMatchup(db, matchups[i].ID);
	  if err != nil {
	    return scoresheet, err
	  }

		scoresForMatchup, err := getScoresForMatchup(db, matchups[i].ID);
	  if err != nil {
	    return scoresheet, err
	  }

		totalHolesPlayed += len(scoresForMatchup)

		pairingScoreInfos := make([]api.PairingScoreInfo, len(pairingsForMatchup))

		for x := range pairingsForMatchup {

			totalHolesWon := getTotalHolesWonByPairing(pairingsForMatchup[x].ID, scoresForMatchup)

			pairingScoreInfos[x] = api.PairingScoreInfo{
				 ID: strconv.Itoa(pairingsForMatchup[x].ID),
				 Name: fmt.Sprintf("%s / %s", pairingsForMatchup[x].Player1Name, pairingsForMatchup[x].Player2Name),
				 CaptainID: strconv.Itoa(pairingsForMatchup[x].CaptainID),
				 TotalHolesWon: totalHolesWon,
			}
		}

		if (pairingScoreInfos[0].TotalHolesWon > pairingScoreInfos[1].TotalHolesWon){
			matchupScoreInfos[i].LeaderPairingID = pairingScoreInfos[0].ID
		} else if (pairingScoreInfos[0].TotalHolesWon < pairingScoreInfos[1].TotalHolesWon) {
			matchupScoreInfos[i].LeaderPairingID = pairingScoreInfos[1].ID
		}

		matchupScoreInfos[i].Pairings = pairingScoreInfos
		matchupScoreInfos[i].Name = matchups[i].Name
		matchupScoreInfos[i].HoleNumberLastPlayed = getHoleLastPlayedForMatchup(matchups[i].ID, scoresForMatchup)
	}

	holeCount, err := getHoleCount(db, gameId)
	if err != nil {
		return scoresheet, err
	}

	scoreInfo := api.ScoreInfo{}
	scoreInfo.Matchups = matchupScoreInfos
	scoreInfo.TotalNumOfHoles = holeCount
	scoreInfo.NumOfHolesRemaining = holeCount - totalHolesPlayed

	scoresheet.Score = scoreInfo

	scoresheet.CaptainsList = captainsList

  return scoresheet, nil;

}
