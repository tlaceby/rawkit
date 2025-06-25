//go:build darwin && amd64

package rawkit

/*
#cgo LDFLAGS: -L${SRCDIR}/libs/darwin_amd64/0.0.5 -lraw_wrapper -lraw -lz -lm -lc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
