package controllers

import (
	"encoding/json"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/responses"
	"github.com/ichtrojan/muzz/services"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Discover(w http.ResponseWriter, r *http.Request) {
	rules := govalidator.MapData{
		"age":    []string{"numeric"},
		"gender": []string{"alpha", "in:male,female"},
	}

	opts := govalidator.Options{
		Request:         r,
		Rules:           rules,
		RequiredDefault: false,
	}

	validationErrors := helpers.ValidateRequest(opts, "query")

	if len(validationErrors) != 0 {
		helpers.ReturnValidationErrors(w, validationErrors)
		return
	}

	users, pagination, err, status := services.Discover(r)

	if err != nil {
		message := helpers.PrepareMessage(err.Error())

		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(message)

		return
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(responses.DiscoverPaginatedResponse{Users: responses.GenerateDiscoverResponse(r.Context(), users), Meta: pagination})

	return
}
