#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for Windows-amd64 (MinGW-w64)
# Usage:  ./scripts/build_windows_mingw.sh [VERSION]
# ------------------------------------------------------------

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -f "${PROJECT_ROOT}/.env" ]] && { set -a; . "${PROJECT_ROOT}/.env"; set +a; }

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH=amd64
HOST=x86_64-w64-mingw32

OUT="${PROJECT_ROOT}/libs/windows_${ARCH}/${VERSION}"
CURRENT_DIR="${PROJECT_ROOT}/libs/windows_${ARCH}/current"
WRAPPER_DIR="${PROJECT_ROOT}/wrapper"
INC_DIR="${PROJECT_ROOT}/include/libraw"
LIBRAW_DIR="${PROJECT_ROOT}/LibRaw"

echo "▶ Building rawkit for Windows ${ARCH} → ${OUT}"
mkdir -p "${OUT}" "${INC_DIR}"

# ── 1. build LibRaw ──────────────────────────────────────────
pushd "${LIBRAW_DIR}" >/dev/null

if [[ ! -f configure ]]; then
  echo "• Generating configure script"
  make -f Makefile.dist
fi

echo "• Configuring LibRaw"
CC=${HOST}-gcc  CXX=${HOST}-g++  RANLIB=${HOST}-ranlib \
./configure --host="${HOST}" \
            --disable-shared --enable-static --disable-examples \
            CFLAGS="-O3" \
            LDFLAGS="-lz -static" >/dev/null

echo "• Compiling LibRaw (library only)"
make -C lib -j"$(nproc 2>/dev/null || echo 1)" libraw.a >/dev/null
popd >/dev/null

echo "• Syncing LibRaw headers"
for h in "${LIBRAW_DIR}"/libraw/*.h; do
  dst="${INC_DIR}/$(basename "$h")"
  [[ -f "$dst" && cmp -s "$h" "$dst" ]] || cp "$h" "$dst"
done

cp "${LIBRAW_DIR}/lib/libraw.a" "${OUT}/"

echo "• Building libraw_wrapper.a"
${HOST}-g++ -std=c++17 -O3 \
  -I"${INC_DIR}" -I"${LIBRAW_DIR}/libraw" \
  -c "${WRAPPER_DIR}/libraw_wrapper.cpp" \
  -o "${OUT}/libraw_wrapper.o"

${HOST}-ar rcs "${OUT}/libraw_wrapper.a" "${OUT}/libraw_wrapper.o"
rm "${OUT}/libraw_wrapper.o"

echo "• Updating current/ directory"
rm -rf "${CURRENT_DIR}"
mkdir -p "${CURRENT_DIR}"
cp "${OUT}"/*.a "${CURRENT_DIR}/"

echo "✔ rawkit done  (current folder refreshed → ${VERSION})"
