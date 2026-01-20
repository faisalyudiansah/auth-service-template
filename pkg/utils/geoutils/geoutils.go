package geoutils

import (
	"fmt"
)

func GeoFromText(longitude, latitude string) string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", longitude, latitude)
}
