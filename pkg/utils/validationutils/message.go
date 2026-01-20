package validationutils

import (
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"

	"github.com/go-playground/validator/v10"
)

func TagToMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "birth_date":
		return "birth date cannot be in the future"
	case "oneof":
		return fmt.Sprintf("must be one of %s", fe.Param())
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "password":
		return fmt.Sprintf("%s must contain at least 8 characters including 1 number, 1 special character, and 1 capital letter excluding whitespaces", fe.Field())
	case "len":
		return fmt.Sprintf("%s length or value must be exactly %v", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s length or value %v must be at most", fe.Field(), fe.Param())
	case "dgte":
		return fmt.Sprintf("%s must be greater than or equal to %v", fe.Field(), fe.Param())
	case "dlte":
		return fmt.Sprintf("%s must be less than or equal to %v", fe.Field(), fe.Param())
	case "dgt":
		return fmt.Sprintf("%s must be greater than to %v", fe.Field(), fe.Param())
	case "dlt":
		return fmt.Sprintf("%s must be less than to %v", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %v", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be lower than or equal to %v", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s has invalid email format", fe.Field())
	case "eq":
		return fmt.Sprintf("%s must be equal to %v", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("%s length or value must be at least %v", fe.Field(), fe.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be an URL", fe.Field())
	case "boolean":
		return fmt.Sprintf("%s must be a boolean", fe.Field())
	case "latitude":
		return fmt.Sprintf("%s must be a latitude", fe.Field())
	case "longitude":
		return fmt.Sprintf("%s must be a longitude", fe.Field())
	case "base64":
		return fmt.Sprintf("%s must be in base64 format", fe.Field())
	case "gtecsfield":
		return fmt.Sprintf("%s must be greater than or equal to %v", fe.Field(), fe.Param())
	case "ltecsfield":
		return fmt.Sprintf("%s must be lower than or equal to %v", fe.Field(), fe.Param())
	case "phone_number":
		return fmt.Sprintf("%s has an invalid phone number format", fe.Field())
	case "time_format":
		return fmt.Sprintf("please send time in format of %s", constant.ConvertGoTimeLayoutToReadable(fe.Param()))
	case "no_duplicates":
		return fmt.Sprintf("%s duplicate values are not allowed", fe.Field())
	case "day_of_weeks":
		return fmt.Sprintf("%s must be a valid day of the week (e.g, Sunday, Monday, etc)", fe.Field())
	case "clean_input":
		return fmt.Sprintf("%s must contain at least 4 characters with excluding whitespaces and special character", fe.Field())
	default:
		return "invalid input"
	}
}
