package report

import (
	"fmt"
)

func PrintTerminalReport(data ReportData) {
	fmt.Println("========================================")
	fmt.Println("         DOCKER DOCTOR REPORT           ")
	fmt.Println("========================================")
	
	if data.System.IsReachable {
		fmt.Println("[OK] Sistema Docker: Operativo y Accesible")
	} else {
		fmt.Println("[FAIL] Sistema Docker: Inaccesible (¿Está corriendo el daemon?)")
	}

	fmt.Printf("[INFO] Contenedores: %d en total, %d detenidos\n", data.Containers.Total, data.Containers.Stopped)
	fmt.Printf("[INFO] Imágenes: %d en total, %d dangling\n", data.Images.Total, data.Images.Dangling)
	fmt.Printf("[INFO] Volúmenes: %d en total, %d huérfanos\n", data.Volumes.Total, data.Volumes.Orphaned)
	fmt.Printf("[INFO] Redes: %d en total, %d sin uso\n", data.Networks.Total, data.Networks.Unused)
	fmt.Printf("[INFO] Puertos Expuestos: %d\n", data.Ports.TotalExposed)

	if data.Security.TrivyInstalled {
		if data.Security.TotalVulnerabilities == 0 {
			fmt.Printf("[SEC]  Trivy: %d imágenes escaneadas. 0 vulnerabilidades.\n", data.Security.ScannedImages)
		} else {
			fmt.Printf("[SEC]  Trivy: %d imgs escaneadas. %d Vulnerabilidades (CRITICAL: %d, HIGH: %d)\n", 
				data.Security.ScannedImages, data.Security.TotalVulnerabilities, 
				data.Security.CriticalCount, data.Security.HighCount)
		}
	} else {
		fmt.Println("[SEC]  Trivy: No instalado. Omite escaneo de vulnerabilidades.")
	}

	if data.Compose.FileFound {
		fmt.Printf("[YML]  Docker Compose analizado: %d servicios encontrados.\n", data.Compose.TotalServices)
	}

	fmt.Println("========================================")
}
