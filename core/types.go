// Package core provides RAW image reading and processing via LibRaw,
// with fallback support for standard image formats (JPEG, PNG).
package core

// ImageType represents the format category of an image file.
type ImageType int

const (
	IMG_TYPE_RAW     ImageType = iota // RAW camera format
	IMG_TYPE_JPG                      // JPEG format
	IMG_TYPE_PNG                      // PNG format
	IMG_TYPE_UNKNOWN                  // Unknown or unsupported format
)

// String returns the string representation of ImageType.
func (t ImageType) String() string {
	switch t {
	case IMG_TYPE_RAW:
		return "RAW"
	case IMG_TYPE_JPG:
		return "JPEG"
	case IMG_TYPE_PNG:
		return "PNG"
	case IMG_TYPE_UNKNOWN:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// RAWImageType represents specific RAW format variants by manufacturer.
type RAWImageType int

const (
	RAW_TYPE_ARW     RAWImageType = iota // Sony
	RAW_TYPE_CR2                         // Canon
	RAW_TYPE_CR3                         // Canon
	RAW_TYPE_NEF                         // Nikon
	RAW_TYPE_DNG                         // Adobe/Universal
	RAW_TYPE_ORF                         // Olympus
	RAW_TYPE_RAF                         // Fujifilm
	RAW_TYPE_RW2                         // Panasonic/Lumix
	RAW_TYPE_UNKNOWN                     // Unknown RAW format
)

// String returns the string representation of RAWImageType.
func (t RAWImageType) String() string {
	switch t {
	case RAW_TYPE_ARW:
		return "ARW"
	case RAW_TYPE_CR2:
		return "CR2"
	case RAW_TYPE_CR3:
		return "CR3"
	case RAW_TYPE_NEF:
		return "NEF"
	case RAW_TYPE_DNG:
		return "DNG"
	case RAW_TYPE_ORF:
		return "ORF"
	case RAW_TYPE_RAF:
		return "RAF"
	case RAW_TYPE_RW2:
		return "RW2"
	case RAW_TYPE_UNKNOWN:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// represents the color profile of image data  --  Values correspond to LibRaw's colorspace enum.
type Colorspace int

const (
	LIBRAW_COLORSPACE_NotFound     Colorspace = iota // Colorspace not found
	LIBRAW_COLORSPACE_sRGB                           // Standard RGB
	LIBRAW_COLORSPACE_AdobeRGB                       // Adobe RGB (1998)
	LIBRAW_COLORSPACE_WideGamutRGB                   // Wide Gamut RGB
	LIBRAW_COLORSPACE_ProPhotoRGB                    // ProPhoto RGB
	LIBRAW_COLORSPACE_Unknown      Colorspace = 255  // Unknown colorspace
)

// Valid reports whether this colorspace is supported
func (c Colorspace) Supported() bool {
	switch c {
	case LIBRAW_COLORSPACE_NotFound,
		LIBRAW_COLORSPACE_sRGB,
		LIBRAW_COLORSPACE_AdobeRGB,
		LIBRAW_COLORSPACE_WideGamutRGB,
		LIBRAW_COLORSPACE_ProPhotoRGB:
		return true
	default:
		return false
	}
}

// Either RGBA OR RGB (4, 3)
type Channels int

const (
	LIBRAW_CHANNELS_RGB  Channels = 3
	LIBRAW_CHANNELS_RGBA Channels = 4
)

// Image represents a loaded image from the filesystem.
// Meta is only populated for RAW files; it will be nil for JPEG/PNG.
type Image struct {
	Type    ImageType    // Format category (RAW, JPG, PNG, etc.)
	RawType RAWImageType // Specific RAW format (ARW, CR2, etc.)
	Path    string       // Original file path
	Data    *ImageData   // Processed pixel data
	Meta    *EXIFData    // Camera metadata (nil for non-RAW)
}

// ImageData stores processed pixel data in 16-bit RGB format.
// Data is laid out as [R0, G0, B0, R1, G1, B1, ...] in row-major order.
type ImageData struct {
	Width      int        // heidth in pixels
	Height     int        // height in pixels
	Colorspace Colorspace // Color profile
	Channels   Channels   // Number of channels (3 for RGB)
	Data       []uint16   // Raw pixel data stored with a 16 bit-depth
}

// contains EXIF metadata extracted from RAW files.
type EXIFData struct {
	ISO          int     // 400
	Aperture     float32 // 2.8
	ShutterSpeed float32 // 0.001 for 1/1000s
	FocalLength  float32 // 35.0mm

	LensModel   string // "Sigma 16mm f1.4"
	CameraMake  string // "Sony"
	CameraModel string // "A6700"
}
