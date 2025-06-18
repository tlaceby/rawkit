#include "libraw_wrapper.h"
#include "../include/libraw/libraw.h"
#include <cstddef>
#include <cstdlib>

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
    return NULL;
}

