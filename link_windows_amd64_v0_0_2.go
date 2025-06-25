//go:build windows && amd64

package rawkit

/*
#cgo LDFLAGS: -L${SRCDIR}/libs/windows_amd64/v0.0.2 -lraw_wrapper -lraw -lz -lm -lstdc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
