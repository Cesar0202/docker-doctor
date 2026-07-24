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
	Security   analyzer.SecurityStatus
	Compose    analyzer.ComposeAnalysis
}

type Recommendation struct {
	Level                 string
	Message               string
	Command               string
	Why                   string
	Impact                string
	Risk                  string
	RecoverableSpaceBytes int64
}

type HealthCategory struct {
	Name  string
	Stars string
	Score int
}

type HealthReport struct {
	GlobalScore      int
	StatusText       string
	Categories       []HealthCategory
	TotalRecoverable int64
}
