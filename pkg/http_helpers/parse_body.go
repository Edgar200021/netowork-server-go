package http_helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ParseBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var parsed T

	err := json.NewDecoder(r.Body).Decode(&parsed)
	if err != nil {
		WriteErrorJson(w, "Invalid request body", http.StatusBadRequest)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(parsed); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorsMap := make(map[string][]string)
			structType := reflect.TypeOf(parsed)

			for _, ve := range validationErrors {
				field := jsonFieldName(structType, ve.Field())
				errorsMap[field] = append(errorsMap[field], validationMessage(ve))
			}

			WriteValidationErrorsJson(w, errorsMap)
			return nil, err
		}

		WriteErrorJson(w, "Validation failed", http.StatusBadRequest)
		return nil, err
	}

	return &parsed, nil
}

func jsonFieldName(structType reflect.Type, fieldName string) string {
	field, ok := structType.FieldByName(fieldName)
	if !ok {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	
	name := strings.Split(jsonTag, ",")[0]
	return name
}

func validationMessage(ve validator.FieldError) string {
	isString := ve.Kind() == reflect.String

	switch ve.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		if isString {
			return fmt.Sprintf("Minimum length is %s characters", ve.Param())
		}
		return fmt.Sprintf("Minimum value is %s", ve.Param())
	case "max":
		if isString {
			return fmt.Sprintf("Maximum length is %s characters", ve.Param())
		}
		return fmt.Sprintf("Maximum value is %s", ve.Param())
	case "gte":
		if isString {
			return fmt.Sprintf("Must have at least %s characters", ve.Param())
		}
		return fmt.Sprintf("Must be greater than or equal to %s", ve.Param())
	case "lte":
		if isString {
			return fmt.Sprintf("Must have at most %s characters", ve.Param())
		}
		return fmt.Sprintf("Must be less than or equal to %s", ve.Param())
	case "gt":
		if isString {
			return fmt.Sprintf("Must have more than %s characters", ve.Param())
		}
		return fmt.Sprintf("Must be greater than %s", ve.Param())
	case "lt":
		if isString {
			return fmt.Sprintf("Must have less than %s characters", ve.Param())
		}
		return fmt.Sprintf("Must be less than %s", ve.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", strings.ReplaceAll(ve.Param(), " ", ", "))
	default:
		return fmt.Sprintf("Validation failed on '%s'", ve.Tag())
	}
}
