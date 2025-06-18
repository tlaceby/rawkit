#pragma once
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif

typedef struct RawKitImage {
    uint16_t *buffer;   // stored as packed pixels (16-bit)
    int       width;
    int       height;
    int       channels; // 1 = CFA/raw, 3 = RGB (default is 3)
} RawKitImage;

/**
 * rawkit_load()
 * ------------------------------------
 * path   : full path to a RAW file.
 * returns: pointer to heap-allocated RawKitImage,
 *          or NULL if LibRaw throws/returns an error.
 *
 * The buffer and the RawKitImage struct are allocated together;
 * call rawkit_free(img) when done.
 */
RawKitImage *rawkit_load(const char *path);

/**
 * Frees both the pixel buffer and the RawKitImage container.
 */
void rawkit_free(RawKitImage *img);

/**
  * Returns integer representation of LibRaw version. During LibRaw development, the version number is always increase .
*/
int rawkit_libraw_version_num();

/**
* Returns string representation of LibRaw version in MAJOR.MINOR.PATCH-Status format (i.e. 0.6.0-Alpha2 or 0.6.1-Release).
*/
const char* rawkit_libraw_version_str();

#ifdef __cplusplus
} // extern "C"
#endif
