package api

type PairScoreInfo struct {
	HolesWon int `json:"holesWon"`
}

type GroupScoreInfo struct {
	Name string `json:"name"`
  Pairs map[string]PairScoreInfo `json:"pairs"`
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
  Groups GroupScoreInfo `json:"groups"`
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
