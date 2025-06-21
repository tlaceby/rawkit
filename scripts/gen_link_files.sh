#!/usr/bin/env bash
set -euo pipefail

# ------------------------------------------------------------
# Generate version-specific link_*.go files for CGO
# Usage: ./scripts/gen_link_files.sh v1.3.6
# ------------------------------------------------------------

VERSION="${1:?Missing version argument, e.g. v1.3.6}"
TAG_SAFE=$(echo "$VERSION" | tr '.' '_' | tr -d 'v')
OUT_DIR="."

platforms=(
  "darwin arm64 -lc++"
  "darwin amd64 -lc++"
  "linux amd64 -lstdc++"
  "linux arm64 -lstdc++"
  "windows amd64 -lstdc++"
)

echo "▶ Generating CGO link files for version $VERSION"

for p in "${platforms[@]}"; do
  set -- $p
  os=$1; arch=$2; cpp=$3
  file="${OUT_DIR}/link_${os}_${arch}_v${TAG_SAFE}.go"

  cat <<EOF > "$file"
//go:build $os && $arch

package rawkit

/*
#cgo LDFLAGS: -L\${SRCDIR}/libs/${os}_${arch}/${VERSION} -lraw_wrapper -lraw -lz -lm ${cpp}
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
EOF

  echo "• Created $file"
done

echo "✔ link_*.go files generated"
