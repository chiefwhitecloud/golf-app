package api

type Hole struct {
	Number int `json:"number"`
	Par    int `json:"par"`
	Yards  int `json:"yards"`
}

type CourseInfo struct {
	Name  string `json:"name"`
	Holes []Hole `json:"holes"`
}

type PairInfo struct {
	Players []string `json:"players"`
	Captain string   `json:"captain"`
}

type MatchupInfo struct {
	Name  string     `json:"name"`
	Pairs []PairInfo `json:"pairs"`
}

type MatchInfo struct {
	Course   CourseInfo    `json:"course"`
	Captains []string      `json:"captains"`
	Matchups []MatchupInfo `json:"matchups"`
}

type DataImport struct {
	Match MatchInfo `json:"match"`
}
