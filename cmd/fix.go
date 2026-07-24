package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/docker"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Asistente interactivo para reparar y limpiar el entorno Docker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("========================================")
		fmt.Println("         DOCKER DOCTOR FIX              ")
		fmt.Println("========================================")
		fmt.Println("Buscando problemas reparables en tu entorno...")
		fmt.Println()

		ctx := context.Background()
		client, err := docker.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al inicializar cliente Docker: %v\n", err)
			return
		}

		contStatus := analyzer.AnalyzeContainers(ctx, client)
		imgStatus := analyzer.AnalyzeImages(ctx, client)
		volStatus := analyzer.AnalyzeVolumes(ctx, client)
		netStatus := analyzer.AnalyzeNetworks(ctx, client)

		foundIssues := false
		fmt.Println("Se encontraron:")

		if imgStatus.Dangling > 0 {
			fmt.Printf("✔ %d imágenes dangling (sin uso)\n", imgStatus.Dangling)
			foundIssues = true
		}
		if contStatus.Stopped > 0 {
			fmt.Printf("✔ %d contenedores detenidos\n", contStatus.Stopped)
			foundIssues = true
		}
		if volStatus.Orphaned > 0 {
			fmt.Printf("✔ %d volúmenes huérfanos\n", volStatus.Orphaned)
			foundIssues = true
		}
		if netStatus.Unused > 0 {
			fmt.Printf("✔ %d redes sin uso\n", netStatus.Unused)
			foundIssues = true
		}

		if !foundIssues {
			fmt.Println("\n¡Buenas noticias! No hay nada de basura que limpiar. Tu sistema está sano.")
			return
		}

		fmt.Printf("\n¿Desea eliminarlos y liberar espacio? (Y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			fmt.Println("\nEjecutando reparaciones...")

			if imgStatus.Dangling > 0 {
				fmt.Println("Limpiando imágenes...")
				runSystemCommand("docker", "image", "prune", "-f")
			}
			if contStatus.Stopped > 0 {
				fmt.Println("Limpiando contenedores...")
				runSystemCommand("docker", "container", "prune", "-f")
			}
			if volStatus.Orphaned > 0 {
				fmt.Println("Limpiando volúmenes...")
				runSystemCommand("docker", "volume", "prune", "-f")
			}
			if netStatus.Unused > 0 {
				fmt.Println("Limpiando redes...")
				runSystemCommand("docker", "network", "prune", "-f")
			}

			fmt.Println("\n✅ ¡Reparación completada! Tu entorno está limpio.")
		} else {
			fmt.Println("\nOperación cancelada. No se ha modificado nada.")
		}
	},
}

func runSystemCommand(name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	_ = c.Run()
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
