package rk

type ImageType string
type RAWImageType string

const (
	IMG_TYPE_RAW     ImageType = "RAW"
	IMG_TYPE_JPG     ImageType = "JPG"
	IMG_TYPE_PNG     ImageType = "PNG"
	IMG_TYPE_UNKNOWN ImageType = "UNKNOWN"
)

const (
	RAW_TYPE_ARW     RAWImageType = "ARW" // Sony
	RAW_TYPE_CR2     RAWImageType = "CR2" // Canon
	RAW_TYPE_CR3     RAWImageType = "CR3" // Canon
	RAW_TYPE_NEF     RAWImageType = "NEF" // Nikon
	RAW_TYPE_DNG     RAWImageType = "DNG" // Adobe/Universal
	RAW_TYPE_ORF     RAWImageType = "ORF" // Olympus
	RAW_TYPE_RAF     RAWImageType = "RAF" // Fujifilm
	RAW_TYPE_RW2     RAWImageType = "RW2" // Panasonic/Lumix
	RAW_TYPE_UNKNOWN RAWImageType = "UNKNOWN"
)

// Represents a loaded Image from the filesystem
type Image struct {
	Type      ImageType
	RawType   RAWImageType
	Path      string
	Thumbnail *ImageData
	Data      *ImageData
	Meta      *ImageMeta // only RAW images will contain this info (otherwise will be nil)
}

// Stores processed image data
type ImageData struct {
	Width      int
	Height     int
	Colorspace string // "sRGB", "AdobeRGB"
	Channels   int    // 3=RGB, 4=RGBA
	BitDepth   int    // original source: 8 (JPEG/PNG), [12, 14, 16 (RAW)]
	Data       []uint16
}

// Relevant metadata for RAW images
type ImageMeta struct {
	Aperture    float32 // f2.8
	ISO         int     // eg: 4000
	SS          float32 // 1/8s (0.125)
	FocalLength float32

	LensModel   string
	CameraMake  string // Sony
	CameraModel string // A6700
}
