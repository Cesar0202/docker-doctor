# Docker Doctor - Alcance del Proyecto

Este documento detalla las funcionalidades incluidas en el Producto Mínimo Viable (MVP) y la visión completa para el futuro del proyecto **Docker Doctor**.

## MVP (Producto Mínimo Viable)

El MVP se centra en funcionar como una herramienta de diagnóstico por CLI rápida y eficiente para detectar problemas comunes y usos excesivos de recursos.

1. **Análisis del sistema Docker**
   - Verificación de instalación de Docker, Docker Engine y Docker Compose.
   - Detección de versión y estado del daemon.

2. **Diagnóstico de Contenedores**
   - Detección de contenedores detenidos, reiniciándose constantemente, huérfanos o duplicados.
   - Monitoreo básico de consumo de CPU y RAM.

3. **Diagnóstico de Imágenes**
   - Detección de imágenes sin utilizar (dangling), duplicadas, de gran tamaño (enormes), muy antiguas y sin etiquetas.

4. **Diagnóstico de Volúmenes**
   - Identificación de volúmenes huérfanos (no adjuntos a ningún contenedor).
   - Cálculo del espacio ocupado y fecha de última utilización.

5. **Diagnóstico de Redes**
   - Detección de redes no utilizadas, redes duplicadas y posibles conflictos.

6. **Diagnóstico de Puertos**
   - Mapeo de puertos expuestos y detección de puertos ocupados (ej. 80, 3306, 5432).

7. **Análisis de Espacio en Disco**
   - Desglose total del espacio ocupado por Docker (Imágenes, Volúmenes, Build Cache, Logs).

8. **Generación de Reportes**
   - Exportación de los resultados del análisis en formatos: **HTML**, **JSON** y **Markdown**.

---

## Proyecto Completo (Funcionalidades Futuras)

El proyecto completo evolucionará de un analizador estático a un ecosistema de monitoreo, optimización y seguridad continuo.

1. **Análisis de Vulnerabilidades**
   - Integración con herramientas como **Trivy** o **Docker Scout** para escanear vulnerabilidades en las imágenes instaladas.

2. **Motor de Recomendaciones Inteligentes**
   - Sugerencias accionables tras el análisis (ej. advertir sobre imágenes antiguas).
   - Estimación de ahorro de espacio e impresión del comando exacto sugerido para solucionar el problema (ej. `docker image prune`).

3. **Historial y Seguimiento**
   - Uso de una base de datos local (SQLite) para guardar los análisis a lo largo del tiempo.
   - Comparativas visuales del crecimiento del uso del disco (ej. Hoy vs. Hace un mes).

4. **Dashboard Web interactivo (Fase 2)**
   - Interfaz gráfica (Backend en Go, Frontend en React/Vite/Tailwind CSS).
   - Visualización amigable de métricas de CPU, RAM, Disco, alertas e historial.

5. **Monitor en Tiempo Real (TUI)**
   - Modo de visualización en la terminal similar a `htop`.
   - Actualización de métricas en vivo cada segundo.

6. **Docker Compose Analyzer**
   - Lectura y validación estática de archivos `docker-compose.yml`.
   - Detección de variables faltantes, puertos repetidos, imágenes inexistentes y antipatrones/malas prácticas.

7. **CI/CD y Distribución Global**
   - Flujos de trabajo automatizados con **GitHub Actions**.
   - Generación y distribución de binarios multiplataforma listos para usar (Windows, Linux, macOS) sin dependencias adicionales.
