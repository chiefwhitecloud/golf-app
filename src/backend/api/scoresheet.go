package api

type ScoreHoleInfoSaveRequest struct {
	Scores []HoleScoreInfo `json:scores`
}

type PairingScoreInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CaptainID     string `json:"captainId"`
	TotalHolesWon int    `json:"totalHolesWon"`
}

type HoleScoreInfo struct {
	PairingID string `json:"pairingId"`
	Score     int    `json:"score"`
}

type HoleInfo struct {
	ID               int
	HoleNumber       int             `json:"number"`
	HoleYards        int             `json:"yards"`
	HolePar          int             `json:"par"`
	SelfPath         string          `json:"selfPath"`
	Scores           []HoleScoreInfo `json:"scores,omitempty"`
	WinningPairingID string          `json:"winningPairingID,omitempty"`
}

type MatchupScoreInfo struct {
	Name                 string             `json:"name"`
	Pairings             []PairingScoreInfo `json:"pairings"`
	HoleNumberLastPlayed int                `json:"holeNumberLastPlayed"`
	LeaderPairingID      string             `json:"LeaderPairingId"`
	SelfPath             string             `json:"selfPath"`
	ScoreDetailsPath     string             `json:"scoreDetailsPath"`
	ID                   int								`json:"-"`
}

type ScoreDetailInfo struct {
	HolesInfo []HoleInfo `json:"holes"`
}

type CaptainScores struct {
	TotalHolesWon int `json:"totalHolesWon"`
}

type MatchupScoreInfoResponse struct {
	Matchup        MatchupScoreInfo         `json:"matchup"`
	CaptainsIndent map[string]CaptainIndent `json:"captainIdent"`
}

type ScoreDetailResponse struct {
	ScoreDetail ScoreDetailInfo `json:"scoreDetail"`
}

type ScoreInfo struct {
	Captains            map[string]CaptainScores `json:"captains"`
	Matchups            []MatchupScoreInfo       `json:"matchups"`
	NumOfHolesRemaining int                      `json:"numOfHolesRemaining"`
	TotalNumOfHoles     int                      `json:"totalNumOfHoles"`
}

type CaptainIndent struct {
	Name string `json:"name"`
}

type Scoresheet struct {
	Score          ScoreInfo                `json:"scoresheet"`
	CaptainsIndent map[string]CaptainIndent `json:"captainIdent"`
}
