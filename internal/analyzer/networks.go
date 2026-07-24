package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/network"
)

type NetworkAnalysis struct {
	Total  int
	Unused int
}

func AnalyzeNetworks(ctx context.Context, client *docker.Client) NetworkAnalysis {
	networks, err := client.Cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return NetworkAnalysis{}
	}

	total := len(networks)
	unused := 0

	for _, n := range networks {
		// Ignorar las redes por defecto de docker
		if n.Name == "bridge" || n.Name == "host" || n.Name == "none" {
			continue
		}

		// Si no tiene contenedores conectados
		if len(n.Containers) == 0 {
			unused++
		}
	}

	return NetworkAnalysis{
		Total:  total,
		Unused: unused,
	}
}
