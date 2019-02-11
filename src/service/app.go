package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chiefwhitecloud/golf-app/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	log.Print(connectionString)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	err = a.databaseConnectionTest()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Database connection established")
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/import", a.handleCreateNewGame()).Methods("POST")
	a.Router.HandleFunc("/feeds/default/scoresheet", a.handleGetScoresheet()).Methods("POST")
}

func (a *App) CreateTables() {
	err := database.CreateTables(a.DB)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) databaseConnectionTest() error {
	return retry(5, time.Second, func() error {
		log.Print("Attempting database connection")
		rows, err := a.DB.Query("SELECT WHERE 1=0")
		if err != nil {
			return err
		}
		defer rows.Close()

		//if err = rows.Err(); err != nil {
		//  return err
		//}

		log.Print("Database connected")
		return nil
	})
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

func respondWithStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
