package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"docker-doctor/internal/db"
	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)

var showTrend bool

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Muestra el historial de análisis y la tendencia del Health Score",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.InitDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al conectar con la base de datos: %v\n", err)
			return
		}

		scans, err := db.GetLatestScans(15) // Fetch up to 15
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al obtener historial: %v\n", err)
			return
		}

		if len(scans) == 0 {
			fmt.Println("No hay historial de escaneos. Ejecuta 'docker-doctor scan' primero.")
			return
		}

		// Reverse slice so chronological order is left-to-right for the graph
		for i, j := 0, len(scans)-1; i < j; i, j = i+1, j-1 {
			scans[i], scans[j] = scans[j], scans[i]
		}

		if showTrend {
			fmt.Println("========================================")
			fmt.Println("    Tendencia del Health Score (15)")
			fmt.Println("========================================")
			fmt.Println()

			data := make([]float64, len(scans))
			for i, s := range scans {
				data[i] = float64(s.HealthScore)
			}

			graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50), asciigraph.Caption("Health Score en el tiempo"))
			fmt.Println(graph)
			fmt.Println()
			return
		}

		fmt.Println("========================================")
		fmt.Println("       Historial de Docker Doctor       ")
		fmt.Println("========================================")
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "FECHA\tSCORE\tCONTENEDORES\tIMÁGENES\tVOLÚMENES")
		for _, s := range scans {
			dateStr := s.CreatedAt.Format("02/01 15:04")
			
			scoreColor := "\033[36m"
			if s.HealthScore >= 90 {
				scoreColor = "\033[32m"
			} else if s.HealthScore >= 70 {
				scoreColor = "\033[33m"
			} else {
				scoreColor = "\033[31m"
			}

			fmt.Fprintf(w, "%s\t%s%d\033[0m\t%d\t%d\t%d\n", 
				dateStr, 
				scoreColor, s.HealthScore, 
				s.TotalContainers, 
				s.TotalImages, 
				s.TotalVolumes)
		}
		w.Flush()
		fmt.Println()
		fmt.Println("Tip: Usa 'docker-doctor history --trend' para ver la gráfica visual.")
	},
}

func init() {
	historyCmd.Flags().BoolVarP(&showTrend, "trend", "t", false, "Muestra una gráfica ASCII de la tendencia")
	rootCmd.AddCommand(historyCmd)
}
