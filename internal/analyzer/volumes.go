package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/volume"
)

type VolumeAnalysis struct {
	Total    int
	Orphaned int
}

func AnalyzeVolumes(ctx context.Context, client *docker.Client) VolumeAnalysis {
	// Obtenemos todos los volúmenes
	volumes, err := client.Cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return VolumeAnalysis{}
	}

	total := len(volumes.Volumes)
	orphaned := 0

	// Un volumen se considera huérfano si no está siendo usado por ningún contenedor (UsageData no presente o RefCount == 0)
	// Pero VolumeList básico no siempre trae UsageData a menos que se invoque DiskUsage.
	// Haremos una aproximación con UsageData si está disponible.
	for _, v := range volumes.Volumes {
		if v.UsageData == nil || v.UsageData.RefCount == 0 {
			orphaned++
		}
	}

	return VolumeAnalysis{
		Total:    total,
		Orphaned: orphaned,
	}
}
