package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"fmt"
	"net"
	"github.com/docker/docker/api/types/container"
)

type PortAnalysis struct {
	TotalExposed      int
	InUse             []uint16
	HostOccupiedPorts map[int]string
}

func AnalyzePorts(ctx context.Context, client *docker.Client) PortAnalysis {
	containers, err := client.Cli.ContainerList(ctx, container.ListOptions{All: false})
	if err != nil {
		return PortAnalysis{HostOccupiedPorts: make(map[int]string)}
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

	hostOccupied := make(map[int]string)
	commonPorts := map[int]string{
		80:    "HTTP (Apache/Nginx)",
		443:   "HTTPS",
		3306:  "MySQL / MariaDB",
		5432:  "PostgreSQL",
		6379:  "Redis",
		27017: "MongoDB",
		8080:  "HTTP Alt",
	}

	for port, service := range commonPorts {
		// Si docker ya lo está usando, lo ignoramos para el conflicto de host
		if portSet[uint16(port)] {
			continue
		}
		
		// Comprobamos si el puerto está ocupado localmente
		address := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			// Puerto ocupado
			hostOccupied[port] = service
		} else {
			listener.Close()
		}
	}

	return PortAnalysis{
		TotalExposed:      len(inUse),
		InUse:             inUse,
		HostOccupiedPorts: hostOccupied,
	}
}
