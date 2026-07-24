#!/bin/bash
set -e

echo "========================================="
echo "   Instalando Docker Doctor..."
echo "========================================="

# Detectar OS y Arquitectura
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Arquitectura no soportada: $ARCH"
    exit 1
fi

REPO="Cesar0202/docker-doctor"
BINARY_NAME="docker-doctor"

echo "Detectado: $OS-$ARCH"

# Aquí normalmente se descargaría el binario desde los Releases de GitHub.
# Ejemplo:
# DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
# curl -sSL -o $BINARY_NAME $DOWNLOAD_URL

echo ""
echo "[AVISO PARA CESAR]: Aún no has creado un 'Release' en GitHub con los binarios compilados."
echo "Por ahora, este script intentará usar 'go install' si Go está instalado."
echo ""

if command -v go &> /dev/null; then
    echo "-> Go está instalado. Compilando..."
    go install github.com/$REPO@latest
    echo "¡Instalado con éxito usando Go!"
else
    echo "ERROR: No se encontró Go en este sistema y aún no hay binarios pre-compilados en GitHub Releases."
    echo "Por favor, para probarlo en este servidor sin instalar Go, debes compilarlo en tu PC local (GOOS=linux GOARCH=amd64 go build) y subir el archivo al servidor."
    exit 1
fi
