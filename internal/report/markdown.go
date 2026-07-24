package report

import (
	"fmt"
	"os"
	"strings"
)

func ExportMarkdown(data ReportData, filename string) error {
	var builder strings.Builder

	builder.WriteString("# Docker Doctor Report\n\n")
	
	if data.System.IsReachable {
		builder.WriteString("## [OK] Sistema Docker\nEstado: Operativo y Accesible\n\n")
	} else {
		builder.WriteString("## [FAIL] Sistema Docker\nEstado: Inaccesible (¿Está corriendo el daemon?)\n\n")
	}

	builder.WriteString("## Contenedores\n")
	builder.WriteString(fmt.Sprintf("- Total: %d\n", data.Containers.Total))
	builder.WriteString(fmt.Sprintf("- Detenidos: %d\n\n", data.Containers.Stopped))

	builder.WriteString("## Imágenes\n")
	builder.WriteString(fmt.Sprintf("- Total: %d\n", data.Images.Total))
	builder.WriteString(fmt.Sprintf("- Dangling (Huérfanas): %d\n\n", data.Images.Dangling))

	builder.WriteString("## Volúmenes\n")
	builder.WriteString(fmt.Sprintf("- Total: %d\n", data.Volumes.Total))
	builder.WriteString(fmt.Sprintf("- Huérfanos: %d\n\n", data.Volumes.Orphaned))

	builder.WriteString("## Redes\n")
	builder.WriteString(fmt.Sprintf("- Total: %d\n", data.Networks.Total))
	builder.WriteString(fmt.Sprintf("- Sin uso: %d\n\n", data.Networks.Unused))

	builder.WriteString("## Puertos\n")
	builder.WriteString(fmt.Sprintf("- Total Expuestos: %d\n", data.Ports.TotalExposed))

	err := os.WriteFile(filename, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo Markdown: %w", err)
	}

	fmt.Printf("✅ Reporte Markdown exportado exitosamente a %s\n", filename)
	return nil
}
