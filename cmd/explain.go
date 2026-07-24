package cmd

import (
	"fmt"
	"strings"

	"docker-doctor/internal/explain"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [mensaje_de_error]",
	Short: "Interpreta y da solución a mensajes de error comunes de Docker",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errorMsg := strings.Join(args, " ")

		fmt.Println("========================================")
		fmt.Println("         DOCKER DOCTOR EXPLAIN          ")
		fmt.Println("========================================")
		fmt.Printf("Analizando error:\n\"%s\"\n\n", errorMsg)

		result := explain.AnalyzeError(errorMsg)

		if result == nil {
			fmt.Println("No tengo este error específico en mi base de datos de conocimiento actual.")
			fmt.Println("Intenta buscar en Google o ChatGPT el mensaje exacto.")
			return
		}

		fmt.Println("Diagnóstico:")
		fmt.Println(result.Diagnosis)
		fmt.Println()

		fmt.Println("Posibles causas:")
		for _, causa := range result.Causes {
			fmt.Printf("✔ %s\n", causa)
		}
		fmt.Println()

		fmt.Println("Comandos útiles:")
		for _, comando := range result.Commands {
			fmt.Printf("➜ \033[1m%s\033[0m\n", comando)
		}
		fmt.Println()

		fmt.Printf("Nivel de confianza:\n%s\n", result.Confidence)
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
