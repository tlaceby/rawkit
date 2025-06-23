#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for macOS
# Usage:  ./scripts/build_darwin.sh [VERSION] [ARCH]
#         VERSION defaults to .env VERSION or v0.0.1
#         ARCH    defaults to uname -m
# ------------------------------------------------------------

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -f "$PROJECT_ROOT/.env" ]]; then
  set -a
  source "$PROJECT_ROOT/.env"
  set +a
fi

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH="${2:-$(uname -m)}"
OUT="$PROJECT_ROOT/libs/darwin_${ARCH}/${VERSION}"
CURRENT="$PROJECT_ROOT/libs/darwin_${ARCH}/current"
WRAPPER_DIR="$PROJECT_ROOT/wrapper"
INC_DIR="$PROJECT_ROOT/include/libraw"
LIBRAW_DIR="$PROJECT_ROOT/LibRaw"

mkdir -p LibRaw/object
mkdir -p LibRaw/bin
mkdir -p "$OUT" "$INC_DIR"
echo "▶ Building rawkit for macOS $ARCH → $OUT"

pushd "$LIBRAW_DIR" > /dev/null

if [[ ! -f configure ]]; then
  echo "• Generating configure script"
  make -f Makefile.dist
fi

if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  ./configure --disable-shared --enable-static CFLAGS="-O3" > /dev/null
  echo "• Compiling LibRaw"
  make -j"$(sysctl -n hw.ncpu)" > /dev/null
else
  echo "• LibRaw already built — skipping"
fi

popd > /dev/null

echo "• Syncing LibRaw headers"
for h in "$LIBRAW_DIR"/libraw/*.h; do
  cp "$h" "$INC_DIR/"
done

cp "$LIBRAW_DIR/lib/libraw.a" "$OUT/"

echo "• Building libraw_wrapper.a"
c++ -std=c++17 -O3 \
    -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
    -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
    -o "$OUT/libraw_wrapper.o"

ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

echo "✔ rawkit build complete for $VERSION"