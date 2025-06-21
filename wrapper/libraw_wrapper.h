#pragma once
#include <stdint.h>
#include <sys/types.h>
#ifdef __cplusplus
extern "C" {
#endif

typedef struct RawKitImage {
  const char* format;
  uint16_t *buffer;   // stored as packed pixels (16-bit)
  ushort       width;       // Size of visible ("meaningful") part of the image (without the frame).
  ushort       height;

  int colors; // 1 = CFA/raw, 3 = RGB (default is 3)
  int flip; // Image orientation (0 if does not require rotation; 3 if requires 180-deg rotation; 5 if 90 deg counterclockwise, 6 if 90 deg clockwise). 

  int asShotWBApplied; // Set to 1 if WB already applied in camera (multishot modes; small raw)
  unsigned rawBitsPerPixel;
  unsigned rawCount;
  unsigned dngVersion;
  // const char* colorSpace;

  float iso;
  float shutterSpeed;
  float aperature;
  float focalLength;
  char artist[64];

  char cameraMake[64];
  char cameraModel[64];
  char cameraNormalizedMake[64];
  char cameraNormalizedModel[64];
  char cameraSoftware[64];

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
