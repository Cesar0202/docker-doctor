package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/container"
)

type ContainerAnalysis struct {
	Total                 int
	Stopped               int
	RecoverableSpaceBytes int64
}

func AnalyzeContainers(ctx context.Context, client *docker.Client) ContainerAnalysis {
	containers, err := client.Cli.ContainerList(ctx, container.ListOptions{All: true, Size: true})
	if err != nil {
		return ContainerAnalysis{}
	}

	total := len(containers)
	stopped := 0
	var recoverable int64

	for _, c := range containers {
		if c.State != "running" {
			stopped++
			recoverable += c.SizeRw
		}
	}

	return ContainerAnalysis{
		Total:                 total,
		Stopped:               stopped,
		RecoverableSpaceBytes: recoverable,
	}
}
