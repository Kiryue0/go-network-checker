package export

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Kiryue0/go-network-checker/internal/model"
)

func TestSaveJSON(t *testing.T) {
	dir := t.TempDir()
	report := model.ScanReport{
		TotalHost: 2,
		AliveHost: 1,
	}

	err := SaveJSON(report, dir)
	if err != nil {
		t.Fatal(err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	// dosyayı oku, içeriği kontrol et
	data, err := os.ReadFile(dir + "/" + files[0].Name())
	if err != nil {
		t.Fatal(err)
	}

	var result model.ScanReport
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	if result.TotalHost != 2 {
		t.Errorf("expected TotalHost 2, got %d", result.TotalHost)
	}
	if result.AliveHost != 1 {
		t.Errorf("expected AliveHost 1, got %d", result.AliveHost)
	}
}
