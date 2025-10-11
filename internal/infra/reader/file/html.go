package file

import (
	"fmt"
	"os"
)

func (r *FileReader) LoadHTML(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML file: %w", err)
	}

	return string(data), nil
}
