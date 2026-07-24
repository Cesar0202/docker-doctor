package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/container"
)

type SecurityStatus struct {
	TrivyInstalled    bool
	ScannedImages     int
	TotalVulnerabilities int
	CriticalCount     int
	HighCount         int
	Details           []string
}

// Estructura simplificada del JSON de Trivy
type TrivyReport struct {
	Results []struct {
		Vulnerabilities []struct {
			VulnerabilityID string `json:"VulnerabilityID"`
			Severity        string `json:"Severity"`
			Title           string `json:"Title"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func AnalyzeSecurity(ctx context.Context, client *docker.Client) SecurityStatus {
	status := SecurityStatus{}

	// Verificar si trivy está instalado
	_, err := exec.LookPath("trivy")
	if err != nil {
		status.TrivyInstalled = false
		return status
	}
	status.TrivyInstalled = true

	// Obtener contenedores en ejecución para escanear sus imágenes
	containers, err := client.Cli.ContainerList(ctx, container.ListOptions{})
	if err != nil || len(containers) == 0 {
		return status
	}

	// Extraer imágenes únicas
	uniqueImages := make(map[string]bool)
	for _, c := range containers {
		uniqueImages[c.Image] = true
	}

	// Escanear hasta 3 imágenes para no hacer el comando lento
	maxScans := 3
	for img := range uniqueImages {
		if status.ScannedImages >= maxScans {
			break
		}

		// Ejecutar Trivy silenciosamente devolviendo JSON
		cmd := exec.CommandContext(ctx, "trivy", "image", "--format", "json", "--quiet", "--severity", "HIGH,CRITICAL", img)
		out, err := cmd.Output()
		if err != nil {
			continue // Ignorar errores de escaneo individual
		}

		var report TrivyReport
		if err := json.Unmarshal(out, &report); err == nil {
			for _, result := range report.Results {
				for _, vuln := range result.Vulnerabilities {
					status.TotalVulnerabilities++
					if vuln.Severity == "CRITICAL" {
						status.CriticalCount++
					} else if vuln.Severity == "HIGH" {
						status.HighCount++
					}
				}
			}
			if status.TotalVulnerabilities > 0 {
				status.Details = append(status.Details, fmt.Sprintf("%s tiene %d vulnerabilidades altas/críticas", img, status.TotalVulnerabilities))
			}
		}
		status.ScannedImages++
	}

	return status
}
