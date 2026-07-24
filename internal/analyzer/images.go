package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
)

type ImageAnalysis struct {
	Total     int
	Dangling  int
}

func AnalyzeImages(ctx context.Context, client *docker.Client) ImageAnalysis {
	// TODO: Implementar lógica de análisis con client.Cli.ImageList
	return ImageAnalysis{
		Total:    0,
		Dangling: 0,
	}
}
