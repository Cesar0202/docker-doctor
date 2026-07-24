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

	return recs
}
