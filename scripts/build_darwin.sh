#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for macOS
# Usage:  ./scripts/build_darwin.sh [VERSION] [ARCH]
#         VERSION defaults to the one in .env (or v0.0.1)
#         ARCH    defaults to your host  (arm64 on Apple Silicon)
# ------------------------------------------------------------

# ---------- 0. load .env (if present) ------------------------
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -f "${PROJECT_ROOT}/.env" ]]; then
  set -a; . "${PROJECT_ROOT}/.env"; set +a
fi

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH="${2:-$(uname -m)}"
OUT="$PROJECT_ROOT/libs/darwin_${ARCH}/${VERSION}"
WRAPPER_DIR="$PROJECT_ROOT/wrapper"
INC_DIR="$PROJECT_ROOT/include/libraw"
LIBRAW_DIR="$PROJECT_ROOT/LibRaw"

echo "▶ Building rawkit for macOS ${ARCH} → ${OUT}"
mkdir -p "$OUT" "$INC_DIR"

pushd "$LIBRAW_DIR" >/dev/null

if [[ ! -f configure ]]; then
  echo "• Generating configure script"
  make -f Makefile.dist
fi

if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  ./configure --disable-shared --enable-static CFLAGS="-O3" >/dev/null
  echo "• Compiling LibRaw"
  make -j"$(sysctl -n hw.ncpu)"    >/dev/null
fi
popd >/dev/null

echo "• Syncing LibRaw public headers"
for h in "$LIBRAW_DIR"/libraw/*.h; do
  dst="$INC_DIR/$(basename "$h")"
  if [[ ! -f "$dst" ]] || ! cmp -s "$h" "$dst"; then
    cp "$h" "$dst"
  fi
done

cp "$LIBRAW_DIR/lib/libraw.a" "$OUT/"

echo "• Building libraw_wrapper.a"
c++ -std=c++17 -O3 \
    -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
    -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
    -o "$OUT/libraw_wrapper.o"

ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

echo "✔ rawkit done: $OUT/{libraw.a,libraw_wrapper.a}"
