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
	fmt.Println("========================================")
}
