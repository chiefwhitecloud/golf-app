package service

import (
	"encoding/json"
	"fmt"
	"github.com/chiefwhitecloud/golf-app/api"
	"github.com/chiefwhitecloud/golf-app/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (a *App) handleGetScoresheet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		scoreSheetResult, err := database.PopulateScoresheet(a.DB, 1)
		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		for i := 0; i < len(scoreSheetResult.Score.Matchups); i++ {
			scoreSheetResult.Score.Matchups[i].SelfPath = fmt.Sprintf("http://%s/feeds/default/scoresheet/matchup/%d",
				r.Host, scoreSheetResult.Score.Matchups[i].ID)
			scoreSheetResult.Score.Matchups[i].ScoreDetailsPath = fmt.Sprintf("http://%s/feeds/default/scoresheet/matchup/%d/scoredetail",
				r.Host, scoreSheetResult.Score.Matchups[i].ID)
		}

		respondWithJSON(w, http.StatusOK, scoreSheetResult)

	}
}

func (a *App) handleGetMatchupScoresheet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var matchupScoreInfoResponse api.MatchupScoreInfoResponse

		matchupScoreInfoResponse, err := database.PopulateMatchupScoresheet(a.DB, 1)
		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		matchupScoreInfoResponse.Matchup.SelfPath = fmt.Sprintf("http://%s/feeds/default/scoresheet/matchup/%d",
			r.Host, matchupScoreInfoResponse.Matchup.ID)
		matchupScoreInfoResponse.Matchup.ScoreDetailsPath = fmt.Sprintf("http://%s/feeds/default/scoresheet/matchup/%d/scoredetail",
			r.Host, matchupScoreInfoResponse.Matchup.ID)

		respondWithJSON(w, http.StatusOK, matchupScoreInfoResponse)

	}
}

func (a *App) handleGetMatchupScoreDetail() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var scoreDetailResponse api.ScoreDetailResponse

		vars := mux.Vars(r)

		matchupID, err := strconv.Atoi(vars["matchup_id"])

		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		scoreDetailResponse, err = database.PopulateMatchupScoreDetail(a.DB, matchupID)
		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		for i := range scoreDetailResponse.ScoreDetail.HolesInfo {
			scoreDetailResponse.ScoreDetail.HolesInfo[i].SelfPath = fmt.Sprintf("http://%s/feeds/default/scoresheet/matchup/%d/scoredetail/%d",
				r.Host, matchupID, scoreDetailResponse.ScoreDetail.HolesInfo[i].ID)
		}

		respondWithJSON(w, http.StatusOK, scoreDetailResponse)

	}
}

func (a *App) handleSaveScoreDetail() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		contentTypeHeader := r.Header.Get("Content-Type")
		if contentTypeHeader != "application/json" {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)

		matchupId, err := strconv.Atoi(vars["matchup_id"])

		if err != nil {
			http.Error(w, "matchup id not valid", 400)
			log.Println(err.Error())
			return
		}

		holeId, err := strconv.Atoi(vars["hole_id"])

		if err != nil {
			http.Error(w, "hole id not valid", 400)
			log.Println(err.Error())
			return
		}

		var scoreSaveRequest api.ScoreHoleInfoSaveRequest

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&scoreSaveRequest)

		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		var scores []api.HoleScoreInfo
		scores = scoreSaveRequest.Scores

		if len(scores) != 2 {
			http.Error(w, "Scores must be in pairs", 400)
		}

		err = database.SaveScoreDetail(a.DB, holeId, matchupId, scores)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to save score", http.StatusInternalServerError)
			return
		}

		respondWithStatus(w, 200)

	}

}
