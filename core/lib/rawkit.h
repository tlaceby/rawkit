#ifndef _rk_core_
#define _rk_core_

#include <libraw/libraw.h>
#include <stdint.h>
#include <stdbool.h>

typedef struct {
    int width;
    int height;
    int colorspace;
    int channels;
    int bit_depth;
    uint16_t* data;
    size_t data_size;
} ImageData;

typedef struct {
    int iso;
    float aperture;
    float shutter_speed;
    float focal_length;
    char* lens_model;
    char* camera_model;
    char* camera_make;
} ImageMeta;

// Returns the current libraw version
const char* get_libraw_version();

// Read metadata only (fast)
ImageMeta* read_raw_meta(const char* file_path, int* err_out);

// Read image data only
ImageData* read_raw_image(const char* file_path, int* err_out);

// Read both metadata and image data
ImageData* read_raw_full(const char* file_path, ImageMeta** meta_out, int* err_out);

// Memory cleanup
void image_data_free(ImageData* data);
void image_meta_free(ImageMeta* meta);

#endif