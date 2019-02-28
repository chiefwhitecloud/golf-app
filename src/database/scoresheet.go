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

	var totalHolesPlayed int

	holesWonByCaptainID := make(map[string]int)

	for i := range matchups {
		matchupScoreInfo, err := populateMatchupScoreInfo(db, matchups[i].Name, matchups[i].ID)

		if err != nil {
			return scoresheet, err
		}

		matchupScoreInfos[i] = matchupScoreInfo
	}



	holeCount, err := getHoleCount(db, gameId)
	if err != nil {
		return scoresheet, err
	}

	captainScores := make(map[string]api.CaptainScores)
	for k, v := range holesWonByCaptainID {
    captainScores[k] = api.CaptainScores{TotalHolesWon: v}
	}

	scoreInfo := api.ScoreInfo{}
	scoreInfo.Captains = captainScores
	scoreInfo.Matchups = matchupScoreInfos
	scoreInfo.TotalNumOfHoles = holeCount
	scoreInfo.NumOfHolesRemaining = holeCount - totalHolesPlayed
	scoresheet.Score = scoreInfo
	scoresheet.CaptainsList = captainsList

  return scoresheet, nil;

}

//Populates the scoresheet for an individual matchup
//The individual matchup scoresheet contains the hole score and information
func PopulateMatchupScoresheet(db *sql.DB, matchupID int, gameId int) (api.MatchupScoreInfo, error) {

	matchupScoreInfo := api.MatchupScoreInfo{}
	matchup := matchup{ID: matchupID}

	if err := matchup.getMatchup(db); err != nil {
		return matchupScoreInfo, err
	}

	matchupScoreInfo, err := populateMatchupScoreInfo(db, matchup.Name, matchup.ID)
	if err != nil {
		return matchupScoreInfo, err
	}

	holes, err := getHoles(db, gameId);
	if err := matchup.getMatchup(db); err != nil {
		return matchupScoreInfo, err
	}

	holeInfos := make([]api.HoleInfo, len(holes))

	for i := range holes {
		holeInfos[i] = api.HoleInfo{
			HoleNumber: holes[i].Number,
			HoleYards: holes[i].Yards,
			HolePar: holes[i].Par,
		}
	}

	matchupScoreInfo.Holes = holeInfos

	return matchupScoreInfo, nil

}

func populateMatchupScoreInfo(db *sql.DB, matchupName string, matchupID int) (api.MatchupScoreInfo, error) {

	matchupScoreInfo := api.MatchupScoreInfo{}

	pairingsForMatchup, err := getPairingsForMatchup(db, matchupID);
	if err != nil {
		return matchupScoreInfo, err
	}

	scoresForMatchup, err := getScoresForMatchup(db, matchupID);
	if err != nil {
		return matchupScoreInfo, err
	}

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
		matchupScoreInfo.LeaderPairingID = pairingScoreInfos[0].ID
	} else if (pairingScoreInfos[0].TotalHolesWon < pairingScoreInfos[1].TotalHolesWon) {
		matchupScoreInfo.LeaderPairingID = pairingScoreInfos[1].ID
	}

	matchupScoreInfo.Pairings = pairingScoreInfos
	matchupScoreInfo.Name = matchupName
	matchupScoreInfo.HoleNumberLastPlayed = getHoleLastPlayedForMatchup(matchupID, scoresForMatchup)
	matchupScoreInfo.SelfPath = fmt.Sprintf("/feeds/default/scoresheet/matchup/%d", matchupID)

	return matchupScoreInfo, nil

}
