package csvutils

import (
	"encoding/csv"
	"os"
)

func OpenCsvWriter(name string) (*os.File, *csv.Writer, error) {
	f, err := os.Create(name)
	return f, csv.NewWriter(f), err
}
