package rk

/*
#cgo CFLAGS: -I${SRCDIR}/core

// Apple Silicon OSX (Install w/ Homebrew)
#cgo darwin,arm64 CFLAGS: -I/opt/homebrew/opt/libraw/include
#cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/opt/libraw/lib -lraw

// Intel OSX (Install w/ Homebrew)
#cgo darwin,amd64 CFLAGS: -I/usr/local/opt/libraw/include
#cgo darwin,amd64 LDFLAGS: -L/usr/local/opt/libraw/lib -lraw

// Linux
#cgo linux pkg-config: libraw

// Window expects LIBRAW_PATH env or vcpkg default
#cgo windows CFLAGS: -I${LIBRAW_PATH}/include
#cgo windows LDFLAGS: -L${LIBRAW_PATH}/lib -lraw

#include "helpers.h"
#include "helpers.c"
*/
import "C"

// Get the version of Libraw currently in use
func Version() string {
	return C.GoString(C.get_libraw_version())
}
