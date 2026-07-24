package analyzer

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ComposeAnalysis struct {
	FileFound      bool
	TotalServices  int
	MissingTags    []string
	ExposedPorts   []string
	PrivilegedSvcs []string
}

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	Image      string   `yaml:"image"`
	Ports      []string `yaml:"ports"`
	Privileged bool     `yaml:"privileged"`
}

func AnalyzeCompose() ComposeAnalysis {
	status := ComposeAnalysis{}

	// Buscar docker-compose.yml o compose.yaml
	filename := "docker-compose.yml"
	data, err := os.ReadFile(filename)
	if err != nil {
		filename = "compose.yaml"
		data, err = os.ReadFile(filename)
		if err != nil {
			return status
		}
	}

	status.FileFound = true

	var compose composeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return status
	}

	status.TotalServices = len(compose.Services)

	for name, svc := range compose.Services {
		// Analizar etiquetas de imagen (latest o vacía)
		if svc.Image != "" {
			parts := strings.Split(svc.Image, ":")
			if len(parts) == 1 || parts[len(parts)-1] == "latest" {
				status.MissingTags = append(status.MissingTags, name)
			}
		}

		// Analizar puertos (verificar si están atados a 0.0.0.0 implícitamente)
		for _, port := range svc.Ports {
			// Un puerto sin IP especificada (ej. "8080:80") se une a 0.0.0.0, lo cual puede ser inseguro
			// Si tiene IP, suele tener 3 partes separadas por ':' (ej "127.0.0.1:8080:80")
			if !strings.Contains(port, "127.0.0.1") && strings.Count(port, ":") == 1 {
				status.ExposedPorts = append(status.ExposedPorts, fmt.Sprintf("%s (%s)", name, port))
			}
		}

		// Analizar privilegios
		if svc.Privileged {
			status.PrivilegedSvcs = append(status.PrivilegedSvcs, name)
		}
	}

	return status
}
