package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Docker Doctor - Diagnóstico inteligente para Docker",
	Long: `Docker Doctor es una herramienta CLI para diagnosticar
problemas de rendimiento, recursos y malas prácticas en tu entorno Docker local.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
