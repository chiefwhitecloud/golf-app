package service

import (
	"net/http"
)

func (a *App) handleGetScoresheet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		respondWithStatus(w, http.StatusCreated)

	}
}
