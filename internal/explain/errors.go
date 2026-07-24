package explain

import (
	"regexp"
)

type Explanation struct {
	Diagnosis  string
	Causes     []string
	Commands   []string
	Confidence string
}

type ErrorPattern struct {
	Regex    *regexp.Regexp
	Response Explanation
}

var Database = []ErrorPattern{
	{
		Regex: regexp.MustCompile(`(?i)bind: address already in use|bind for .*:(\d+) failed`),
		Response: Explanation{
			Diagnosis: "El puerto que intentas usar ya está siendo utilizado por otro proceso.",
			Causes: []string{
				"Un servicio local (como Apache, Nginx, PostgreSQL) está ocupando el puerto",
				"Ya existe otro contenedor de Docker corriendo y usando ese mismo puerto",
			},
			Commands: []string{
				"docker ps (para ver si otro contenedor lo usa)",
				"netstat -ano (en Windows) o lsof -i :PUERTO (en Linux/Mac) para encontrar el proceso culpable",
			},
			Confidence: "98%",
		},
	},
	{
		Regex: regexp.MustCompile(`(?i)manifest .* not found|manifest unknown`),
		Response: Explanation{
			Diagnosis: "Docker no puede encontrar la imagen solicitada en el registro (Docker Hub).",
			Causes: []string{
				"Escribiste mal el nombre de la imagen",
				"Escribiste mal el tag o versión (ej. usaste :latest y no existe)",
				"La imagen es privada y no has hecho 'docker login'",
			},
			Commands: []string{
				"docker login (si es un registro privado)",
				"Revisa la sintaxis en tu docker-compose.yml o comando run",
			},
			Confidence: "95%",
		},
	},
	{
		Regex: regexp.MustCompile(`(?i)no space left on device`),
		Response: Explanation{
			Diagnosis: "El disco duro (o la partición donde Docker guarda sus datos) está completamente lleno.",
			Causes: []string{
				"Tienes demasiadas imágenes huérfanas (dangling)",
				"Un contenedor generó logs masivos que llenaron el disco",
				"Volúmenes de datos muy pesados sin limpiar",
			},
			Commands: []string{
				"docker system prune -a --volumes (⚠️ Borrará todo lo que no esté corriendo)",
				"docker-doctor fix (Para usar el asistente interactivo de limpieza)",
			},
			Confidence: "99%",
		},
	},
	{
		Regex: regexp.MustCompile(`(?i)Cannot connect to the Docker daemon`),
		Response: Explanation{
			Diagnosis: "La línea de comandos no puede comunicarse con el motor interno de Docker.",
			Causes: []string{
				"Docker Desktop no está abierto",
				"El servicio de Docker no está iniciado en Linux",
				"No tienes permisos (necesitas sudo o pertenecer al grupo docker)",
			},
			Commands: []string{
				"Abre la aplicación Docker Desktop (en Windows/Mac)",
				"sudo systemctl start docker (en Linux)",
			},
			Confidence: "99%",
		},
	},
}

func AnalyzeError(errorMsg string) *Explanation {
	for _, pattern := range Database {
		if pattern.Regex.MatchString(errorMsg) {
			return &pattern.Response
		}
	}
	return nil
}
