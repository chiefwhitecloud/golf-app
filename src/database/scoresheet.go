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

	matchups, err := getMatchups(db, gameId);
  if err != nil {
    return scoresheet, err
  }

	matchupScoreInfos := make([]api.MatchupScoreInfo, len(matchups))

	for i := range matchups {

		pairingsForMatchup, err := getPairingsForMatchup(db, matchups[i].ID);
	  if err != nil {
	    return scoresheet, err
	  }

		scoresForMatchup, err := getScoresForMatchup(db, matchups[i].ID);
	  if err != nil {
	    return scoresheet, err
	  }

		pairingScoreInfos := make([]api.PairingScoreInfo, len(pairingsForMatchup))

		for x := range pairingsForMatchup {

			holesWon := getTotalHolesWonByPairing(pairingsForMatchup[x].ID, scoresForMatchup)

			pairingScoreInfos[x] = api.PairingScoreInfo{
				 ID: strconv.Itoa(pairingsForMatchup[x].ID),
				 Name: fmt.Sprintf("%s / %s", pairingsForMatchup[x].Player1Name, pairingsForMatchup[x].Player2Name),
				 CaptainID: strconv.Itoa(pairingsForMatchup[x].CaptainID),
				 HolesWon: holesWon,
			}
		}

		if (pairingScoreInfos[0].HolesWon > pairingScoreInfos[1].HolesWon){
			//pairingScoreInfos[0].isWinningMatchup = true
			//pairingScoreInfos[1].isWinningMatchup = false
		} else if (pairingScoreInfos[0].HolesWon < pairingScoreInfos[1].HolesWon) {
			//pairingScoreInfos[0].isWinningMatchup = false
			//pairingScoreInfos[1].isWinningMatchup = true 
		}

		matchupScoreInfos[i].Pairings = pairingScoreInfos
		matchupScoreInfos[i].Name = matchups[i].Name
		matchupScoreInfos[i].HoleNumberLastPlayed = getHoleLastPlayedForMatchup(matchups[i].ID, scoresForMatchup)
	}

	scoreInfo := api.ScoreInfo{}

	scoreInfo.Matchups = matchupScoreInfos

	scoresheet.Score = scoreInfo

  return scoresheet, nil;

}
