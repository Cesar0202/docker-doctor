package report

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExportJSON(data ReportData, filename string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando JSON: %w", err)
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo JSON: %w", err)
	}

	fmt.Printf("✅ Reporte JSON exportado exitosamente a %s\n", filename)
	return nil
}
