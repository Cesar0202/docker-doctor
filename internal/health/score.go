package health

import (
	"docker-doctor/internal/report"
)

func CalculateHealth(data report.ReportData) report.HealthReport {
	hr := report.HealthReport{}

	// Calcular estrellas por categoría
	catEngine := calculateEngineScore(data)
	catContainers := calculateContainersScore(data)
	catImages := calculateImagesScore(data)
	catVolumes := calculateVolumesScore(data)
	catSecurity := calculateSecurityScore(data)
	catCompose := calculateComposeScore(data)

	hr.Categories = []report.HealthCategory{
		catEngine, catContainers, catImages, catVolumes, catSecurity, catCompose,
	}

	// Score global es el promedio (excluyendo compose si no aplica)
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

func calculateEngineScore(data report.ReportData) report.HealthCategory {
	score := 100
	if !data.System.IsReachable {
		score = 0
	}
	return report.HealthCategory{Name: "Docker Engine", Score: score, Stars: getStars(score)}
}

func calculateContainersScore(data report.ReportData) report.HealthCategory {
	score := 100
	if data.Containers.Stopped > 0 {
		score -= (data.Containers.Stopped * 10)
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Contenedores", Score: score, Stars: getStars(score)}
}

func calculateImagesScore(data report.ReportData) report.HealthCategory {
	score := 100
	if data.Images.Dangling > 0 {
		score -= (data.Images.Dangling * 15)
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Imágenes", Score: score, Stars: getStars(score)}
}

func calculateVolumesScore(data report.ReportData) report.HealthCategory {
	score := 100
	if data.Volumes.Orphaned > 0 {
		score -= (data.Volumes.Orphaned * 15)
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Volúmenes", Score: score, Stars: getStars(score)}
}

func calculateSecurityScore(data report.ReportData) report.HealthCategory {
	if !data.Security.TrivyInstalled {
		return report.HealthCategory{Name: "Seguridad", Score: 50, Stars: "N/A (Trivy no instalado)"}
	}
	score := 100
	if data.Security.TotalVulnerabilities > 0 {
		score -= (data.Security.CriticalCount * 20)
		score -= (data.Security.HighCount * 5)
	}
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Seguridad", Score: score, Stars: getStars(score)}
}

func calculateComposeScore(data report.ReportData) report.HealthCategory {
	if !data.Compose.FileFound {
		return report.HealthCategory{Name: "Docker Compose", Score: -1, Stars: "N/A"}
	}
	score := 100
	score -= (len(data.Compose.MissingTags) * 10)
	score -= (len(data.Compose.ExposedPorts) * 20)
	score -= (len(data.Compose.PrivilegedSvcs) * 30)
	if score < 0 { score = 0 }
	return report.HealthCategory{Name: "Docker Compose", Score: score, Stars: getStars(score)}
}
