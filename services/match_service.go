package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/middleware"
	"github.com/ichtrojan/muzz/models"
	"github.com/ichtrojan/muzz/requests"
	"github.com/ichtrojan/muzz/responses"
	"net/http"
	"time"
)

func Swipe(r *http.Request, payload requests.SwipeRequest) (responses.SwipeResponse, error, int) {
	user := middleware.GetUser(r.Context())

	// prevent self swiping
	if payload.UserId == user.Id {
		return responses.SwipeResponse{}, errors.New("you cannot swipe yourself"), http.StatusNotAcceptable
	}

	// check if the user ID provided belongs to an existing user
	var swipedUser models.User

	database.Connection.Where("id = ?", payload.UserId).First(&swipedUser)

	if swipedUser.Empty() {
		return responses.SwipeResponse{}, errors.New("user does not exist"), http.StatusNotAcceptable
	}

	// check against double swiping
	var swipeExists int64

	database.Connection.Model(&models.Swipe{}).Where("user_id = ?", middleware.GetUser(r.Context()).Id).Where("swiped_on = ?", swipedUser.Id).Count(&swipeExists)

	if swipeExists > 0 {
		return responses.SwipeResponse{}, errors.New("you already swiped on this user"), http.StatusNotAcceptable
	}

	// log `no` and return `false`
	if payload.Preference == "no" {
		err := database.Connection.Create(&models.Swipe{
			Id:         uuid.New().String(),
			UserId:     user.Id,
			SwipedOn:   swipedUser.Id,
			Preference: "no",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}).Error

		if err != nil {
			return responses.SwipeResponse{}, helpers.ServerError(), http.StatusInternalServerError
		}

		return responses.SwipeResponse{
			Matched: false,
			MatchId: "",
		}, nil, http.StatusOK
	}

	// log `yes` and check if there is a mutual swipe before
	// returning `true` otherwise return false
	err := database.Connection.Create(&models.Swipe{
		Id:         uuid.New().String(),
		UserId:     user.Id,
		SwipedOn:   swipedUser.Id,
		Preference: "yes",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}).Error

	if err != nil {
		return responses.SwipeResponse{}, helpers.ServerError(), http.StatusInternalServerError
	}

	var exisingSwipe models.Swipe

	database.Connection.Where("user_id = ?", swipedUser.Id).Where("swiped_on = ?", user.Id).First(&exisingSwipe)

	if exisingSwipe.Empty() {
		return responses.SwipeResponse{
			Matched: false,
			MatchId: "",
		}, nil, http.StatusOK
	} else {
		match := models.SwipeMatch{
			Id:        uuid.New().String(),
			MatchOne:  user.Id,
			MatchTwo:  swipedUser.Id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		database.Connection.Create(&match)

		return responses.SwipeResponse{
			Matched: true,
			MatchId: match.Id,
		}, nil, http.StatusOK
	}
}
