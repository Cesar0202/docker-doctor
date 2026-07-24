package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/docker"
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

		data := report.ReportData{
			System:     sysStatus,
			Containers: contStatus,
			Images:     imgStatus,
			Volumes:    volStatus,
			Networks:   netStatus,
			Ports:      portStatus,
		}

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
			report.PrintTerminalReport(data)
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "terminal", "Formato de salida (terminal, json, html, md)")
	scanCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Archivo de salida (por defecto report.[formato])")
	rootCmd.AddCommand(scanCmd)
}
