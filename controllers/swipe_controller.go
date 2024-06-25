package controllers

import (
	"encoding/json"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/requests"
	"github.com/ichtrojan/muzz/services"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Swipe(w http.ResponseWriter, r *http.Request) {
	var request requests.SwipeRequest

	rules := govalidator.MapData{
		"userId":     []string{"uuid", "required"},
		"preference": []string{"alpha", "in:yes,no", "required"},
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

	match, err, status := services.Swipe(r, request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(match)
	return
}
