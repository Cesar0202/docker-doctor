package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/db"
	"docker-doctor/internal/docker"
	"docker-doctor/internal/health"
	"docker-doctor/internal/recommender"
	"docker-doctor/internal/report"
	"docker-doctor/web"
)

// Extendemos ReportData para la API para incluir recomendaciones
type APIResponse struct {
	report.ReportData
	Recommendations []report.Recommendation
	Health          report.HealthReport
	LastScan        db.ScanHistory
}

func Serve(port int) error {
	// 1. Configurar rutas de la API
	http.HandleFunc("/api/report", handleReport)
	http.HandleFunc("/api/history", handleHistory)

	// 2. Extraer y servir la carpeta dist/
	distFs, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		return fmt.Errorf("no se encontró el build del frontend (carpeta dist): %w", err)
	}

	fileServer := http.FileServer(http.FS(distFs))
	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("[INFO] Servidor Web iniciado. Abre tu navegador en http://localhost%s\n", addr)

	return http.ListenAndServe(addr, nil)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()
	client, err := docker.NewClient()
	if err != nil {
		http.Error(w, `{"error": "No se pudo conectar a Docker"}`, http.StatusInternalServerError)
		return
	}

	data := report.ReportData{
		System:     analyzer.AnalyzeSystem(ctx, client),
		Containers: analyzer.AnalyzeContainers(ctx, client),
		Images:     analyzer.AnalyzeImages(ctx, client),
		Volumes:    analyzer.AnalyzeVolumes(ctx, client),
		Networks:   analyzer.AnalyzeNetworks(ctx, client),
		Ports:      analyzer.AnalyzePorts(ctx, client),
		Security:   analyzer.AnalyzeSecurity(ctx, client),
		Compose:    analyzer.AnalyzeCompose(),
	}

	recs := recommender.GenerateRecommendations(data)
	hr := health.CalculateHealth(data)
	lastScan, _ := db.GetLastScan()

	response := APIResponse{
		ReportData:      data,
		Recommendations: recs,
		Health:          hr,
		LastScan:        lastScan,
	}

	json.NewEncoder(w).Encode(response)
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	scans, err := db.GetLatestScans(10)
	if err != nil {
		http.Error(w, `{"error": "No se pudo obtener el historial"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(scans)
}
