package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
)

type SystemStatus struct {
	IsReachable bool
}

func AnalyzeSystem(ctx context.Context, client *docker.Client) SystemStatus {
	err := client.Ping(ctx)
	return SystemStatus{
		IsReachable: err == nil,
	}
}
