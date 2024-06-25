package helpers

import (
	"encoding/json"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"net/url"
)

func ValidateRequest(options govalidator.Options, method string) url.Values {
	var e url.Values

	v := govalidator.New(options)

	switch method {
	case "json":
		e = v.ValidateJSON()
		break
	case "query":
		e = v.Validate()
	}

	return e
}

func ReturnValidationErrors(w http.ResponseWriter, e url.Values) {
	err := map[string]interface{}{"message": "The given data was invalid", "errors": e}
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(err)
}
