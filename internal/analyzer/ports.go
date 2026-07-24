package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/container"
)

type PortAnalysis struct {
	TotalExposed int
	InUse        []uint16
}

func AnalyzePorts(ctx context.Context, client *docker.Client) PortAnalysis {
	containers, err := client.Cli.ContainerList(ctx, container.ListOptions{All: false})
	if err != nil {
		return PortAnalysis{}
	}

	inUse := []uint16{}
	portSet := make(map[uint16]bool)

	for _, c := range containers {
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				if !portSet[p.PublicPort] {
					portSet[p.PublicPort] = true
					inUse = append(inUse, p.PublicPort)
				}
			}
		}
	}

	return PortAnalysis{
		TotalExposed: len(inUse),
		InUse:        inUse,
	}
}
