package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/db"
	"docker-doctor/internal/docker"
	"docker-doctor/internal/health"
	"docker-doctor/internal/recommender"
	"docker-doctor/internal/report"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	outputFile   string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Ejecuta un diagnóstico completo del entorno Docker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Iniciando escaneo de Docker...")
		ctx := context.Background()

		client, err := docker.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al inicializar cliente Docker: %v\n", err)
			return
		}

		sysStatus := analyzer.AnalyzeSystem(ctx, client)
		contStatus := analyzer.AnalyzeContainers(ctx, client)
		imgStatus := analyzer.AnalyzeImages(ctx, client)
		volStatus := analyzer.AnalyzeVolumes(ctx, client)
		netStatus := analyzer.AnalyzeNetworks(ctx, client)
		portStatus := analyzer.AnalyzePorts(ctx, client)
		secStatus := analyzer.AnalyzeSecurity(ctx, client)
		compStatus := analyzer.AnalyzeCompose()

		data := report.ReportData{
			System:     sysStatus,
			Containers: contStatus,
			Images:     imgStatus,
			Volumes:    volStatus,
			Networks:   netStatus,
			Ports:      portStatus,
			Security:   secStatus,
			Compose:    compStatus,
		}

		// Generar Health Score y Recomendaciones
		hr := health.CalculateHealth(data)
		recs := recommender.GenerateRecommendations(data)

		// Obtener último escaneo antes de guardar este
		lastScan, _ := db.GetLastScan()

		// Guardar el historial en la base de datos (ignorar errores para no detener la ejecución)
		_ = db.SaveScan(db.ScanHistory{
			TotalContainers:       contStatus.Total,
			StoppedContainers:     contStatus.Stopped,
			TotalImages:           imgStatus.Total,
			DanglingImages:        imgStatus.Dangling,
			TotalVolumes:          volStatus.Total,
			OrphanedVolumes:       volStatus.Orphaned,
			TotalNetworks:         netStatus.Total,
			UnusedNetworks:        netStatus.Unused,
			HealthScore:           hr.GlobalScore,
			RecoverableSpaceBytes: hr.TotalRecoverable,
		})

		switch strings.ToLower(outputFormat) {
		case "json":
			if outputFile == "" {
				outputFile = "report.json"
			}
			report.ExportJSON(data, outputFile)
		case "html":
			if outputFile == "" {
				outputFile = "report.html"
			}
			report.ExportHTML(data, outputFile)
		case "markdown", "md":
			if outputFile == "" {
				outputFile = "report.md"
			}
			report.ExportMarkdown(data, outputFile)
		default:
			report.PrintTerminalReport(data, hr, lastScan, recs)
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "terminal", "Formato de salida (terminal, json, html, md)")
	scanCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Archivo de salida (por defecto report.[formato])")
	rootCmd.AddCommand(scanCmd)
}
