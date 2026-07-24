package report

import (
	"docker-doctor/internal/db"
	"fmt"
	"os"
	"strings"
)

func ExportMarkdown(data ReportData, hr HealthReport, lastScan db.ScanHistory, recs []Recommendation, filename string) error {
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

	builder.WriteString(fmt.Sprintf("\n## Docker Health Score: %d/100 (%s)\n\n", hr.GlobalScore, hr.StatusText))

	if lastScan.ID != 0 {
		deltaScore := hr.GlobalScore - lastScan.HealthScore
		if deltaScore > 0 {
			builder.WriteString(fmt.Sprintf("**Último análisis:** %d -> %d (Mejoró %d pts)\n\n", lastScan.HealthScore, hr.GlobalScore, deltaScore))
		} else if deltaScore < 0 {
			builder.WriteString(fmt.Sprintf("**Último análisis:** %d -> %d (Empeoró %d pts)\n\n", lastScan.HealthScore, hr.GlobalScore, -deltaScore))
		}
	}

	builder.WriteString("### Categorías\n")
	for _, cat := range hr.Categories {
		if cat.Score >= 0 {
			builder.WriteString(fmt.Sprintf("- %s: %s\n", cat.Name, cat.Stars))
		}
	}

	if len(recs) > 0 {
		builder.WriteString("\n## 💡 Recomendaciones de la IA\n\n")
		for _, r := range recs {
			builder.WriteString(fmt.Sprintf("### [%s] %s\n\n", r.Level, r.Message))
			if r.Why != "" {
				builder.WriteString(fmt.Sprintf("**¿Por qué importa?**\n%s\n\n", r.Why))
			}
			if r.Impact != "" {
				builder.WriteString(fmt.Sprintf("- **Impacto:** %s\n", r.Impact))
			}
			if r.Risk != "" {
				builder.WriteString(fmt.Sprintf("- **Riesgo:** %s\n", r.Risk))
			}
			if r.RecoverableSpaceBytes > 0 {
				builder.WriteString(fmt.Sprintf("- **Espacio recuperable:** %d MB\n", r.RecoverableSpaceBytes/1024/1024))
			}
			builder.WriteString(fmt.Sprintf("\n```bash\n%s\n```\n\n---\n", r.Command))
		}
	} else {
		builder.WriteString("\n## 💡 Recomendaciones de la IA\n\n¡Tu entorno está limpio y optimizado! No hay recomendaciones.\n")
	}

	err := os.WriteFile(filename, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo Markdown: %w", err)
	}

	fmt.Printf("[OK] Reporte Markdown exportado exitosamente a %s\n", filename)
	return nil
}
