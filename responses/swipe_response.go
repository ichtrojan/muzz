package responses

type SwipeResponse struct {
	Matched bool   `json:"matched"`
	MatchId string `json:"matchId"`
}
