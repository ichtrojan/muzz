package services

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/models"
	"github.com/ichtrojan/muzz/requests"
	"github.com/ichtrojan/muzz/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Login(request requests.LoginRequest) (responses.LoginResponse, error, int) {
	var user models.User

	database.Connection.Where("email = ?", request.Email).First(&user)

	if user.Empty() {
		return responses.LoginResponse{}, errors.New("invalid credentials"), http.StatusUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return responses.LoginResponse{}, errors.New("invalid credentials"), http.StatusUnauthorized
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * (24 * 7)).Unix(),
		"user_id": user.Id,
	})

	token, err := at.SignedString([]byte("testing"))

	if err != nil {
		return responses.LoginResponse{}, helpers.ServerError(), http.StatusInternalServerError
	}

	return responses.LoginResponse{
		Token: token,
		User:  responses.GenerateUserResponse(user),
	}, nil, http.StatusOK
}
