package constant

import "time"

const (
	DEBUG   = "debug"
	RELEASE = "release"
)

const (
	ResponseSuccessMessage = "success"
)

const (
	WIB = 7 * time.Hour
)

const (
	DEFAULT_LIMIT = 10
	DEFAULT_PAGE  = 1
)

const (
	DEFAULT_DATE_ONLY        = "2006-01-02"
	DEFAULT_DATE_ONLY_FORMAT = "YYYY-MM-DD"
)

var timeLayoutTranslate map[string]string = map[string]string{
	"02-01-2006":      "DD-MM-YYYY",
	DEFAULT_DATE_ONLY: DEFAULT_DATE_ONLY_FORMAT,
	"2006":            "YYYY",
	"15:04":           "hh:mm",
}

func ConvertGoTimeLayoutToReadable(layout string) string {
	return timeLayoutTranslate[layout]
}
