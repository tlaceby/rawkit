#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for macOS
# Usage:  ./scripts/build_darwin.sh [VERSION] [ARCH]
# ------------------------------------------------------------

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -f "${PROJECT_ROOT}/.env" ]] && { set -a; . "${PROJECT_ROOT}/.env"; set +a; }

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH="${2:-$(uname -m)}"            # arm64 | x86_64
OUT="$PROJECT_ROOT/libs/darwin_${ARCH}/${VERSION}"

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

echo "▶ Building rawkit for macOS ${ARCH} → ${OUT}"
mkdir -p "$OUT" "$INC_DIR"


pushd "$LIBRAW_DIR" >/dev/null
if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  ./configure --disable-shared --enable-static CFLAGS="-O3" >/dev/null
  echo "• Compiling LibRaw"
  make -j"$(sysctl -n hw.ncpu)"        >/dev/null
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
c++ -std=c++17 -O3 \
    -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
    -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
    -o "$OUT/libraw_wrapper.o"

ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

ln -sfn "$OUT" "$PROJECT_ROOT/libs/darwin_${ARCH}
