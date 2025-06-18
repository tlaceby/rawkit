package rawkit

/*
#cgo LDFLAGS: -lraw_wrapper -lraw -lz -lm
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func LoadRAW(path string) (RawKitImage, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	out := C.rawkit_load(cpath)
	if out == nil {
		return RawKitImage{}, fmt.Errorf("rawkit: failed to load %s", path)
	}

	defer C.rawkit_free(out)

	// size in pixels * channels
	size := int(out.width * out.height * out.channels)
	data := unsafe.Slice((*uint16)(unsafe.Pointer(out.buffer)), size)

	rawkit_img := RawKitImage{
		buffer:   make([]uint16, size),
		width:    int(out.width),
		height:   int(out.height),
		channels: int(out.channels),
	}

	copy(rawkit_img.buffer, data)

	return rawkit_img, nil
}
