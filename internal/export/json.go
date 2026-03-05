package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Kiryue0/go-network-checker/internal/model"
)

func SaveJSON(report model.ScanReport, outputDir string) error {
	scanDate := time.Now().Format("2006-01-02_15-04-05")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to json: %w", err)
	}
	filePath := outputDir + "/scan_" + scanDate + ".json"

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
