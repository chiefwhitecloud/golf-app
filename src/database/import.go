package database

import (
	"database/sql"
	"github.com/chiefwhitecloud/golf-app/api"
	_ "github.com/lib/pq"
)

func CreateNewGame(db *sql.DB, importData api.DataImport) error {

	game := game{}
	err := game.createGame(db, "Name")
	if err != nil {
		return err
	}

	return nil
}
