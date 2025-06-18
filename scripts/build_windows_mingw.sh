#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for Windows (MinGW-w64)
# Usage:  ./scripts/build_windows_mingw.sh [VERSION]
# Always builds for x86_64-w64-mingw32 (amd64)
# ------------------------------------------------------------

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -f "${PROJECT_ROOT}/.env" ]] && { set -a; . "${PROJECT_ROOT}/.env"; set +a; }

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH="amd64"
OUT="$PROJECT_ROOT/libs/windows_${ARCH}/${VERSION}"

WRAPPER_DIR="$PROJECT_ROOT/wrapper"
INC_DIR="$PROJECT_ROOT/include/libraw"

LIBRAW_DIR="$PROJECT_ROOT/LibRaw"
if [[ ! -d "$LIBRAW_DIR" ]]; then
  echo "• Cloning LibRaw repository"
  git clone --depth 1 https://github.com/LibRaw/LibRaw.git "$LIBRAW_DIR"
fi

if [[ ! -f "$LIBRAW_DIR/configure" ]]; then
  echo "• Running autoreconf"
  (cd "$LIBRAW_DIR" && autoreconf -fi)
fi

echo "▶ Building rawkit for Windows amd64 → ${OUT}"
mkdir -p "$OUT" "$INC_DIR"

HOST=x86_64-w64-mingw32

pushd "$LIBRAW_DIR" >/dev/null
if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  CC=${HOST}-gcc  CXX=${HOST}-g++  RANLIB=${HOST}-ranlib \
  ./configure --host="$HOST" --disable-shared --enable-static CFLAGS="-O3"
  echo "• Compiling LibRaw"
  make -j"$(nproc || echo 4)"
fi
popd >/dev/null

echo "• Syncing LibRaw headers"
for h in "$LIBRAW_DIR"/libraw/*.h; do
  dst="$INC_DIR/$(basename "$h")"
  [[ -f "$dst" && -z "$(cmp -s "$h" "$dst"; echo $?)" ]] && continue
  cp "$h" "$INC_DIR/"
done
cp "$LIBRAW_DIR/lib/libraw.a" "$OUT/"

echo "• Building libraw_wrapper.a"
${HOST}-g++ -std=c++17 -O3 \
  -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
  -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
  -o "$OUT/libraw_wrapper.o"

${HOST}-ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

ln -sfn "$OUT" "$PROJECT_ROOT/libs/windows_${ARCH}/current"
echo "✔ rawkit done  (current → ${VERSION})"
