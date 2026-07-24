package health

import (
	"docker-doctor/internal/report"
	"fmt"
)

func CalculateHealth(data report.ReportData) report.HealthReport {
	hr := report.HealthReport{}

	// Calcular estrellas por categoría y detalles
	catEngine, detEngine := calculateEngineScore(data)
	catContainers, detContainers := calculateContainersScore(data)
	catImages, detImages := calculateImagesScore(data)
	catVolumes, detVolumes := calculateVolumesScore(data)
	catSecurity, detSecurity := calculateSecurityScore(data)
	catCompose, detCompose := calculateComposeScore(data)

	hr.Categories = []report.HealthCategory{
		catEngine, catContainers, catImages, catVolumes, catSecurity, catCompose,
	}

	hr.ScoreDetails = append(hr.ScoreDetails, detEngine...)
	hr.ScoreDetails = append(hr.ScoreDetails, detContainers...)
	hr.ScoreDetails = append(hr.ScoreDetails, detImages...)
	hr.ScoreDetails = append(hr.ScoreDetails, detVolumes...)
	hr.ScoreDetails = append(hr.ScoreDetails, detSecurity...)
	hr.ScoreDetails = append(hr.ScoreDetails, detCompose...)

	// Score global es el promedio
	totalScore := 0
	count := 0
	for _, c := range hr.Categories {
		if c.Score >= 0 { // -1 significa N/A
			totalScore += c.Score
			count++
		}
	}

	if count > 0 {
		hr.GlobalScore = totalScore / count
	}

	if hr.GlobalScore >= 90 {
		hr.StatusText = "Excelente"
	} else if hr.GlobalScore >= 70 {
		hr.StatusText = "Bueno"
	} else if hr.GlobalScore >= 50 {
		hr.StatusText = "Atención"
	} else {
		hr.StatusText = "Crítico"
	}

	hr.TotalRecoverable = data.Containers.RecoverableSpaceBytes +
		data.Images.RecoverableSpaceBytes +
		data.Volumes.RecoverableSpaceBytes

	return hr
}

func getStars(score int) string {
	if score < 0 {
		return "N/A"
	}
	if score >= 90 { return "★★★★★" }
	if score >= 75 { return "★★★★☆" }
	if score >= 50 { return "★★★☆☆" }
	if score >= 30 { return "★★☆☆☆" }
	return "★☆☆☆☆"
}

func calculateEngineScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	score := 100
	var details []report.ScoreDetail
	if !data.System.IsReachable {
		score = 0
		details = append(details, report.ScoreDetail{Points: -100, Reason: "Docker Engine no responde"})
	} else {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "Docker Engine funcionando correctamente"})
	}
	return report.HealthCategory{Name: "Docker Engine", Score: score, Stars: getStars(score)}, details
}

func calculateContainersScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	score := 100
	var details []report.ScoreDetail
	if data.Containers.Stopped > 0 {
		penalty := data.Containers.Stopped * 10
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d contenedor(es) detenido(s)", data.Containers.Stopped)})
	} else {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "No hay contenedores detenidos"})
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Contenedores", Score: score, Stars: getStars(score)}, details
}

func calculateImagesScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	score := 100
	var details []report.ScoreDetail
	if data.Images.Dangling > 0 {
		penalty := data.Images.Dangling * 15
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d imagen(es) dangling", data.Images.Dangling)})
	} else {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "No hay imágenes dangling"})
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Imágenes", Score: score, Stars: getStars(score)}, details
}

func calculateVolumesScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	score := 100
	var details []report.ScoreDetail
	if data.Volumes.Orphaned > 0 {
		penalty := data.Volumes.Orphaned * 15
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d volumen(es) huérfano(s)", data.Volumes.Orphaned)})
	} else {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "No hay volúmenes huérfanos"})
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Volúmenes", Score: score, Stars: getStars(score)}, details
}

func calculateSecurityScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	var details []report.ScoreDetail
	if !data.Security.TrivyInstalled {
		return report.HealthCategory{Name: "Seguridad", Score: -1, Stars: "N/A (Trivy no instalado)"}, details
	}
	score := 100
	if data.Security.TotalVulnerabilities > 0 {
		if data.Security.CriticalCount > 0 {
			penalty := data.Security.CriticalCount * 20
			score -= penalty
			details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d vulnerabilidad(es) CRITICAL", data.Security.CriticalCount)})
		}
		if data.Security.HighCount > 0 {
			penalty := data.Security.HighCount * 5
			score -= penalty
			details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d vulnerabilidad(es) HIGH", data.Security.HighCount)})
		}
	} else {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "Sin vulnerabilidades detectadas"})
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Seguridad", Score: score, Stars: getStars(score)}, details
}

func calculateComposeScore(data report.ReportData) (report.HealthCategory, []report.ScoreDetail) {
	var details []report.ScoreDetail
	if !data.Compose.FileFound {
		return report.HealthCategory{Name: "Docker Compose", Score: -1, Stars: "N/A"}, details
	}
	score := 100
	
	if len(data.Compose.MissingTags) > 0 {
		penalty := len(data.Compose.MissingTags) * 10
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d imagen(es) sin tag específico", len(data.Compose.MissingTags))})
	}
	if len(data.Compose.ExposedPorts) > 0 {
		penalty := len(data.Compose.ExposedPorts) * 20
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d puerto(s) expuesto(s) globalmente", len(data.Compose.ExposedPorts))})
	}
	if len(data.Compose.PrivilegedSvcs) > 0 {
		penalty := len(data.Compose.PrivilegedSvcs) * 30
		score -= penalty
		details = append(details, report.ScoreDetail{Points: -penalty, Reason: fmt.Sprintf("%d servicio(s) en modo privilegiado", len(data.Compose.PrivilegedSvcs))})
	}
	
	if score == 100 {
		details = append(details, report.ScoreDetail{Points: 100, Reason: "Docker Compose configurado correctamente"})
	}

	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Docker Compose", Score: score, Stars: getStars(score)}, details
}
