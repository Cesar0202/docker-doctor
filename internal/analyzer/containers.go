package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
)

type ContainerAnalysis struct {
	Total   int
	Stopped int
}

func AnalyzeContainers(ctx context.Context, client *docker.Client) ContainerAnalysis {
	// TODO: Implementar lógica de análisis con client.Cli.ContainerList
	return ContainerAnalysis{
		Total:   0,
		Stopped: 0,
	}
}
