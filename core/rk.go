// Package core provides RAW image reading and processing via LibRaw,
// with fallback support for standard image formats (JPEG, PNG).
package core

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/lib
// apple arm
#cgo darwin,arm64 CFLAGS: -I/opt/homebrew/opt/libraw/include
#cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/opt/libraw/lib -lraw
// apple intel
#cgo darwin,amd64 CFLAGS: -I/usr/local/opt/libraw/include
#cgo darwin,amd64 LDFLAGS: -L/usr/local/opt/libraw/lib -lraw

// linux all
#cgo linux pkg-config: libraw

// windows all
#cgo windows CFLAGS: -I${LIBRAW_PATH}/include
#cgo windows LDFLAGS: -L${LIBRAW_PATH}/lib -lraw

#include "rawkit.h"
#include "helpers.c"
#include "processor.c"
#include <stdlib.h>
*/
import "C"

// Returns type information based on image file extension
func DetectImageType(ext string) (ImageType, RAWImageType) {
	switch strings.ToLower(ext) {
	case ".arw":
		return IMG_TYPE_RAW, RAW_TYPE_ARW
	case ".cr2":
		return IMG_TYPE_RAW, RAW_TYPE_CR2
	case ".cr3":
		return IMG_TYPE_RAW, RAW_TYPE_CR3
	case ".nef":
		return IMG_TYPE_RAW, RAW_TYPE_NEF
	case ".dng":
		return IMG_TYPE_RAW, RAW_TYPE_DNG
	case ".orf":
		return IMG_TYPE_RAW, RAW_TYPE_ORF
	case ".raf":
		return IMG_TYPE_RAW, RAW_TYPE_RAF
	case ".rw2":
		return IMG_TYPE_RAW, RAW_TYPE_RW2
	case ".jpg", ".jpeg":
		return IMG_TYPE_JPG, RAW_TYPE_UNKNOWN
	case ".png":
		return IMG_TYPE_PNG, RAW_TYPE_UNKNOWN
	default:
		return IMG_TYPE_UNKNOWN, RAW_TYPE_UNKNOWN
	}
}

// LibrawVersion returns the version string of the underlying LibRaw library.
func LibrawVersion() string {
	return C.GoString(C.get_libraw_version())
}

// ReadRAWMetadata extracts metadata from a RAW image file without processing pixels.
// This is significantly faster than [ReadAll] when only metadata is needed.
// Returns an error if the file is not a supported RAW format.
//
// Example:
//
//	meta, err := core.ReadRAWMetadata("/path/to/photo.arw")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Camera: %s %s\n", meta.CameraMake, meta.CameraModel)
//	fmt.Printf("Settings: ISO %d, f/%.1f, %gs\n", meta.ISO, meta.Aperture, meta.SS)
func ReadRAWMetadata(filePath string) (*EXIFData, error) {
	if !isRawFile(filePath) {
		return nil, errors.New("file is not raw")
	}

	var errCode C.int
	path := C.CString(filePath)
	defer C.free(unsafe.Pointer(path))

	cMeta := C.read_raw_meta(path, &errCode)
	if LibrawError(errCode) != LibrawSuccess {
		return nil, LibrawError(errCode)
	}

	if cMeta == nil {
		return nil, nil
	}

	meta := &EXIFData{
		ISO:          int(cMeta.iso),
		Aperture:     float32(cMeta.aperture),
		ShutterSpeed: float32(cMeta.shutter_speed),
		FocalLength:  float32(cMeta.focal_length),
		LensModel:    C.GoString(cMeta.lens_model),
		CameraMake:   C.GoString(cMeta.camera_make),
		CameraModel:  C.GoString(cMeta.camera_model),
	}

	C.image_meta_free(cMeta)
	return meta, nil
}

// ReadImageData reads pixel data from an image file.
// Supports both RAW formats (via LibRaw) and standard formats (JPEG, PNG).
// Use this when metadata is not needed.
//
// Example:
//
//	data, err := core.ReadImageData("/path/to/photo.arw")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Dimensions: %dx%d\n", data.Width, data.Height)
//	r, g, b := data.Pixel(100, 200)
//	fmt.Printf("Pixel at (100,200): R=%d G=%d B=%d\n", r, g, b)
func ReadImageData(filePath string) (*ImageData, error) {
	if isRawFile(filePath) {
		return readRawImageData(filePath)
	}

	return readStandardImage(filePath)
}

