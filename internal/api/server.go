package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"docker-doctor/internal/analyzer"
	"docker-doctor/internal/docker"
	"docker-doctor/internal/report"
	"docker-doctor/web"
)

func Serve(port int) error {
	// 1. Configurar ruta de la API
	http.HandleFunc("/api/report", handleReport)

	// 2. Extraer y servir la carpeta dist/
	distFs, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		return fmt.Errorf("no se encontró el build del frontend (carpeta dist): %w", err)
	}
	
	fileServer := http.FileServer(http.FS(distFs))
	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 Servidor Web iniciado. Abre tu navegador en http://localhost%s\n", addr)
	
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
	}

	json.NewEncoder(w).Encode(data)
}
