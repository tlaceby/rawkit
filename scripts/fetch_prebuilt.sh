#!/usr/bin/env bash
set -euo pipefail
OS=$1   # linux | darwin | windows
ARCH=$2 # amd64 | arm64
TAG=${3:-latest}

ASSET="${OS}_${ARCH}.tar.gz"

echo "⤵️  Fetching ${ASSET} from release ${TAG}"
curl -sL -o "${ASSET}" \
  "https://github.com/tlaceby/rawkit/releases/download/${TAG}/${ASSET}"

echo "• unpacking"
DIR="libs/${OS}_${ARCH}/current"
rm -rf "$DIR"
mkdir -p "$DIR"
tar -xzf "${ASSET}" -C "$DIR"
rm "${ASSET}"
