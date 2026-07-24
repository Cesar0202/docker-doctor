package report

import (
	"fmt"
	"docker-doctor/internal/analyzer"
)

func PrintTerminalReport(sys analyzer.SystemStatus, cont analyzer.ContainerAnalysis, img analyzer.ImageAnalysis) {
	fmt.Println("========================================")
	fmt.Println("         DOCKER DOCTOR REPORT           ")
	fmt.Println("========================================")
	
	if sys.IsReachable {
		fmt.Println("✅ Sistema Docker: Operativo y Accesible")
	} else {
		fmt.Println("❌ Sistema Docker: Inaccesible (¿Está corriendo el daemon?)")
	}

	fmt.Printf("📦 Contenedores: %d en total, %d detenidos\n", cont.Total, cont.Stopped)
	fmt.Printf("🖼️  Imágenes: %d en total, %d dangling\n", img.Total, img.Dangling)
	fmt.Println("========================================")
}
