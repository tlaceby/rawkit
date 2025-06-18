#!/usr/bin/env bash
set -euo pipefail
# ------------------------------------------------------------
# rawkit — build static LibRaw + wrapper for Linux targets
# Usage:  ./scripts/build_linux.sh [VERSION] [ARCH]
# Example: ./scripts/build_linux.sh v0.2.0 amd64
# ------------------------------------------------------------

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -f "${PROJECT_ROOT}/.env" ]] && { set -a; . "${PROJECT_ROOT}/.env"; set +a; }

VERSION="${1:-${VERSION:-v0.0.1}}"
ARCH="${2:-amd64}"               # amd64 | arm64

OUT="$PROJECT_ROOT/libs/linux_${ARCH}/${VERSION}"
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

echo "▶ Building rawkit for Linux ${ARCH} → ${OUT}"
mkdir -p "$OUT" "$INC_DIR"

# cross-compile tuple
case "$ARCH" in
  amd64) HOST="x86_64-linux-gnu"  ; CC=gcc  ; CXX=g++  ;;
  arm64) HOST="aarch64-linux-gnu"; CC=${HOST}-gcc ; CXX=${HOST}-g++ ;;
  *) echo "✖ unknown arch $ARCH"; exit 1 ;;
esac

pushd "$LIBRAW_DIR" >/dev/null
if [[ ! -f lib/libraw.a ]]; then
  echo "• Configuring LibRaw"
  CC=$CC CXX=$CXX RANLIB=${HOST}-ranlib \
  ./configure --host="$HOST" --disable-shared --enable-static CFLAGS="-O3"
  echo "• Compiling LibRaw"
  make -j"$(nproc)"
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
${CXX} -std=c++17 -O3 \
  -I"$INC_DIR" -I"$LIBRAW_DIR/libraw" \
  -c "$WRAPPER_DIR/libraw_wrapper.cpp" \
  -o "$OUT/libraw_wrapper.o"

${HOST}-ar rcs "$OUT/libraw_wrapper.a" "$OUT/libraw_wrapper.o"
rm "$OUT/libraw_wrapper.o"

ln -sfn "$OUT" "$PROJECT_ROOT/libs/linux_${ARCH}/current"
echo "✔ rawkit done  (current → ${VERSION})"
