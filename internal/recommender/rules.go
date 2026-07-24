package recommender

import (
	"docker-doctor/internal/report"
)

type Recommendation struct {
	Level   string // "INFO", "WARNING", "CRITICAL"
	Message string
	Command string
}

func GenerateRecommendations(data report.ReportData) []Recommendation {
	var recs []Recommendation

	// Reglas de contenedores
	if data.Containers.Stopped > 0 {
		recs = append(recs, Recommendation{
			Level:   "WARNING",
			Message: "Tienes contenedores detenidos consumiendo espacio en disco. Considera eliminarlos si ya no los necesitas.",
			Command: "docker container prune",
		})
	}

	// Reglas de imágenes
	if data.Images.Dangling > 0 {
		recs = append(recs, Recommendation{
			Level:   "CRITICAL",
			Message: "Tienes imágenes 'dangling' (huérfanas/basura). Esto suele pasar al hacer builds. Bórralas para liberar disco rápidamente.",
			Command: "docker image prune",
		})
	}

	// Reglas de volúmenes
	if data.Volumes.Orphaned > 0 {
		recs = append(recs, Recommendation{
			Level:   "CRITICAL",
			Message: "Tienes volúmenes de datos que no están conectados a ningún contenedor. Si no respaldaste nada ahí, bórralos.",
			Command: "docker volume prune",
		})
	}

	// Reglas de redes
	if data.Networks.Unused > 0 {
		recs = append(recs, Recommendation{
			Level:   "INFO",
			Message: "Tienes redes de Docker que no están en uso.",
			Command: "docker network prune",
		})
	}

	// Reglas de seguridad (Trivy)
	if data.Security.TotalVulnerabilities > 0 {
		recs = append(recs, Recommendation{
			Level:   "CRITICAL",
			Message: "Trivy detectó vulnerabilidades HIGH/CRITICAL en las imágenes de tus contenedores. Revisa el reporte para más detalles y actualiza tus imágenes.",
			Command: "docker pull <imagen> (luego reconstruye los contenedores)",
		})
	}

	// Reglas de Docker Compose
	if data.Compose.FileFound {
		if len(data.Compose.MissingTags) > 0 {
			recs = append(recs, Recommendation{
				Level:   "WARNING",
				Message: "En tu docker-compose usas imágenes sin tag específico (o con :latest). Esto causa inconsistencias en producción. Fija una versión (ej: node:18.16.0).",
				Command: "Edita docker-compose.yml",
			})
		}
		if len(data.Compose.ExposedPorts) > 0 {
			recs = append(recs, Recommendation{
				Level:   "WARNING",
				Message: "Tienes puertos expuestos globalmente (0.0.0.0) en tu docker-compose.yml. Considera atarlos a 127.0.0.1 (ej: '127.0.0.1:8080:80') para mayor seguridad.",
				Command: "Edita docker-compose.yml",
			})
		}
		if len(data.Compose.PrivilegedSvcs) > 0 {
			recs = append(recs, Recommendation{
				Level:   "CRITICAL",
				Message: "Tienes servicios corriendo en modo privilegiado (privileged: true). Esto es un riesgo de seguridad enorme. Evítalo si no es estrictamente necesario.",
				Command: "Edita docker-compose.yml",
			})
		}
	}

	return recs
}
