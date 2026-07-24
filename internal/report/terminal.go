package report

import (
	"fmt"
	"sort"
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

	// Desglose del Puntaje
	if len(hr.ScoreDetails) > 0 {
		fmt.Println("Desglose:")
		fmt.Println()
		for _, d := range hr.ScoreDetails {
			if d.Points > 0 {
				fmt.Printf("\033[32m+%d\033[0m %s\n", d.Points, d.Reason)
			} else if d.Points < 0 {
				fmt.Printf("\033[31m%d\033[0m %s\n", d.Points, d.Reason)
			}
		}
		fmt.Printf("\nResultado final:\n%d/100\n\n", hr.GlobalScore)
		fmt.Println(strings.Repeat("-", 40))
	}

	// Comparativa con el último escaneo
	if lastScan.ID != 0 {
		fmt.Println("Último análisis:")
		deltaScore := hr.GlobalScore - lastScan.HealthScore
		if deltaScore > 0 {
			fmt.Printf("Health Score: %d ➔ %d (\033[32mMejoró %d pts ↑\033[0m)\n", lastScan.HealthScore, hr.GlobalScore, deltaScore)
		} else if deltaScore < 0 {
			fmt.Printf("Health Score: %d ➔ %d (\033[31mEmpeoró %d pts ↓\033[0m)\n", lastScan.HealthScore, hr.GlobalScore, -deltaScore)
		} else {
			fmt.Printf("Health Score: %d ➔ %d (Sin cambios)\n", lastScan.HealthScore, hr.GlobalScore)
		}

		deltaSpace := hr.TotalRecoverable - lastScan.RecoverableSpaceBytes
		if deltaSpace != 0 {
			deltaMB := deltaSpace / 1024 / 1024
			if deltaMB < 0 {
				fmt.Printf("Espacio recuperable: \033[32mBajó %d MB ↓\033[0m\n", -deltaMB)
			} else if deltaMB > 0 {
				fmt.Printf("Espacio recuperable: \033[31mSubió %d MB ↑\033[0m\n", deltaMB)
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	}

	// Ordenar recomendaciones por Prioridad descendente
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Priority > recs[j].Priority
	})

	// Recomendaciones
	if len(recs) > 0 {
		fmt.Println("\nProblemas encontrados y Recomendaciones:")
		for _, r := range recs {
			fmt.Printf("\n%s\n%s\n", r.Level, r.Message)

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

		// Acciones prioritarias
		fmt.Println("\n==========================")
		fmt.Println("   Acciones Prioritarias")
		fmt.Println("==========================")
		for i, r := range recs {
			if i >= 3 {
				break
			}
			fmt.Printf("\n%d. Ejecutar:\n\033[1m%s\033[0m\n", i+1, r.Command)
			if r.RecoverableSpaceBytes > 0 {
				fmt.Printf("Ganancia estimada:\n%d MB\n", r.RecoverableSpaceBytes/1024/1024)
			}
			if r.Risk != "" {
				fmt.Printf("Riesgo:\n%s\n", r.Risk)
			}
		}

	} else {
		fmt.Println("\n¡Tu entorno está limpio y optimizado! No hay recomendaciones.")
	}
}
