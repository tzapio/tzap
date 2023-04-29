#!/bin/bash

GITHUB_USER="tzapio"
GITHUB_REPO="tzap"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Mapping uname output to correct architecture names
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

# Get the latest version number from GitHub API
VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ "$OS" == "darwin" ] || [ "$OS" == "linux" ]; then
    FILENAME="tzap-${VERSION}-${OS}-${ARCH}"
    URL="https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/download/${VERSION}/${FILENAME}"
    echo "Downloading Tzap ${VERSION} for ${OS}-${ARCH}..."
    curl -L -o tzap "${URL}"
    chmod +x tzap
    sudo mv tzap /usr/local/bin/
    echo "Tzap ${VERSION} installed successfully."
elif [ "$OS" == "windows" ]; then
    FILENAME="tzap-${VERSION}-windows-${ARCH}.exe"
    URL="https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/download/${VERSION}/${FILENAME}"
    echo "Downloading Tzap ${VERSION} for ${OS}-${ARCH}..."
    curl -L -o tzap.exe "${URL}"
    mv tzap.exe %USERPROFILE%\\AppData\\Local\\Programs\\tzap\\
    setx PATH "%PATH%;%USERPROFILE%\\AppData\\Local\\Programs\\tzap\\"
    echo "Tzap ${VERSION} installed successfully."
else
    echo "Error: Unsupported operating system."
    exit 1
fi