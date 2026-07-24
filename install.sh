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

# Ajuste si es darwin (macOS) - aunque ahorita solo compilamos linux y windows
# Lo dejamos preparado para el futuro.
if [ "$OS" = "darwin" ]; then
    echo "Aún no hay binarios para macOS. Por favor usa 'go install github.com/Cesar0202/docker-doctor@latest'"
    exit 1
fi

REPO="Cesar0202/docker-doctor"
BINARY_NAME="docker-doctor"
FILENAME="${BINARY_NAME}-${OS}-${ARCH}"

DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/${FILENAME}"

echo "Detectado: $OS-$ARCH"
echo "Descargando desde: $DOWNLOAD_URL"

# Descargar el binario en un archivo temporal
TMP_FILE=$(mktemp)
curl -sSL -f -o "$TMP_FILE" "$DOWNLOAD_URL" || {
    echo "Error: No se pudo descargar el binario (¿Revisaste si el Release existe en GitHub?)"
    rm -f "$TMP_FILE"
    exit 1
}

# Dar permisos de ejecución
chmod +x "$TMP_FILE"

# Mover a /usr/local/bin
INSTALL_DIR="/usr/local/bin"
echo "Instalando en $INSTALL_DIR (puede pedir contraseña de sudo)..."
sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "¡Instalación completada con éxito! 🎉"
echo "Puedes probarlo ejecutando: docker-doctor"
