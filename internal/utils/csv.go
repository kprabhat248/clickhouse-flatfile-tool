package utils

import (
	"encoding/csv"
	"os"
)

// ReadCSVHeaders opens the given CSV file and returns the first row (headers).
func ReadCSVHeaders(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return headers, nil
}

// GetColumnIndices maps selected columns to their indices in headers.
func GetColumnIndices(headers []string, selected []string) []int {
	var indices []int
	for _, sel := range selected {
		for idx, header := range headers {
			if header == sel {
				indices = append(indices, idx)
			}
		}
	}
	return indices
}
