package custom_type

import (
	"fmt"
	"strings"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	constantPkg "github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(`"` + t.Format(constantPkg.DEFAULT_DATE_ONLY) + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	t, err := time.Parse(constantPkg.DEFAULT_DATE_ONLY, s)
	if err != nil {
		return &apperror.FieldBindError{
			Field:   "birth_date",
			Message: fmt.Sprintf("invalid date format, use %s", constantPkg.DEFAULT_DATE_ONLY_FORMAT),
		}
	}

	*d = DateOnly(t)
	return nil
}

func (d DateOnly) Format(layout string) string {
	t := time.Time(d)
	return t.Format(layout)
}

func (d DateOnly) String() string {
	t := time.Time(d)
	if t.IsZero() {
		return ""
	}
	return t.Format(constantPkg.DEFAULT_DATE_ONLY)
}

func (d DateOnly) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d DateOnly) IsAfter(t time.Time) bool {
	return time.Time(d).After(t)
}

func (d DateOnly) IsBefore(t time.Time) bool {
	return time.Time(d).Before(t)
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}
