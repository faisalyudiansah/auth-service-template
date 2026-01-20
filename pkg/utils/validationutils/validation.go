package validationutils

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	uppercaseRegexp   = regexp.MustCompile(`[A-Z]`)
	numberRegexp      = regexp.MustCompile(`[0-9]`)
	specialCharRegexp = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':",\.<>\/\\\|` + "`" + `~]`)
	phoneRegexp       = regexp.MustCompile(`^(\+62|62|08)[0-9]{8,12}$`)
	dayOfWeeks        = map[string]bool{
		"Sunday":    true,
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
	}
)

func BirthDateValidator(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(custom_typePkg.DateOnly)
	if !ok {
		return false
	}

	birthDate := time.Time(date)

	if birthDate.IsZero() {
		return false
	}

	now := time.Now()
	if birthDate.After(now) {
		return false
	}

	return true
}

func URLValidator(fl validator.FieldLevel) bool {
	urlStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	if parsedURL.Host == "" {
		return false
	}

	return true
}

func DecimalGT(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.GreaterThan(baseValue)
}

func DecimalLT(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.LessThan(baseValue)
}
func DecimalGTE(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.GreaterThanOrEqual(baseValue)
}

func DecimalLTE(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.LessThanOrEqual(baseValue)
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if ok {
		if strings.Contains(password, " ") {
			return false
		}

		if !uppercaseRegexp.MatchString(password) {
			return false
		}
		if !numberRegexp.MatchString(password) {
			return false
		}
		if !specialCharRegexp.MatchString(password) {
			return false
		}

		if len(password) < 8 || len(password) > 255 {
			return false
		}
		return true
	}

	return false
}

func CleanInputValidator(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if ok {
		if specialCharRegexp.MatchString(input) {
			return false
		}
		if len(input) < 4 || len(input) > 255 {
			return false
		}
		return true
	}
	return false
}

func PhoneNumberValidator(fl validator.FieldLevel) bool {
	phoneNumber, ok := fl.Field().Interface().(string)
	if ok {
		phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
		phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")

		if phoneRegexp.MatchString(phoneNumber) {
			return true
		}
	}

	return false
}

func TimeFormatValidator(fl validator.FieldLevel) bool {
	format := fl.Param()
	_, err := time.Parse(format, fl.Field().String())
	return err == nil
}

func DayOfWeekValidator(fl validator.FieldLevel) bool {
	day := fl.Field().String()
	return dayOfWeeks[day]
}

func NoDuplicatesValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	seen := make(map[interface{}]bool)

	for i := 0; i < field.Len(); i++ {
		value := field.Index(i).Interface()
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}
