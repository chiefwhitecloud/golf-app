package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/chiefwhitecloud/golf-app/service"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	ensureTablesExist()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	payload := []byte(`{
		"match": {
			"captains": ["chief", "madden"],
			"course": {
				"name": "Twin Rivers",
				"holes": [{
						"number": 1,
						"par": 4,
						"yards": 345
					},
					{
						"number": 2,
						"par": 5,
						"yards": 445
					},
					{
						"number": 3,
						"par": 3,
						"yards": 145
					}]
			},
			"matchups": [{
        "name": "group 1",
        "pairs":  [{
            "players": ["White", "Campbell"],
            "captain": "chief"
          },
          {
            "players": ["Drover", "Rogers"],
            "captain": "madden"
          }]
      }]
		}
	}`)

	req, _ := http.NewRequest("POST", "/import", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

}

func ensureTablesExist() {
	a.CreateTables()
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
