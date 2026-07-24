package report

import (
	"docker-doctor/internal/analyzer"
)

type ReportData struct {
	System     analyzer.SystemStatus
	Containers analyzer.ContainerAnalysis
	Images     analyzer.ImageAnalysis
	Volumes    analyzer.VolumeAnalysis
	Networks   analyzer.NetworkAnalysis
	Ports      analyzer.PortAnalysis
}
