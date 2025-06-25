//go:build linux && arm64

package rawkit

/*
#cgo LDFLAGS: -L${SRCDIR}/libs/linux_arm64/v0.0.5 -lraw_wrapper -lraw -lz -lm -lstdc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
