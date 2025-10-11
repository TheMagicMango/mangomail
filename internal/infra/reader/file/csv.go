package file

import (
	"encoding/csv"
	"fmt"
	"os"
)

func (r *FileReader) LoadCSV(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	headers := records[0]

	emailIdx := -1
	for i, header := range headers {
		if header == "email" {
			emailIdx = i
			break
		}
	}

	if emailIdx == -1 {
		return nil, fmt.Errorf("CSV must have an 'email' column")
	}

	rows := make([]map[string]interface{}, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		record := records[i]

		if len(record) <= emailIdx || record[emailIdx] == "" {
			continue
		}

		row := make(map[string]interface{})
		for j, value := range record {
			if j < len(headers) {
				row[headers[j]] = value
			}
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no valid rows found in CSV")
	}

	return rows, nil
}
