# 🩺 Docker Doctor

Docker Doctor no es solo un visor de datos, es un **Sistema Experto** diseñado para diagnosticar, monitorear y optimizar tus entornos Docker. Actúa como el experto DevOps de la sala de servidores: analiza tu infraestructura, te da un puntaje de salud y te dice exactamente qué comandos ejecutar para solucionar problemas, ahorrar disco y mitigar riesgos de seguridad.

## ✨ Características Principales

- **📊 Health Score Global (0-100)**: Recibe una calificación general del estado de tu entorno Docker ("Excelente", "Bueno", "Atención" o "Crítico"), además de un sistema de calificación por estrellas (`★★★★★`) para cada categoría (Contenedores, Imágenes, Volúmenes, Seguridad, etc.).
- **📈 Comparativa Histórica (Deltas)**: Almacenamiento local mediante SQLite. Compara automáticamente tu escaneo actual con el anterior y te dice si has mejorado o empeorado tu puntaje (ej. `88 -> 94 (Mejoró 6 pts)`).
- **💡 Motor de Recomendaciones de IA**: No solo detecta el problema, te lo explica. Cada sugerencia incluye:
  - El motivo (¿Por qué importa?).
  - El impacto (Bajo, Medio, Alto) y el riesgo (Rendimiento, Seguridad, Estabilidad).
  - El comando exacto recomendado.
- **💰 Cálculo Real de Ahorro de Disco**: Consulta la API de uso de disco de Docker para decirte exactamente cuántos Megabytes (MB) o Gigabytes (GB) recuperarás al limpiar la basura.
- **🔒 Escaneo de Seguridad y Vulnerabilidades**: Se integra con **Trivy** (si está instalado en tu sistema) para escanear tus imágenes en busca de vulnerabilidades CRITICAL y HIGH.
- **🐳 Analizador de Docker Compose**: Escanea de forma estática los archivos `docker-compose.yml` locales en busca de malas prácticas (imágenes sin tag específico `:latest`, puertos globales expuestos `0.0.0.0` o contenedores en modo privilegiado).
- **🖥️ Dashboard Web Interactivo y Minimalista**: Un servidor integrado que provee una interfaz gráfica hermosa y neutral, construida en React, que muestra tus datos en tiempo real (Live Polling).
- **🚀 Cero Dependencias**: Escrito en Go (Golang). El servidor web de React, las plantillas y la base de datos (pure-go SQLite) están **compilados dentro de un solo archivo binario**. Puedes desplegarlo en cualquier servidor Linux, Windows o Mac sin instalar Node.js ni Python.

---

## 🛠️ Instalación

### Opción 1: Usando Go (Recomendado)
Si tienes Go instalado en tu máquina, puedes instalar Docker Doctor de forma global para que funcione desde cualquier terminal:

```bash
go install
```
*(Esto compilará el código y lo guardará en tu carpeta `go/bin`, que suele estar en tu `PATH`).*

### Opción 2: Compilar el binario manualmente
```bash
# Para Windows:
go build -o docker-doctor.exe main.go

# Para Servidores Linux (Cross-Compilation):
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o docker-doctor-linux main.go
```

---

## 🚀 Uso y Comandos

Una vez instalado (o usando el ejecutable), tienes a tu disposición los siguientes comandos:

### 1. Diagnóstico de Consola (CLI)
Realiza un análisis rápido y genera recomendaciones directamente en tu terminal:
```bash
docker-doctor scan
```

### 2. Exportación a Formatos Externos (CI/CD, Documentación)
Ideal para enviar el reporte a otros sistemas o guardarlo como evidencia:
```bash
# Formato JSON (Ideal para integraciones):
docker-doctor scan --output json --file reporte.json

# Formato Markdown (Ideal para wikis o Notion):
docker-doctor scan --output md --file reporte.md

# Formato HTML (Un archivo autónomo estilizado):
docker-doctor scan --output html --file reporte.html
```

### 3. Dashboard Web en Tiempo Real
Levanta el servidor integrado para explorar los datos gráficamente. Incluye gráficas de consumo, historial, y las tarjetas de recomendaciones de la IA.
```bash
docker-doctor serve
# Por defecto se abre en el puerto 8080 (http://localhost:8080)
```

---

## 🏗️ Arquitectura del Proyecto

- **Golang**: Motor principal, concurrencia, integración con el socket/API de Docker, CLI (con la librería Cobra) y servidor HTTP.
- **React + TypeScript + Vite**: Frontend embebido (`go:embed`) servido dinámicamente.
- **SQLite (Pure-Go)**: Base de datos sin dependencias en C (`CGO_ENABLED=0`), facilitando la compilación cruzada.
- **Trivy CLI**: Integración externa (solo si se detecta en el sistema) para un análisis de seguridad robusto sin engordar el binario.

## 🤝 Contribuir
¡Las contribuciones (pull requests) y las estrellas (⭐) en GitHub son siempre bienvenidas!
