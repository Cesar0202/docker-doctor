package cmd

import (
	"fmt"
	"os"

	"docker-doctor/internal/api"

	"github.com/spf13/cobra"
)

var port int

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Inicia el Dashboard Web interactivo",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Serve(port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al iniciar el servidor web: %v\n", err)
		}
	},
}

func init() {
	webCmd.Flags().IntVarP(&port, "port", "p", 8080, "Puerto para el servidor web")
	rootCmd.AddCommand(webCmd)
}
