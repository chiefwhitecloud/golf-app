package service

import (
	"net/http"
	"github.com/chiefwhitecloud/golf-app/database"
	"log"
)

func (a *App) handleGetScoresheet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		scoreSheetResult, err := database.PopulateScoresheet(a.DB, 1)
		if err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

 		respondWithJSON(w, http.StatusOK, scoreSheetResult)

	}
}
