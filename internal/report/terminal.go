package report

import (
	"fmt"
	"strings"

	"docker-doctor/internal/db"
)

func PrintTerminalReport(data ReportData, hr HealthReport, lastScan db.ScanHistory, recs []Recommendation) {
	fmt.Println("========================================")
	fmt.Println("         DOCKER DOCTOR                  ")
	fmt.Println("========================================")

	// Health Score
	fmt.Printf("\nDocker Health Score\n\n")

	// Determinar color del score
	scoreColor := "\033[36m" // cyan
	if hr.GlobalScore >= 90 {
		scoreColor = "\033[32m" // green
	} else if hr.GlobalScore >= 70 {
		scoreColor = "\033[33m" // yellow
	} else {
		scoreColor = "\033[31m" // red
	}

	fmt.Printf("%s[%s] %d/100\033[0m\n\n", scoreColor, hr.StatusText, hr.GlobalScore)

	// Comparativa con el último escaneo
	if lastScan.ID != 0 {
		fmt.Println("Último análisis:")
		deltaScore := hr.GlobalScore - lastScan.HealthScore
		if deltaScore > 0 {
			fmt.Printf("Health Score: %d ↓ %d (\033[32mMejoró %d pts\033[0m)\n", lastScan.HealthScore, hr.GlobalScore, deltaScore)
		} else if deltaScore < 0 {
			fmt.Printf("Health Score: %d ↓ %d (\033[31mEmpeoró %d pts\033[0m)\n", lastScan.HealthScore, hr.GlobalScore, -deltaScore)
		} else {
			fmt.Printf("Health Score: %d ↓ %d (Sin cambios)\n", lastScan.HealthScore, hr.GlobalScore)
		}
		
		deltaSpace := hr.TotalRecoverable - lastScan.RecoverableSpaceBytes
		if deltaSpace != 0 {
			deltaMB := deltaSpace / 1024 / 1024
			if deltaMB < 0 {
				fmt.Printf("Espacio recuperable: \033[32mBajó %d MB\033[0m\n", -deltaMB)
			} else if deltaMB > 0 {
				fmt.Printf("Espacio recuperable: \033[31mSubió %d MB\033[0m\n", deltaMB)
			}
		}
		fmt.Println()
	}

	// Puntaje por Categorías
	fmt.Println("Puntaje por categorías:")
	for _, cat := range hr.Categories {
		if cat.Score >= 0 {
			fmt.Printf("%-15s %s\n", cat.Name, cat.Stars)
		}
	}
	fmt.Println("========================================")

	// Recomendaciones
	if len(recs) > 0 {
		fmt.Println("\nProblemas encontrados y Recomendaciones:")
		for _, r := range recs {
			color := "\033[36m"
			prefix := "[INFO]"
			if r.Level == "WARNING" {
				color = "\033[33m"
				prefix = "[WARNING]"
			} else if r.Level == "CRITICAL" {
				color = "\033[31m"
				prefix = "[CRITICAL]"
			}

			fmt.Printf("\n%s%s %s\033[0m\n", color, prefix, r.Message)

			if r.Why != "" {
				fmt.Printf("\n¿Por qué importa?\n%s\n", r.Why)
			}

			if r.Impact != "" {
				fmt.Printf("\nImpacto:\n%s\n", r.Impact)
			}

			if r.Risk != "" {
				fmt.Printf("\nRiesgo:\n%s\n", r.Risk)
			}

			if r.RecoverableSpaceBytes > 0 {
				mb := r.RecoverableSpaceBytes / 1024 / 1024
				if mb > 0 {
					fmt.Printf("\nEspacio recuperable:\n\033[32m%d MB\033[0m\n", mb)
				}
			}

			fmt.Printf("\nComando recomendado:\n\033[1m%s\033[0m\n", r.Command)
			fmt.Println(strings.Repeat("-", 40))
		}
	} else {
		fmt.Println("\n¡Tu entorno está limpio y optimizado! No hay recomendaciones.")
	}
}
