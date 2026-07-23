package util

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

func InvalidRequest(err error) types.APIError {
	errors := make(map[string]string)
	validationErros, ok := err.(validator.ValidationErrors)

	if !ok {
		return types.NewAPIError(http.StatusBadRequest, err)
	}

	for _, verr := range validationErros {
		switch verr.Tag() {
		case "required":
			errors[verr.Field()] = "is required"
		case "email":
			errors[verr.Field()] = "must be an email"
		case "url":
			errors[verr.Field()] = "must be an url"
		case "min":
			errors[verr.Field()] = fmt.Sprintf("must have at least %s characters", verr.Param())
		case "max":
			errors[verr.Field()] = fmt.Sprintf("must have a maximum of %s characters", verr.Param())
		case "gte":
			errors[verr.Field()] = fmt.Sprintf("must be greater or equal than %s", verr.Param())
		case "gt":
			errors[verr.Field()] = fmt.Sprintf("must be greater than %s", verr.Param())
		case "oneof":
			errors[verr.Field()] = fmt.Sprintf("must be one of the following: %s", verr.Param())
		default:
			errors[verr.Field()] = "not valid"
		}
	}

	return types.NewAPIErrorWithDetail(http.StatusUnprocessableEntity, fmt.Errorf("request has invalid data"), errors)
}

func InvalidParamID(paramID string) types.APIError {
	return types.NewAPIError(http.StatusBadRequest, fmt.Errorf("%s must be an integer", paramID))
}

func ResourceNotFound(resource string, ID int) types.APIError {
	return types.NewAPIError(http.StatusNotFound, fmt.Errorf("%s with ID %d does not exist", resource, ID))
}
