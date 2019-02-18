package api

type PairScoreInfo struct {
	Name string `json:"name"`
	Captain string `json:"captainName"`
	HolesWon int `json:"holesWon"`
}

type MatchupScoreInfo struct {
	Name string `json:"name"`
  Pair1 PairScoreInfo `json:"pair1"`
	Pair2 PairScoreInfo `json:"pair2"`
  LastHolePlayed string `json:"lastHolePlayed"`
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

type PairingIndent struct {
	Name string `json:"name"`
  CaptainID string `json:"captainId"`
}

type Scoresheet struct {
	Score ScoreInfo `json:"scoresheet"`
  PairingsList map[string]PairingIndent `json:"pairingIdent"`
  CaptainsList map[string]CaptainIndent `json:"captainIdent"`
}
