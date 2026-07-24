# Docker Doctor

Docker Doctor no es solo un visor de datos, es un **Sistema Experto** diseñado para diagnosticar, monitorear y optimizar tus entornos Docker. Actúa como el experto DevOps de la sala de servidores: analiza tu infraestructura, te da un puntaje de salud y te dice exactamente qué comandos ejecutar para solucionar problemas, ahorrar disco y mitigar riesgos de seguridad.

## Características Principales

- **Health Score Global (0-100)**: Recibe una calificación general del estado de tu entorno Docker ("Excelente", "Bueno", "Atención" o "Crítico"), además de un sistema de calificación por estrellas (`★★★★★`) para cada categoría (Contenedores, Imágenes, Volúmenes, Seguridad, etc.).
- **Comparativa Histórica (Deltas)**: Almacenamiento local mediante SQLite. Compara automáticamente tu escaneo actual con el anterior y te dice si has mejorado o empeorado tu puntaje (ej. `88 -> 94 (Mejoró 6 pts)`).
- **Asistente Proactivo Interactivo**: Permite detectar basura y limpiar el entorno automáticamente con el comando `fix`, interpretando además errores complejos usando una base de datos de conocimiento con `explain`.
- **Motor de Recomendaciones**: No solo detecta el problema, te lo explica. Cada sugerencia incluye el motivo, el riesgo (Rendimiento, Seguridad, Estabilidad) y un nivel de prioridad visual.
- **Cálculo Real de Ahorro de Disco**: Consulta la API de uso de disco de Docker para decirte exactamente cuántos Megabytes (MB) o Gigabytes (GB) recuperarás al limpiar la basura.
- **Escaneo de Puertos Locales**: Detecta si hay servicios locales (como PostgreSQL en el puerto 5432) que puedan causar conflictos antes de levantar un contenedor.
- **Analizador de Docker Compose**: Escanea de forma estática los archivos `docker-compose.yml` locales en busca de malas prácticas.
- **Dashboard Web Interactivo y Minimalista**: Un servidor integrado que provee una interfaz gráfica hermosa y neutral, construida en React, que muestra tus datos en tiempo real (Live Polling).
- **Cero Dependencias**: Escrito en Go (Golang). El servidor web de React, las plantillas y la base de datos (pure-go SQLite) están **compilados dentro de un solo archivo binario**. Puedes desplegarlo en cualquier servidor Linux, Windows o Mac sin instalar Node.js ni Python.

---

## Instalación

### Opción 1: Instalar directamente el Binario (Recomendado)
**¡No necesitas tener Go instalado!** Docker Doctor es un binario independiente. Puedes instalarlo en cualquier servidor Linux o macOS en segundos con nuestro script de instalación rápida:

```bash
curl -sSL https://raw.githubusercontent.com/Cesar0202/docker-doctor/main/install.sh | bash
```
*(Este script detecta tu sistema operativo y arquitectura, descarga el ejecutable oficial desde GitHub Releases y lo coloca en `/usr/local/bin` para que puedas usarlo de inmediato).*

### Opción 2: Compilar desde el código fuente (Desarrolladores)
Si prefieres compilarlo tú mismo o ya tienes un entorno de Go (`go1.21+`) configurado, puedes usar:

```bash
go install github.com/Cesar0202/docker-doctor@latest
```
*(Esto compilará el código y lo guardará en tu carpeta `go/bin`, que suele estar en tu `PATH`).*

---

## Uso y Comandos

Una vez instalado, tienes a tu disposición los siguientes comandos principales:

### 1. Diagnóstico de Consola (CLI)
Realiza un análisis rápido y genera recomendaciones de expertos directamente en tu terminal:
```bash
docker-doctor scan
```

### 2. Asistente de Reparación Interactiva
Limpia tu entorno automáticamente sin necesidad de teclear comandos manuales de prune:
```bash
docker-doctor fix
```

### 3. Experto Solucionador de Errores
Explica errores crípticos de Docker (como `manifest unknown` o `conflict`) y te da comandos sugeridos:
```bash
docker-doctor explain "mensaje de error de docker"
```

### 4. Historial y Tendencias
Muestra una tabla con el progreso de tu salud a lo largo del tiempo, o genera un gráfico ASCII en la terminal:
```bash
docker-doctor history
docker-doctor history --trend
```

### 5. Exportación a Formatos Externos (CI/CD, Documentación)
Ideal para enviar el reporte a otros sistemas o guardarlo como evidencia:
```bash
docker-doctor scan --output json --file reporte.json
docker-doctor scan --output md --file reporte.md
docker-doctor scan --output html --file reporte.html
```

### 6. Dashboard Web en Tiempo Real
Levanta el servidor integrado para explorar los datos gráficamente. Incluye gráficas de consumo, historial y recomendaciones.
```bash
docker-doctor web
# Por defecto se abre en el puerto 8080 (http://localhost:8080)

# Puedes usar otro puerto si el 8080 está ocupado:
docker-doctor web --port 3000
```

---

## Arquitectura del Proyecto

- **Golang**: Motor principal, concurrencia, integración con el socket/API de Docker, CLI (con la librería Cobra) y servidor HTTP.
- **React + TypeScript + Vite**: Frontend embebido (`go:embed`) servido dinámicamente.
- **SQLite (Pure-Go)**: Base de datos sin dependencias en C (`CGO_ENABLED=0`), facilitando la compilación cruzada.

## Contribuir
¡Las contribuciones (pull requests) y las estrellas en GitHub son siempre bienvenidas!
