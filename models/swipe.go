package models

import "time"

type Swipe struct {
	Id         string
	UserId     string
	SwipedOn   string
	Preference string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (swipe Swipe) Empty() bool {
	if swipe.Id == "" {
		return true
	}

	return false
}
