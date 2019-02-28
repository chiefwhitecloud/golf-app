package service

import (
	"fmt"
	"github.com/chiefwhitecloud/golf-app/api"
	"github.com/chiefwhitecloud/golf-app/database"
	"log"
	"net/http"
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
			scoreSheetResult.Score.Matchups[i].SelfPath = fmt.Sprintf("http://%s%s",
				r.Host, scoreSheetResult.Score.Matchups[i].SelfPath)
		}

		respondWithJSON(w, http.StatusOK, scoreSheetResult)

	}
}

func (a *App) handleGetMatchupScoresheet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var matchupScoreInfo api.MatchupScoreInfo

		matchupScoreInfo, err := database.PopulateMatchupScoresheet(a.DB, 1, 1)
		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, matchupScoreInfo)

	}
}
