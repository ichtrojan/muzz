package controllers

import (
	"encoding/json"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/requests"
	"github.com/ichtrojan/muzz/responses"
	"github.com/ichtrojan/muzz/services"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

type CreateUserResponse struct {
	User responses.FakeUserResponse `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var request requests.LoginRequest

	rules := govalidator.MapData{
		"email":    []string{"email", "required"},
		"password": []string{"alpha_space", "required"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &request,
		Rules:   rules,
	}

	validationErrors := helpers.ValidateRequest(opts, "json")

	if len(validationErrors) != 0 {
		helpers.ReturnValidationErrors(w, validationErrors)
		return
	}

	response, err, status := services.Login(request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fakeUser, err, status := services.GenerateUser()

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(CreateUserResponse{User: fakeUser})
	return
}
