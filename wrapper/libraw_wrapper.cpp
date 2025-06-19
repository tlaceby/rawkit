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

    memcpy(rki->buffer, img->data, img->data_size);
    rki->width = img->width;
    rki->height = img->height;
    rki->channels = img->colors;

    lr.dcraw_clear_mem(img);
    lr.recycle();

    return rki;
}
