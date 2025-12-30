#ifndef _rk_core_
#define _rk_core_

#include <libraw/libraw.h>
#include "stdint.h"

typedef enum ImageType {
	IMG_TYPE_RAW,
	IMG_TYPE_JPG,
	IMG_TYPE_PNG,
	IMG_TYPE_UNKNOWN,
} ImageType;

typedef enum RAWImageType {
	RAW_TYPE_ARW,        // Sony
	RAW_TYPE_CR2,        // Canon
	RAW_TYPE_CR3,        // Canon
	RAW_TYPE_NEF,        // Nikon
	RAW_TYPE_DNG,        // Adobe/Universal
	RAW_TYPE_ORF,        // Olympus
	RAW_TYPE_RAF,        // Fujifilm
	RAW_TYPE_RW2,        // Panasonic/Lumix
	RAW_TYPE_UNKNOWN,
} RAWImageType;

typedef struct {
    ImageType type;
	RAWImageType raw_type;
	char* path;
	ImageData* thumbnail;
	ImageData* data;
	ImageMeta* meta; // only RAW images will contain this info (otherwise will be nil)
} Image;

typedef struct {
    int width;
    int height;
    int colorspace; // Libraw_colorspace
    int channels; // 3=RGB, 4=RGBA
    int bit_depth; // original source: 8 (JPEG/PNG), [12, 14, 16 (RAW)]
    uint16_t* data; // (size= bitDepth * channels * width * height * sizeof(uint16_t))
} ImageData;

typedef struct {
    int iso;
    float aperture;
    float ss;
    float focal_length;

    char* lens_model;
    char* camera_model;
    char* camera_make;
} ImageMeta;


typedef struct {
    Image* img;
    const char* file_path;
    libraw_data_t* handle;
} ImageProcessor;

// ------------------------
// -- PUBLIC API MEMBERS --
// ------------------------

// Returns the current libraw version as a string
const char* get_libraw_version();
// Attempt to read and process a RAW image with Libraw
ImageProcessor* image_processor_create(const char* file_path, int* err_out);
// Perform memory cleanup on ImageProcessor
void image_processor_cleanup(ImageProcessor* processor);

#endif