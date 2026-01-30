package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ParseBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {

	var data T
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return nil, err
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			ValidationErrorResponse(w, getErrorsMap(ve))
		} else {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
		}

		return nil, err
	}

	return &data, nil
}

func getErrorsMap(ve validator.ValidationErrors) map[string]string {
	errors := map[string]string{}

	for _, err := range ve {
		field := err.Field()
		param := err.Param()
		tag := err.Tag()

		key := strings.ToLower(field[:1]) + field[1:]

		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", key)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", key)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", key, param)
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters long", key, param)
		case "len":
			message = fmt.Sprintf("%s must be exactly %s characters long", key, param)
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", key, param)
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", key, param)
		case "eq":
			message = fmt.Sprintf("%s must be equal to %s", key, param)
		case "ne":
			message = fmt.Sprintf("%s must not be equal to %s", key, param)
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", key)
		case "uuid":
			message = fmt.Sprintf("%s must be a valid UUID", key)
		case "oneof":
			message = fmt.Sprintf("%s must be one of [%s]", key, param)
		case "numeric":
			message = fmt.Sprintf("%s must be numeric", key)
		default:
			message = err.Error()
		}

		errors[key] = message
	}

	return errors

}
