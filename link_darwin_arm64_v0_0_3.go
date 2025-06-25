//go:build darwin && arm64

package rawkit

/*
#cgo LDFLAGS: -L${SRCDIR}/libs/darwin_arm64/v0.0.3 -lraw_wrapper -lraw -lz -lm -lc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
