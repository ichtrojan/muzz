package models

import "time"

type SwipeMatch struct {
	Id        string
	MatchOne  string
	MatchTwo  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (swipeMatch SwipeMatch) Empty() bool {
	if swipeMatch.Id == "" {
		return true
	}

	return false
}
