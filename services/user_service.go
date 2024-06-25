package services

import (
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/models"
	"github.com/ichtrojan/muzz/responses"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func GenerateUser() (responses.FakeUserResponse, error, int) {
	var fakeUser models.User

	if err := faker.FakeData(&fakeUser); err != nil {
		return responses.FakeUserResponse{}, helpers.ServerError(), http.StatusInternalServerError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fakeUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return responses.FakeUserResponse{}, helpers.ServerError(), http.StatusInternalServerError
	}

	err = database.Connection.Create(&models.User{
		Id:        uuid.New().String(),
		Name:      fakeUser.Name,
		Email:     fakeUser.Email,
		Password:  string(hashedPassword),
		Gender:    fakeUser.Gender,
		Age:       fakeUser.Age,
		Longitude: fakeUser.Longitude,
		Latitude:  fakeUser.Latitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error

	if err != nil {
		return responses.FakeUserResponse{}, helpers.ServerError(), http.StatusInternalServerError
	}

	return responses.GenerateFakeUserResponse(fakeUser), nil, http.StatusOK
}
