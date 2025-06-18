#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -f "${PROJECT_ROOT}/.env" ]]; then
  set -a; . "${PROJECT_ROOT}/.env"; set +a
fi

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH=amd64
OUT="$PROJECT_ROOT/libs/windows_${ARCH}/${VERSION}"
WRAPPER_DIR="$PROJECT_ROOT/wrapper"
INC_DIR="$PROJECT_ROOT/include/libraw"
LIBRAW_DIR="$PROJECT_ROOT/LibRaw"
HOST="x86_64-w64-mingw32"

echo "▶ Building rawkit for Windows ${ARCH} → ${OUT}"
mkdir -p "$OUT" "$INC_DIR"

pushd "$LIBRAW_DIR" >/dev/null
if [[ ! -f configure ]]; then
  echo "• Generating configure script"
  make -f Makefile.dist
fi
if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  CC=${HOST}-gcc  CXX=${HOST}-g++  RANLIB=${HOST}-ranlib \
  ./configure --host="$HOST" --disable-shared --enable-static CFLAGS="-O3" LDFLAGS="-lz"
  echo "• Compiling LibRaw"
  make -j"$(sysctl -n hw.ncpu 2>/dev/null || echo 1)"
fi
popd >/dev/null

echo "• Syncing LibRaw headers"
for h in "$LIBRAW_DIR"/libraw/*.h; do
  dst="$INC_DIR/$(basename "$h")"
  [[ ! -f "$dst" ]] || cmp -s "$h" "$dst" && continue
  cp "$h" "$dst"
done
cp "$LIBRAW_DIR/lib/libraw.a" "$OUT/"

echo "• Building libraw_wrapper.a"
${HOST}-g++ -std=c++17 -O3 \
  -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
  -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
  -o "$OUT/libraw_wrapper.o"

${HOST}-ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

echo "✔ rawkit done: $OUT/"
