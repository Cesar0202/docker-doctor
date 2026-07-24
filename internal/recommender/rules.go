package recommender

import (
	"docker-doctor/internal/report"
	"fmt"
)

func GenerateRecommendations(data report.ReportData) []report.Recommendation {
	var recs []report.Recommendation

	// Puertos ocupados en el host (Conflictos de red local)
	if len(data.Ports.HostOccupiedPorts) > 0 {
		for port, service := range data.Ports.HostOccupiedPorts {
			recs = append(recs, report.Recommendation{
				Priority: 100,
				Level:    "🔴 Crítico",
				Message:  fmt.Sprintf("❌ Puerto %d ocupado localmente", port),
				Command:  fmt.Sprintf("Detener %s local o cambiar el puerto en docker-compose", service),
				Why:      fmt.Sprintf("Docker Doctor detectó que %s u otro proceso ya está usando el puerto %d. Si levantas un contenedor exponiendo ese puerto, fallará.", service, port),
				Impact:   "Alto",
				Risk:     "Estabilidad",
			})
		}
	}

	if data.Containers.Stopped > 0 {
		recs = append(recs, report.Recommendation{
			Priority:              60,
			Level:                 "🟡 Advertencia",
			Message:               "Tienes contenedores detenidos.",
			Command:               "docker container prune",
			Why:                   "Los contenedores detenidos siguen existiendo en el disco y retienen configuraciones antiguas que dificultan la administración.",
			Impact:                "Bajo",
			Risk:                  "Ninguno",
			RecoverableSpaceBytes: data.Containers.RecoverableSpaceBytes,
		})
	}

	// Reglas de imágenes
	if data.Images.Dangling > 0 {
		recs = append(recs, report.Recommendation{
			Priority:              70,
			Level:                 "🟠 Importante",
			Message:               "Tienes imágenes 'dangling' (huérfanas/basura).",
			Command:               "docker image prune",
			Why:                   "Estas imágenes ocurren al re-construir un Dockerfile. No tienen etiqueta y ocupan enormes cantidades de disco.",
			Impact:                "Alto",
			Risk:                  "Rendimiento",
			RecoverableSpaceBytes: data.Images.RecoverableSpaceBytes,
		})
	}

	// Reglas de volúmenes
	if data.Volumes.Orphaned > 0 {
		recs = append(recs, report.Recommendation{
			Priority:              80,
			Level:                 "🟠 Importante",
			Message:               "Tienes volúmenes de datos que no están conectados a ningún contenedor.",
			Command:               "docker volume prune",
			Why:                   "Si ya borraste el contenedor original y no necesitas los datos (logs, dbs viejas), bórralos para liberar espacio masivamente.",
			Impact:                "Alto",
			Risk:                  "Ninguno",
			RecoverableSpaceBytes: data.Volumes.RecoverableSpaceBytes,
		})
	}

	// Reglas de redes
	if data.Networks.Unused > 0 {
		recs = append(recs, report.Recommendation{
			Priority: 40,
			Level:    "🟢 Información",
			Message:  "Tienes redes de Docker que no están en uso.",
			Command:  "docker network prune",
			Why:      "Las redes innecesarias complican el enrutamiento interno de Docker.",
			Impact:   "Bajo",
			Risk:     "Ninguno",
		})
	}

	// Reglas de seguridad (Trivy)
	if data.Security.TotalVulnerabilities > 0 {
		recs = append(recs, report.Recommendation{
			Priority: 90,
			Level:    "🔴 Crítico",
			Message:  "Trivy detectó vulnerabilidades HIGH/CRITICAL en las imágenes de tus contenedores.",
			Command:  "docker pull <imagen> (luego reconstruye los contenedores)",
			Why:      "Las vulnerabilidades altas permiten a atacantes comprometer el contenedor o incluso el sistema host si logran escapar.",
			Impact:   "Alto",
			Risk:     "Seguridad",
		})
	}

	// Reglas de Docker Compose
	if data.Compose.FileFound {
		if len(data.Compose.MissingTags) > 0 {
			recs = append(recs, report.Recommendation{
				Priority: 50,
				Level:    "🟡 Advertencia",
				Message:  "En tu docker-compose usas imágenes sin tag específico (o con :latest).",
				Command:  "Edita docker-compose.yml",
				Why:      "Usar 'latest' causa que tu entorno cambie sin previo aviso si la imagen base se actualiza, rompiendo tu app.",
				Impact:   "Medio",
				Risk:     "Estabilidad",
			})
		}
		if len(data.Compose.ExposedPorts) > 0 {
			recs = append(recs, report.Recommendation{
				Priority: 85,
				Level:    "🟠 Importante",
				Message:  "Tienes puertos expuestos globalmente (0.0.0.0) en tu docker-compose.yml.",
				Command:  "Edita docker-compose.yml",
				Why:      "Cualquiera en tu red local (o internet si tu PC está expuesta) puede acceder a la base de datos o servicio. Átalo a '127.0.0.1:PUERTO:PUERTO'.",
				Impact:   "Alto",
				Risk:     "Seguridad",
			})
		}
		if len(data.Compose.PrivilegedSvcs) > 0 {
			recs = append(recs, report.Recommendation{
				Priority: 95,
				Level:    "🔴 Crítico",
				Message:  "Tienes servicios corriendo en modo privilegiado (privileged: true).",
				Command:  "Edita docker-compose.yml",
				Why:      "Un contenedor privilegiado es básicamente 'root' en tu máquina host. Si es comprometido, el atacante toma control de toda tu computadora.",
				Impact:   "Alto",
				Risk:     "Seguridad",
			})
		}
	}

	return recs
}
