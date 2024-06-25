package requests

type SwipeRequest struct {
	UserId     string `json:"userId"`
	Preference string `json:"preference"`
}
