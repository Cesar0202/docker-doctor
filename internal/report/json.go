package report

import (
	"encoding/json"
	"fmt"
	"os"
	"docker-doctor/internal/db"
)

func ExportJSON(data ReportData, hr HealthReport, lastScan db.ScanHistory, recs []Recommendation, filename string) error {
	exportPayload := struct {
		ReportData
		Health          HealthReport
		LastScan        db.ScanHistory
		Recommendations []Recommendation
	}{
		ReportData:      data,
		Health:          hr,
		LastScan:        lastScan,
		Recommendations: recs,
	}

	bytes, err := json.MarshalIndent(exportPayload, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando JSON: %w", err)
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo JSON: %w", err)
	}

	fmt.Printf("[OK] Reporte JSON exportado exitosamente a %s\n", filename)
	return nil
}
