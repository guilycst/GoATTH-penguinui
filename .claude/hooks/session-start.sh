#!/bin/bash
set -euo pipefail

# Only run in remote (Claude Code on the web) environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"

# Go version required by the project
GO_VERSION="1.26.1"
REQUIRED_MAJOR=1
REQUIRED_MINOR=26

# Check if Go upgrade is needed
CURRENT_GO_VERSION=$(go version 2>/dev/null | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
CURRENT_MAJOR=$(echo "$CURRENT_GO_VERSION" | cut -d. -f1)
CURRENT_MINOR=$(echo "$CURRENT_GO_VERSION" | cut -d. -f2)

if [ "$CURRENT_MAJOR" -lt "$REQUIRED_MAJOR" ] || { [ "$CURRENT_MAJOR" -eq "$REQUIRED_MAJOR" ] && [ "$CURRENT_MINOR" -lt "$REQUIRED_MINOR" ]; }; then
  echo "==> Installing Go ${GO_VERSION} (current: ${CURRENT_GO_VERSION})..."
  curl -sL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
  rm /tmp/go.tar.gz
  export PATH="/usr/local/go/bin:$PATH"
  echo "==> Go installed: $(go version)"
fi

echo "==> Installing Go module dependencies..."
cd "$PROJECT_DIR"
go mod download

echo "==> Installing templ CLI..."
go install github.com/a-h/templ/cmd/templ@latest

echo "==> Installing Tailwind CSS CLI (standalone)..."
if ! command -v tailwindcss &>/dev/null; then
  curl -sL "https://github.com/nicolo-ribaudo/tailwindcss-cli/releases/download/v4.1.4/tailwindcss-linux-x64" -o /usr/local/bin/tailwindcss
  chmod +x /usr/local/bin/tailwindcss
fi

echo "==> Generating templ files..."
export PATH="$PATH:$(go env GOPATH)/bin"
cd "$PROJECT_DIR"
templ generate

echo "==> Building CSS..."
tailwindcss -i css/main.css -o assets/styles.css

echo "==> Persisting PATH for session..."
echo "export PATH=\"/usr/local/go/bin:\$PATH:$(go env GOPATH)/bin\"" >> "$CLAUDE_ENV_FILE"

echo "==> Session start hook completed successfully."
