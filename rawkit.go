package rawkit

//go:generate bash scripts/ensure_current.sh

/*
#cgo darwin,arm64  LDFLAGS: -L${SRCDIR}/libs/darwin_arm64/current  -lraw_wrapper -lraw -lz -lm -lc++
#cgo darwin,amd64  LDFLAGS: -L${SRCDIR}/libs/darwin_amd64/current  -lraw_wrapper -lraw -lz -lm -lc++
#cgo linux,amd64   LDFLAGS: -L${SRCDIR}/libs/linux_amd64/current   -lraw_wrapper -lraw -lz -lm -lstdc++
#cgo linux,arm64   LDFLAGS: -L${SRCDIR}/libs/linux_arm64/current   -lraw_wrapper -lraw -lz -lm -lstdc++
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs/windows_amd64/current -lraw_wrapper -lraw -lz -lstdc++
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

func LoadRAW(path string) (*RawKitImage, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	out := C.rawkit_load(cpath)
	if out == nil {
		return nil, fmt.Errorf("rawkit: failed to load %s", path)
	}

	pixels := int(out.width * out.height * out.channels)
	src := unsafe.Slice((*uint16)(unsafe.Pointer(out.buffer)), pixels)

	buf := make([]uint16, pixels)
	copy(buf, src)

	runtime.KeepAlive(out)
	defer C.rawkit_free(out)

	img := &RawKitImage{
		Buffer:   buf,
		Width:    int(out.width),
		Height:   int(out.height),
		Channels: int(out.channels),
	}

	return img, nil
}

func LibRawVersionNum() int {
	return int(C.rawkit_libraw_version_num())
}

func LibRawVersionStr() string {
	return C.GoString(C.rawkit_libraw_version_str())
}
