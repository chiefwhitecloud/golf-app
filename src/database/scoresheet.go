package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/chiefwhitecloud/golf-app/api"
	_ "github.com/lib/pq"
	"strconv"
)

func PopulateScoresheet(db *sql.DB, gameId int) (api.Scoresheet, error) {

	scoresheet := api.Scoresheet{}

	captainsIdent, err := populateCaptainsIdentList(db, gameId)
	if err != nil {
		return scoresheet, err
	}

	matchups, err := getMatchups(db, gameId)
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
	scoresheet.CaptainsIndent = captainsIdent

	return scoresheet, nil

}

//Populates the scoresheet for an individual matchup
//The individual matchup scoresheet contains the hole score and information
func PopulateMatchupScoresheet(db *sql.DB, matchupID int) (api.MatchupScoreInfoResponse, error) {

	matchupScoreInfoResponse := api.MatchupScoreInfoResponse{}

	matchupScoreInfo := api.MatchupScoreInfo{}
	matchup := matchup{ID: matchupID}

	if err := matchup.getMatchup(db); err != nil {
		return matchupScoreInfoResponse, err
	}

	captainsIdent, err := populateCaptainsIdentList(db, matchup.GameID)
	if err != nil {
		return matchupScoreInfoResponse, err
	}

	matchupScoreInfo, err = populateMatchupScoreInfo(db, matchup.Name, matchup.ID)
	if err != nil {
		return matchupScoreInfoResponse, err
	}

	matchupScoreInfoResponse.Matchup = matchupScoreInfo
	matchupScoreInfoResponse.CaptainsIndent = captainsIdent

	return matchupScoreInfoResponse, nil
}

func PopulateMatchupScoreDetail(db *sql.DB, matchupID int) (api.ScoreDetailResponse, error) {

	scoreDetailResponse := api.ScoreDetailResponse{}

	matchup := matchup{ID: matchupID}

	if err := matchup.getMatchup(db); err != nil {
		return scoreDetailResponse, err
	}

	holes, err := getHoles(db, matchup.GameID)
	if err != nil {
		return scoreDetailResponse, err
	}

	holeInfos := make([]api.HoleInfo, len(holes))

	for i := range holes {
		holeInfos[i] = api.HoleInfo{
			ID:         holes[i].ID,
			HoleNumber: holes[i].Number,
			HoleYards:  holes[i].Yards,
			HolePar:    holes[i].Par,
		}
	}

	scoreDetailInfo := api.ScoreDetailInfo{
		HolesInfo: holeInfos,
	}
	scoreDetailResponse.ScoreDetail = scoreDetailInfo

	return scoreDetailResponse, nil

}

func SaveScoreDetail(db *sql.DB, holeId int, matchupId int, scores []api.HoleScoreInfo) error {

	h := hole{ID: holeId}
	err := h.getHole(db)
	if err != nil {
		return err
	}

	matchup := matchup{ID: matchupId}
	if err := matchup.getMatchup(db); err != nil {
		return err
	}

	pairingsForMatchup, err := GetPairingsForMatchup(db, matchup.ID)
	if err != nil {
		return err
	}

	var pairing1Id int
	var pairing2Id int

	for x := range pairingsForMatchup {
		if pairingsForMatchup[x].isFirstPairingInMatchup {
			pairing1Id = pairingsForMatchup[x].ID
		} else {
			pairing2Id = pairingsForMatchup[x].ID
		}
	}

	var pairing1Score int
	var pairing2Score int

	for x := range scores {
		s, err := strconv.Atoi(scores[x].PairingID)
		if err != nil {
			return err
		}
		if s == pairing1Id {
			pairing1Score = scores[x].Score
		} else if s == pairing2Id {
			pairing2Score = scores[x].Score
		}
	}

	if pairing1Score == 0 {
		return errors.New("Pairing1Score not correct")
	}

	if pairing2Score == 0 {
		return errors.New("Pairing2Score not correct")
	}

	s := score{
		HoleID:        h.ID,
		Pairing1ID:    pairing1Id,
		Pairing2ID:    pairing2Id,
		Pairing1Score: pairing1Score,
		Pairing2Score: pairing2Score,
		MatchupID:     matchup.ID,
		HoleNumber:    h.Number,
	}

	err = s.saveScore(db)
	if err != nil {
		return err
	}

	return nil
}

//fill out the MatchupScoreInfo struct for the given matchupID
func populateMatchupScoreInfo(db *sql.DB, matchupName string, matchupID int) (api.MatchupScoreInfo, error) {

	matchupScoreInfo := api.MatchupScoreInfo{}

	pairingsForMatchup, err := GetPairingsForMatchup(db, matchupID)
	if err != nil {
		return matchupScoreInfo, err
	}

	scoresForMatchup, err := getScoresForMatchup(db, matchupID)
	if err != nil {
		return matchupScoreInfo, err
	}

	pairingScoreInfos := make([]api.PairingScoreInfo, len(pairingsForMatchup))

	for x := range pairingsForMatchup {

		totalHolesWon := getTotalHolesWonByPairing(pairingsForMatchup[x].ID, scoresForMatchup)

		pairingScoreInfos[x] = api.PairingScoreInfo{
			ID:            strconv.Itoa(pairingsForMatchup[x].ID),
			Name:          fmt.Sprintf("%s / %s", pairingsForMatchup[x].Player1Name, pairingsForMatchup[x].Player2Name),
			CaptainID:     strconv.Itoa(pairingsForMatchup[x].CaptainID),
			TotalHolesWon: totalHolesWon,
		}
	}

	if pairingScoreInfos[0].TotalHolesWon > pairingScoreInfos[1].TotalHolesWon {
		matchupScoreInfo.LeaderPairingID = pairingScoreInfos[0].ID
	} else if pairingScoreInfos[0].TotalHolesWon < pairingScoreInfos[1].TotalHolesWon {
		matchupScoreInfo.LeaderPairingID = pairingScoreInfos[1].ID
	}

	matchupScoreInfo.Pairings = pairingScoreInfos
	matchupScoreInfo.Name = matchupName
	matchupScoreInfo.HoleNumberLastPlayed = getHoleLastPlayedForMatchup(matchupID, scoresForMatchup)
	matchupScoreInfo.ID = matchupID

	return matchupScoreInfo, nil

}

func populateCaptainsIdentList(db *sql.DB, gameId int) (map[string]api.CaptainIndent, error) {

	captainsIdent := map[string]api.CaptainIndent{}

	captains, err := getCaptains(db, gameId)
	if err != nil {
		return captainsIdent, err
	}

	for i := range captains {
		captainsIdent[strconv.Itoa(captains[i].ID)] = api.CaptainIndent{Name: captains[i].Name}
	}

	return captainsIdent, nil
}
