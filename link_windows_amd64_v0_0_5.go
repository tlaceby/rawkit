//go:build windows && amd64

package rawkit

/*
#cgo LDFLAGS: -L${SRCDIR}/libs/windows_amd64/0.0.5 -lraw_wrapper -lraw -lz -lm -lstdc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
