package report

import (
	"fmt"
)

func PrintTerminalReport(data ReportData) {
	fmt.Println("========================================")
	fmt.Println("         DOCKER DOCTOR REPORT           ")
	fmt.Println("========================================")
	
	if data.System.IsReachable {
		fmt.Println("✅ Sistema Docker: Operativo y Accesible")
	} else {
		fmt.Println("❌ Sistema Docker: Inaccesible (¿Está corriendo el daemon?)")
	}

	fmt.Printf("📦 Contenedores: %d en total, %d detenidos\n", data.Containers.Total, data.Containers.Stopped)
	fmt.Printf("🖼️  Imágenes: %d en total, %d dangling\n", data.Images.Total, data.Images.Dangling)
	fmt.Printf("💾 Volúmenes: %d en total, %d huérfanos\n", data.Volumes.Total, data.Volumes.Orphaned)
	fmt.Printf("🌐 Redes: %d en total, %d sin uso\n", data.Networks.Total, data.Networks.Unused)
	fmt.Printf("🚪 Puertos Expuestos: %d\n", data.Ports.TotalExposed)
	fmt.Println("========================================")
}
