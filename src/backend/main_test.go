package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/chiefwhitecloud/golf-app/service"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"strings"
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

	code := m.Run()

	os.Exit(code)
}

func TestSimpleImportTable(t *testing.T) {
	createTables()

	payload := []byte(newGameJSON)

	req, _ := http.NewRequest("POST", "/import", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
}

func TestSimpleScoresheet(t *testing.T) {
	createTables()

	payload := []byte(newGameJSON)

	req, _ := http.NewRequest("POST", "/import", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/feeds/default/scoresheet", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	json := response.Body.String()
	captainIdent := gjson.Get(json, "captainIdent")

	if !captainIdent.Exists() {
		t.Errorf("Expected captainIdent key to exist")
	}

	if !captainIdent.Map()["1"].Exists() {
		t.Errorf("Expected captain 1 key to exist")
	}

	if !captainIdent.Map()["2"].Exists() {
		t.Errorf("Expected captain 2 key to exist")
	}

	captain1 := gjson.Get(captainIdent.Map()["1"].String(), "name")
	captain2 := gjson.Get(captainIdent.Map()["2"].String(), "name")

	if captain1.String() != "madden" && captain2.String() != "madden" {
		t.Errorf("Expected captain madden to exist")
	}

	if captain1.String() != "chief" && captain2.String() != "chief" {
		t.Errorf("Expected captain chief to exist")
	}

	scoresheet := gjson.Get(json, "scoresheet")

	if !scoresheet.Exists() {
		t.Errorf("Expected scoresheet key to exist")
	}

	totalNumOfHoles := gjson.Get(json, "scoresheet.totalNumOfHoles")
	if !totalNumOfHoles.Exists() {
		t.Errorf("Expected scoresheet.totalNumOfHoles key to exist")
	}
	if gjson.Get(json, "scoresheet.totalNumOfHoles").Int() != 3 {
		t.Errorf("Expected scoresheet.totalNumOfHoles to be 3")
	}

	captains := gjson.Get(json, "scoresheet.captains")
	if !captains.Exists() {
		t.Errorf("Expected scoresheet.captains key to exist")
	}
	if gjson.Get(json, "scoresheet.captains.1.totalHolesWon").Int() != 0 {
		t.Errorf("Expected scoresheet.captains.1.totalHolesWon to be 0")
	}
	if gjson.Get(json, "scoresheet.captains.2.totalHolesWon").Int() != 0 {
		t.Errorf("Expected scoresheet.captains.2.totalHolesWon to be 0")
	}

	numOfHolesRemaining := gjson.Get(json, "scoresheet.numOfHolesRemaining")
	if !numOfHolesRemaining.Exists() {
		t.Errorf("Expected scoresheet.numOfHolesRemaining key to exist")
	}
	if gjson.Get(json, "scoresheet.numOfHolesRemaining").Int() != 3 {
		t.Errorf("Expected scoresheet.numOfHolesRemaining to be 3")
	}

	matchups := gjson.Get(json, "scoresheet.matchups")

	if !matchups.Exists() {
		t.Errorf("Expected scoresheet.matchups key to exist")
	}

	if gjson.Get(json, "scoresheet.matchups.#").Int() != 1 {
		t.Errorf("Expected scoresheet.matchups length to be 1")
	}

	if gjson.Get(json, "scoresheet.matchups.0.name").String() != "Group 1" {
		t.Errorf("Expected scoresheet.matchups name to be Group 1")
	}

	if gjson.Get(json, "scoresheet.matchups.0.holeNumberLastPlayed").Int() != 0 {
		t.Errorf("Expected scoresheet.matchups.0.holeNumberLastPlayed to be 0")
	}

	if gjson.Get(json, "scoresheet.matchups.0.pairings.#").Int() != 2 {
		t.Errorf("Expected scoresheet.matchups.0.pairings length to be 2")
	}

	pairingExistsValidation := map[string]bool{
		"White / Campbell": false,
		"Drover / Rogers":  false,
	}

	pairingNamesResult := gjson.Get(json, "scoresheet.matchups.0.pairings.#.name")
	for _, name := range pairingNamesResult.Array() {
		pairingExistsValidation[name.String()] = true
	}

	for k, _ := range pairingExistsValidation {
		if pairingExistsValidation[k] == false {
			t.Errorf("Pairing not found")
		}
	}

	if gjson.Get(json, "scoresheet.matchups.0.pairings.0.holesWon").Int() != 0 {
		t.Errorf("Expected scoresheet.matchups.0.pairings.0.holesWon to be 0")
	}

	if gjson.Get(json, "scoresheet.matchups.0.pairings.1.holesWon").Int() != 0 {
		t.Errorf("Expected scoresheet.matchups.0.pairings.1.holesWon to be 0")
	}

	var matchupArray = gjson.Get(json, "scoresheet.matchups").Array()

	if len(matchupArray) != 1 {
		t.Errorf("matchup array is empty")
	}

	if gjson.Get(matchupArray[0].String(), "selfPath").Exists() == false {
		t.Errorf("selfPath should exist")
	}

	//validate the individual matchup response
	matchupSelfPath := gjson.Get(matchupArray[0].String(), "selfPath").String()

	req, _ = http.NewRequest("GET", matchupSelfPath, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	json = response.Body.String()

	if gjson.Get(json, "matchup.name").String() != "Group 1" {
		t.Errorf("Expected name to be Group 1")
	}

	if gjson.Get(json, "captainIdent").Exists() == false {
		t.Errorf("Expected captainIdent to exist")
	}

	if !gjson.Get(json, "matchup.scoreDetailsPath").Exists() {
		t.Errorf("Expected scoreDetailsPath to exist")
	}

	scoreDetailsPath := gjson.Get(json, "matchup.scoreDetailsPath").String()

	req, _ = http.NewRequest("GET", scoreDetailsPath, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	json = response.Body.String()

	if !gjson.Get(json, "scoreDetail.holes").Exists() {
		t.Errorf("Expected scoreDetails.holes to exist")
	}

	var holesArray = gjson.Get(json, "scoreDetail.holes").Array()

	if len(holesArray) != 3 {
		t.Errorf("score details holes property should contain 3 holes")
	}

	hole1DetailSelfPath := gjson.Get(holesArray[0].String(), "selfPath").String()
	t.Log(hole1DetailSelfPath)

	req, _ = http.NewRequest("PUT", hole1DetailSelfPath, strings.NewReader(`
		{
			"scores":
				[{
						"pairingId": "1",
						"score": 4
					},
					{
						"pairingId": "2",
						"score": 5
					}
				]
		}`))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/feeds/default/scoresheet", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)


}

func createTables() {
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

const newGameJSON = `{
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
}`
