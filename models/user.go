package models

import "time"

type User struct {
	Id        string  `gorm:"primaryKey" faker:"uuid_hyphenated"`
	Name      string  `faker:"name"`
	Email     string  `faker:"email"`
	Password  string  `faker:"password"`
	Gender    string  `faker:"oneof: male, female"`
	Age       int     `faker:"oneof: 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30"`
	Latitude  float64 `faker:"lat"`
	Longitude float64 `faker:"long"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user User) Empty() bool {
	if user.Id == "" {
		return true
	}

	return false
}