// ReadAll loads both pixel data and metadata from an image file.
// For RAW files, both Data and Meta are populated.
// For non-RAW files (JPEG, PNG), Meta will be nil.
//
// Example:
//
//	img, err := core.ReadAll("/path/to/photo.arw")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Dimensions: %dx%d\n", img.Data.Width, img.Data.Height)
//	if img.Meta != nil {
//	    fmt.Printf("Camera: %s\n", img.Meta.CameraModel)
//	}
func ReadAll(filePath string) (*Image, error) {
	imgType, rawType := DetectImageType(filepath.Ext(filePath))

	img := &Image{
		Type:    imgType,
		RawType: rawType,
		Path:    filePath,
	}

	if isRawFile(filePath) {
		data, meta, err := readRawFull(filePath)
		if err != nil {
			return nil, err
		}
		img.Data = data
		img.Meta = meta
		return img, nil
	}

	data, err := readStandardImage(filePath)
	if err != nil {
		return nil, err
	}

	img.Data = data
	return img, nil
}

// ======================
// == INTERNAL HELPERS ==
// ======================

func readRawImageData(filePath string) (*ImageData, error) {
	var errCode C.int
	path := C.CString(filePath)
	defer C.free(unsafe.Pointer(path))

	cData := C.read_raw_image(path, &errCode)
	if LibrawError(errCode) != LibrawSuccess {
		return nil, LibrawError(errCode)
	}

	imgData := convertCImageData(cData)
	if !imgData.Colorspace.Supported() {
		return nil, fmt.Errorf("unsupported colorspace %d", imgData.Colorspace)
	}

	C.image_data_free(cData)
	return imgData, nil
}

func readRawFull(filePath string) (*ImageData, *EXIFData, error) {
	var errCode C.int
	var cMeta *C.ImageMeta
	path := C.CString(filePath)
	defer C.free(unsafe.Pointer(path))

	cData := C.read_raw_full(path, &cMeta, &errCode)
	if LibrawError(errCode) != LibrawSuccess {
		return nil, nil, LibrawError(errCode)
	}

	imgData := convertCImageData(cData)
	C.image_data_free(cData)

	var meta *EXIFData
	if cMeta != nil {
		meta = &EXIFData{
			ISO:          int(cMeta.iso),
			Aperture:     float32(cMeta.aperture),
			ShutterSpeed: float32(cMeta.shutter_speed),
			FocalLength:  float32(cMeta.focal_length),
			LensModel:    C.GoString(cMeta.lens_model),
			CameraMake:   C.GoString(cMeta.camera_make),
			CameraModel:  C.GoString(cMeta.camera_model),
		}
		C.image_meta_free(cMeta)
	}

	return imgData, meta, nil
}

func convertCImageData(cData *C.ImageData) *ImageData {
	if cData == nil {
		return nil
	}

	width := int(cData.width)
	height := int(cData.height)
	channels := int(cData.channels)
	pixelCount := width * height * channels

	data := make([]uint16, pixelCount)
	cSlice := (*[1 << 30]C.uint16_t)(unsafe.Pointer(cData.data))[:pixelCount:pixelCount]
	for i := range data {
		data[i] = uint16(cSlice[i])
	}

	return &ImageData{
		Width:      width,
		Height:     height,
		Colorspace: Colorspace(cData.colorspace),
		Channels:   Channels(channels),
		Data:       data,
	}
}

func readStandardImage(filePath string) (*ImageData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	data := make([]uint16, width*height*3)

	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[idx] = uint16(r)
			data[idx+1] = uint16(g)
			data[idx+2] = uint16(b)
			idx += 3
		}
	}

	return &ImageData{
		Width:      width,
		Height:     height,
		Colorspace: LIBRAW_COLORSPACE_sRGB,
		Channels:   LIBRAW_CHANNELS_RGB,
		Data:       data,
	}, nil
}

func isRawFile(filePath string) bool {
	imgType, _ := DetectImageType(filepath.Ext(filePath))
	return imgType == IMG_TYPE_RAW
}
