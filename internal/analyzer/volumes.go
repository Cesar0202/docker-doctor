package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types"
)

type VolumeAnalysis struct {
	Total                 int
	Orphaned              int
	RecoverableSpaceBytes int64
}

func AnalyzeVolumes(ctx context.Context, client *docker.Client) VolumeAnalysis {
	// Usamos DiskUsage para obtener información exacta de uso y tamaño
	du, err := client.Cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return VolumeAnalysis{}
	}

	total := len(du.Volumes)
	orphaned := 0
	var recoverable int64

	for _, v := range du.Volumes {
		if v.UsageData == nil || v.UsageData.RefCount == 0 {
			orphaned++
			if v.UsageData != nil {
				recoverable += v.UsageData.Size
			}
		}
	}

	return VolumeAnalysis{
		Total:                 total,
		Orphaned:              orphaned,
		RecoverableSpaceBytes: recoverable,
	}
}
