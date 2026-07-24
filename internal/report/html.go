package report

import (
	"fmt"
	"html/template"
	"os"
	"docker-doctor/internal/db"
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
        .health-score { font-size: 2em; text-align: center; margin: 20px 0; font-weight: bold; }
        .health-score.good { color: #28a745; }
        .health-score.warn { color: #ffc107; }
        .health-score.crit { color: #dc3545; }
        .recs { margin-top: 30px; }
        .rec-item { background: #fff; padding: 15px; border-left: 4px solid #17a2b8; margin-bottom: 15px; border-radius: 4px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); }
        .rec-item code { display: block; background: #e9ecef; padding: 10px; margin-top: 10px; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Docker Doctor Report</h1>
        
        <div class="card {{if .System.IsReachable}}success{{else}}danger{{end}}">
            <h3>Estado del Sistema</h3>
            <p>{{if .System.IsReachable}}[OK] Operativo y Accesible{{else}}[FAIL] Inaccesible{{end}}</p>
        </div>

        <div class="grid">
            <div class="card">
                <h3>Contenedores</h3>
                <p>Total: <strong>{{.Containers.Total}}</strong></p>
                <p>Detenidos: <strong>{{.Containers.Stopped}}</strong></p>
            </div>
            <div class="card">
                <h3>Imágenes</h3>
                <p>Total: <strong>{{.Images.Total}}</strong></p>
                <p>Dangling: <strong>{{.Images.Dangling}}</strong></p>
            </div>
            <div class="card">
                <h3>Volúmenes</h3>
                <p>Total: <strong>{{.Volumes.Total}}</strong></p>
                <p>Huérfanos: <strong>{{.Volumes.Orphaned}}</strong></p>
            </div>
            <div class="card">
                <h3>Redes</h3>
                <p>Total: <strong>{{.Networks.Total}}</strong></p>
                <p>Sin Uso: <strong>{{.Networks.Unused}}</strong></p>
            </div>
            <div class="card">
                <h3>Puertos</h3>
                <p>Total Expuestos: <strong>{{.ReportData.Ports.TotalExposed}}</strong></p>
            </div>
        </div>

        <h2 style="text-align:center; margin-top:30px;">Health Score</h2>
        <div class="health-score {{if ge .Health.GlobalScore 90}}good{{else if ge .Health.GlobalScore 70}}warn{{else}}crit{{end}}">
            {{.Health.GlobalScore}}/100 ({{.Health.StatusText}})
        </div>

        <div class="recs">
            <h2>Recomendaciones de la IA</h2>
            {{if .Recommendations}}
                {{range .Recommendations}}
                <div class="rec-item">
                    <h4>[{{.Level}}] {{.Message}}</h4>
                    {{if .Why}}<p><strong>¿Por qué importa?</strong> {{.Why}}</p>{{end}}
                    {{if .Impact}}<p><strong>Impacto:</strong> {{.Impact}} | <strong>Riesgo:</strong> {{.Risk}}</p>{{end}}
                    <code>{{.Command}}</code>
                </div>
                {{end}}
            {{else}}
                <p>¡Tu entorno está limpio y optimizado! No hay recomendaciones.</p>
            {{end}}
        </div>
    </div>
</body>
</html>`

func ExportHTML(data ReportData, hr HealthReport, lastScan db.ScanHistory, recs []Recommendation, filename string) error {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parseando plantilla HTML: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo HTML: %w", err)
	}
	defer file.Close()

	exportPayload := struct {
		ReportData
		Health          HealthReport
		LastScan        db.ScanHistory
		Recommendations []Recommendation
	}{
		ReportData:      data,
		Health:          hr,
		LastScan:        lastScan,
		Recommendations: recs,
	}

	err = tmpl.Execute(file, exportPayload)
	if err != nil {
		return fmt.Errorf("error ejecutando plantilla HTML: %w", err)
	}

	fmt.Printf("[OK] Reporte HTML exportado exitosamente a %s\n", filename)
	return nil
}
