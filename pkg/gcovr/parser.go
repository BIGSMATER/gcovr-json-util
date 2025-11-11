package gcovr

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseReport reads and parses a gcovr JSON report file
func ParseReport(filePath string) (*GcovrReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var report GcovrReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", filePath, err)
	}

	return &report, nil
}
