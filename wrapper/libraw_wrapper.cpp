#include "libraw_wrapper.h"
#include "../include/libraw/libraw.h"
#include <cstddef>
#include <cstdint>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include "memory.h"

extern "C" int rawkit_libraw_version_num() {
    return libraw_versionNumber();
}

extern "C" const char* rawkit_libraw_version_str() {
    return libraw_version();
}

extern "C" void rawkit_free(RawKitImage *img) {
    if (!img) return;              
    free(img->buffer);             
    free(img);
}

extern "C" RawKitImage* rawkit_load(const char* filepath) {
    LibRaw lr;
    int status;

    status = lr.open_file(filepath);
    if (status != LIBRAW_SUCCESS) return NULL;

    status = lr.unpack();
    if (status != LIBRAW_SUCCESS) { 
        lr.recycle(); 
        return NULL; 
    }               

    // Number of RAW images in file (0 means that the file has not been recognized).
    if (lr.imgdata.idata.raw_count == 0) {
        lr.recycle();
        return NULL;
    }

    lr.adjust_sizes_info_only();    
    lr.imgdata.params.cropbox[0] = lr.imgdata.sizes.left_margin;
    lr.imgdata.params.cropbox[1] = lr.imgdata.sizes.top_margin; 
    lr.imgdata.params.cropbox[2] = lr.imgdata.sizes.width;         
    lr.imgdata.params.cropbox[3] = lr.imgdata.sizes.height;        
    lr.imgdata.params.user_flip  = lr.imgdata.sizes.flip;          
    lr.imgdata.params.output_bps = 16;

    if (lr.dcraw_process() != LIBRAW_SUCCESS) { lr.recycle(); return NULL; }
    libraw_processed_image_t* img = lr.dcraw_make_mem_image();
    if (!img) {
        lr.recycle();
        return NULL;
    }

    RawKitImage* rki = (RawKitImage*)malloc(sizeof(RawKitImage));
    if (!rki) {
        lr.dcraw_clear_mem(img);
        lr.recycle();
        return NULL;
    }

    rki->buffer = (uint16_t*)malloc(img->data_size);
    if (!rki->buffer) {
        free(rki);
        lr.dcraw_clear_mem(img);
        lr.recycle();
        return NULL;
    }

    libraw_data_t data = lr.imgdata;
    switch (img->type) {
        case LIBRAW_IMAGE_JPEG: rki->format = "img/jpeg"; break;
        case LIBRAW_IMAGE_JPEGXL: rki->format = "img/jpegxl"; break;
        case LIBRAW_IMAGE_BITMAP: rki->format = "img/bitmap"; break;
        case LIBRAW_IMAGE_H265: rki->format = "img/h265"; break;
        default: {
            free(rki);
            lr.dcraw_clear_mem(img);
            lr.recycle();
            return NULL;
        }
    }

    memcpy(rki->buffer, img->data, img->data_size);
    rki->width = img->width;
    rki->height = img->height;
    rki->colors = data.idata.colors;
    rki->flip = data.sizes.flip;

    rki->asShotWBApplied = data.color.as_shot_wb_applied;
    rki->rawBitsPerPixel = data.color.raw_bps;
    rki->rawCount = data.idata.raw_count;

    rki->iso = data.other.iso_speed;
    rki->aperature = data.other.aperture;
    rki->shutterSpeed = data.other.shutter;
    rki->focalLength = data.other.focal_len;

    memcpy(rki->artist, data.other.artist, 64 * sizeof(char));
    memcpy(rki->cameraMake, data.idata.normalized_make, 64 * sizeof(char));
    memcpy(rki->cameraModel, data.idata.model, 64 * sizeof(char));
    memcpy(rki->cameraNormalizedMake, data.idata.normalized_make, 64 * sizeof(char));
    memcpy(rki->cameraNormalizedModel, data.idata.normalized_model, 64 * sizeof(char));
    memcpy(rki->cameraSoftware, data.idata.software, 64 * sizeof(char));

    lr.dcraw_clear_mem(img);
    lr.recycle();

    return rki;
}
