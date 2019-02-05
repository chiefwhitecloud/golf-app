package main_test

import (
	"github.com/chiefwhitecloud/golf-app/service"
	"log"
	"os"
	"testing"
  "net/http"
  "bytes"
  "net/http/httptest"
  "encoding/json"
)

var a service.App

func TestMain(m *testing.M) {
	a = service.App{}
	a.Initialize(
		os.Getenv("POSTGRES_TEST_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test product"}`)

	req, _ := http.NewRequest("POST", "/import", bytes.NewBuffer(payload))
  req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)
    return rr
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS game
(
id SERIAL,
name TEXT NOT NULL,
CONSTRAINT game_pkey PRIMARY KEY (id)
)`
