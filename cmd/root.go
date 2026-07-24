package cmd

import (
	"fmt"
	"os"

	"docker-doctor/internal/db"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Docker Doctor - Diagnóstico inteligente para Docker",
	Long: `Docker Doctor es una herramienta CLI para 
diagnosticar, monitorear y optimizar entornos Docker.
Encuentra problemas comunes y sugiere soluciones.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := db.InitDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Advertencia: no se pudo inicializar SQLite: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
