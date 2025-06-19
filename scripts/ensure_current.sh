#!/usr/bin/env bash
# Ensures libs/<os>/<arch>/current/ is populated, either by downloading
# or by building from source if toolchain is present.

set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[[ $ARCH == x86_64 ]] && ARCH=amd64
[[ $ARCH == arm64  ]] && ARCH=arm64

CUR="libs/${OS}_${ARCH}/current"
if [[ -f "${PROJECT_ROOT}/${CUR}/libraw.a" ]]; then
  echo "✓ ${CUR} already present"
  exit 0
fi

# try toolchain build first
BUILD_SCRIPT="${PROJECT_ROOT}/scripts/build_${OS}.sh"
if [[ -x "$BUILD_SCRIPT" ]]; then
  echo "🔧 building LibRaw for ${OS}_${ARCH}"
  "$BUILD_SCRIPT" v0.0.0 "${ARCH}"
else
  echo "⤵️  downloading pre-built LibRaw for ${OS}_${ARCH}"
  "${PROJECT_ROOT}/scripts/fetch_prebuilt.sh" "${OS}" "${ARCH}"
fi
