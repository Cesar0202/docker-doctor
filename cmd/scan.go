package cmd

import (
	"context"
	"fmt"
	"os"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/docker"
	"docker-doctor/internal/report"

	"github.com/spf13/cobra"
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

		report.PrintTerminalReport(sysStatus, contStatus, imgStatus)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
