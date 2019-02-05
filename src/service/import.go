package service

import (
	"encoding/json"
	"github.com/chiefwhitecloud/golf-app/api"
	"github.com/chiefwhitecloud/golf-app/database"
	"log"
	"net/http"
)

func (a *App) handleCreateNewGame() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ct := r.Header.Get("Content-Type")

		if ct != "application/json" {
			log.Println("Header not set")
			http.Error(w, "Please check HTTP Headers", http.StatusBadRequest)
			return
		}

		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}

		var dataimport api.DataImport
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&dataimport)

		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		err = database.CreateNewGame(a.DB, dataimport)

		respondWithStatus(w, http.StatusCreated)

		return

	}
}
