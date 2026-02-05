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

// Colorspace represents the color profile of image data.
// Values correspond to LibRaw's colorspace enum.
type Colorspace int

const (
	LIBRAW_COLORSPACE_NotFound          Colorspace = iota // Colorspace not found
	LIBRAW_COLORSPACE_sRGB                                // Standard RGB
	LIBRAW_COLORSPACE_AdobeRGB                            // Adobe RGB (1998)
	LIBRAW_COLORSPACE_WideGamutRGB                        // Wide Gamut RGB
	LIBRAW_COLORSPACE_ProPhotoRGB                         // ProPhoto RGB
	LIBRAW_COLORSPACE_ICC                                 // ICC Profile
	LIBRAW_COLORSPACE_Uncalibrated                        // Uncalibrated
	LIBRAW_COLORSPACE_CameraLinearUniWB                   // Camera Linear with UniWB
	LIBRAW_COLORSPACE_CameraLinear                        // Camera Linear
	LIBRAW_COLORSPACE_CameraGammaUniWB                    // Camera Gamma with UniWB
	LIBRAW_COLORSPACE_CameraGamma                         // Camera Gamma
	LIBRAW_COLORSPACE_MonochromeLinear                    // Monochrome Linear
	LIBRAW_COLORSPACE_MonochromeGamma                     // Monochrome Gamma
	LIBRAW_COLORSPACE_Unknown           Colorspace = 255  // Unknown colorspace
)

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
	Width      int        // Image width in pixels
	Height     int        // Image height in pixels
	Colorspace Colorspace // Color profile
	Channels   Channels   // Number of channels (3 for RGB)
	BitDepth   int        // Bits per channel (8 or 16)
	Data       []uint16   // Raw pixel data
}

// EXIFData contains EXIF metadata extracted from RAW files.
type EXIFData struct {
	ISO          int     // ISO sensitivity (e.g., 100, 400, 6400)
	Aperture     float32 // F-number (e.g., 2.8, 5.6)
	ShutterSpeed float32 // Shutter speed in seconds (e.g., 0.001 for 1/1000s)
	FocalLength  float32 // Focal length in millimeters
	LensModel    string  // Lens name/model
	CameraMake   string  // Camera manufacturer (e.g., "Sony")
	CameraModel  string  // Camera model (e.g., "A6700")
}
