package api

type PairingScoreInfo struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CaptainID string `json:"captainId"`
	HolesWon int `json:"holesWon"`
}

type MatchupScoreInfo struct {
	Name string `json:"name"`
  Pairings []PairingScoreInfo `json:"pairings"`
  HoleNumberLastPlayed int `json:"holeNumberLastPlayed"`
  SelfPath string `json:"selfPath"`
}

type CaptainScores struct {
  HolesWon string `json:"holesWon"`
}

type OverallScoreInfo struct {
	Captains map[string]CaptainScores `json:"captains"`
}

type ScoreInfo struct {
	Overall OverallScoreInfo `json:"overall"`
  Matchups []MatchupScoreInfo `json:"matchups"`
  holesToBePlayed int `json:"holesToBePlayed"`
  totalNumOfHoles int `json:"totalNumOfHoles"`
}

type CaptainIndent struct {
	Name string `json:"name"`
}

type Scoresheet struct {
	Score ScoreInfo `json:"scoresheet"`
  CaptainsList map[string]CaptainIndent `json:"captainIdent"`
}
