#!/usr/bin/env bash
set -euo pipefail

GO_VERSION="1.24.0"
GO_ARCH="linux-amd64"
TARBALL="go${GO_VERSION}.${GO_ARCH}.tar.gz"

if [ ! -f "$TARBALL" ]; then
  echo "Downloading Go $GO_VERSION"
  wget -q "https://go.dev/dl/${TARBALL}"
fi

sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "$TARBALL"
rm "$TARBALL"

# Add Go to PATH for current session
export PATH="/usr/local/go/bin:$PATH"

echo "Go $GO_VERSION installed."
