package report

import (
	"fmt"
	"html/template"
	"os"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docker Doctor Report</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f7f6; color: #333; margin: 0; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: #fff; padding: 30px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        h1 { color: #007bff; text-align: center; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-top: 20px; }
        .card { background: #f8f9fa; padding: 20px; border-radius: 6px; border-left: 4px solid #007bff; }
        .card h3 { margin-top: 0; color: #495057; }
        .success { border-color: #28a745; }
        .warning { border-color: #ffc107; }
        .danger { border-color: #dc3545; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Docker Doctor Report 🩺</h1>
        
        <div class="card {{if .System.IsReachable}}success{{else}}danger{{end}}">
            <h3>Estado del Sistema</h3>
            <p>{{if .System.IsReachable}}✅ Operativo y Accesible{{else}}❌ Inaccesible{{end}}</p>
        </div>

        <div class="grid">
            <div class="card">
                <h3>📦 Contenedores</h3>
                <p>Total: <strong>{{.Containers.Total}}</strong></p>
                <p>Detenidos: <strong>{{.Containers.Stopped}}</strong></p>
            </div>
            <div class="card">
                <h3>🖼️ Imágenes</h3>
                <p>Total: <strong>{{.Images.Total}}</strong></p>
                <p>Dangling: <strong>{{.Images.Dangling}}</strong></p>
            </div>
            <div class="card">
                <h3>💾 Volúmenes</h3>
                <p>Total: <strong>{{.Volumes.Total}}</strong></p>
                <p>Huérfanos: <strong>{{.Volumes.Orphaned}}</strong></p>
            </div>
            <div class="card">
                <h3>🌐 Redes</h3>
                <p>Total: <strong>{{.Networks.Total}}</strong></p>
                <p>Sin Uso: <strong>{{.Networks.Unused}}</strong></p>
            </div>
            <div class="card">
                <h3>🚪 Puertos</h3>
                <p>Total Expuestos: <strong>{{.Ports.TotalExposed}}</strong></p>
            </div>
        </div>
    </div>
</body>
</html>`

func ExportHTML(data ReportData, filename string) error {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parseando plantilla HTML: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo HTML: %w", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error ejecutando plantilla HTML: %w", err)
	}

	fmt.Printf("✅ Reporte HTML exportado exitosamente a %s\n", filename)
	return nil
}
