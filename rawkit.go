package rawkit

import (
	"fmt"
	"runtime"
	"unsafe"
)

/*
#include <stdlib.h>
#include "wrapper/libraw_wrapper.h"
*/
import "C"

// LoadRAW opens and decodes a RAW image file at the given path using LibRaw.
// It returns a fully populated *RawKitImage struct containing pixel data and metadata.
// Returns an error if decoding fails or if the file is unsupported.
func LoadRAW(path string) (*RawKitImage, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	out := C.rawkit_load(cpath)
	if out == nil {
		return nil, fmt.Errorf("rawkit: failed to load %s", path)
	}

	pixels := int(out.width) * int(out.height) * int(out.colors)
	src := unsafe.Slice((*uint16)(unsafe.Pointer(out.buffer)), pixels)

	buf := make([]uint16, pixels)
	copy(buf, src)

	runtime.KeepAlive(out)

	img := &RawKitImage{
		Format: LibrawImageFormat(C.GoString(out.format)),
		Buffer: buf,
		Width:  int(out.width),
		Height: int(out.height),

		Colors:          int(out.colors),
		Flip:            int(out.flip),
		AsShotWBApplied: int(out.asShotWBApplied),
		RawBitsPerPixel: uint(out.rawBitsPerPixel),
		RawCount:        uint(out.rawCount),
		DNGVersion:      uint(out.dngVersion),
		ISO:             float32(out.iso),
		ShutterSpeed:    float32(out.shutterSpeed),
		Aperature:       float32(out.aperature),
		FocalLength:     float32(out.focalLength),
		Artist:          C.GoStringN(&out.artist[0], 64),

		CameraMake:            C.GoStringN(&out.cameraMake[0], 64),
		CameraModel:           C.GoStringN(&out.cameraModel[0], 64),
		CameraNormalizedMake:  C.GoStringN(&out.cameraNormalizedMake[0], 64),
		CameraNormalizedModel: C.GoStringN(&out.cameraNormalizedModel[0], 64),
		CameraSoftware:        C.GoStringN(&out.cameraSoftware[0], 64),
	}

	C.rawkit_free(out)
	return img, nil
}

// LibRawVersionNum returns the numeric version of the linked LibRaw library.
func LibRawVersionNum() int {
	return int(C.rawkit_libraw_version_num())
}

// LibRawVersionStr returns the human-readable version string of the linked LibRaw library.
func LibRawVersionStr() string {
	return C.GoString(C.rawkit_libraw_version_str())
}
