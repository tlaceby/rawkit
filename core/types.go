package core

type ImageType int

const (
	IMG_TYPE_RAW ImageType = iota
	IMG_TYPE_JPG
	IMG_TYPE_PNG
	IMG_TYPE_UNKNOWN
)

type RAWImageType int

const (
	RAW_TYPE_ARW RAWImageType = iota // Sony
	RAW_TYPE_CR2                     // Canon
	RAW_TYPE_CR3                     // Canon
	RAW_TYPE_NEF                     // Nikon
	RAW_TYPE_DNG                     // Adobe/Universal
	RAW_TYPE_ORF                     // Olympus
	RAW_TYPE_RAF                     // Fujifilm
	RAW_TYPE_RW2                     // Panasonic/Lumix
	RAW_TYPE_UNKNOWN
)

type Colorspace int

const (
	LIBRAW_COLORSPACE_NotFound Colorspace = iota
	LIBRAW_COLORSPACE_sRGB
	LIBRAW_COLORSPACE_AdobeRGB
	LIBRAW_COLORSPACE_WideGamutRGB
	LIBRAW_COLORSPACE_ProPhotoRGB
	LIBRAW_COLORSPACE_ICC
	LIBRAW_COLORSPACE_Uncalibrated // Tag 0x0001 InteropIndex containing "R03" + LIBRAW_COLORSPACE_Uncalibrated = Adobe RGB
	LIBRAW_COLORSPACE_CameraLinearUniWB
	LIBRAW_COLORSPACE_CameraLinear
	LIBRAW_COLORSPACE_CameraGammaUniWB
	LIBRAW_COLORSPACE_CameraGamma
	LIBRAW_COLORSPACE_MonochromeLinear
	LIBRAW_COLORSPACE_MonochromeGamma
	LIBRAW_COLORSPACE_Unknown Colorspace = 255
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
	Colorspace Colorspace
	Channels   int // 3=RGB, 4=RGBA
	BitDepth   int // original source: 8 (JPEG/PNG), [12, 14, 16 (RAW)]
	Data       []uint16
}

// Relevant metadata for RAW images
type ImageMeta struct {
	ISO         int     // eg: 4000
	Aperture    float32 // f2.8
	SS          float32 // 1/8s (0.125)
	FocalLength float32

	LensModel   string
	CameraMake  string // Sony
	CameraModel string // A6700
}

// ---------------------
// -- IMAGE PROCESSOR --
// ---------------------

type ImageProcessor struct {
	img      *Image // underlying image data
	filePath string
	raw      bool  // whether the image is RAW
	err      error // encountered an error during operation
	handle   *int  // nil when closed/unopened/not raw
}
