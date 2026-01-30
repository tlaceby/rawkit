#include "rawkit.h"
#include "helpers.h"
#include <libraw/libraw_const.h>
#include <stdlib.h>
#include <string.h>

const char* get_libraw_version() {
    return libraw_version();
}

// ==================
// == INTERNAL API ==
// ==================

static char* safe_strdup(const char* src) {
    if (src == NULL || src[0] == '\0') {
        return NULL;
    }
    return strdup(src);
}

static libraw_data_t* open_raw_file(const char* file_path, int* err_out) {
    if (file_path == NULL || file_path[0] == '\0') {
        *err_out = LIBRAW_UNSPECIFIED_ERROR;
        return NULL;
    }

    if (!file_exists(file_path)) {
        *err_out = LIBRAW_IO_ERROR;
        return NULL;
    }

    libraw_data_t* handle = libraw_init(0);
    if (handle == NULL) {
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
        return NULL;
    }

    int res = libraw_open_file(handle, file_path);
    if (res != LIBRAW_SUCCESS) {
        libraw_close(handle);
        *err_out = res;
        return NULL;
    }

    return handle;
}

static ImageMeta* extract_meta(libraw_data_t* raw) {
    if (raw == NULL) {
        return NULL;
    }

    ImageMeta* meta = calloc(1, sizeof(ImageMeta));
    if (meta == NULL) {
        return NULL;
    }

    meta->iso = (int)raw->other.iso_speed;
    meta->aperture = raw->other.aperture;
    meta->shutter_speed = raw->other.shutter;
    meta->focal_length = raw->other.focal_len;
    meta->lens_model = safe_strdup(raw->lens.Lens);
    meta->camera_make = safe_strdup(raw->idata.make);
    meta->camera_model = safe_strdup(raw->idata.model);

    return meta;
}

static ImageData* process_raw(libraw_data_t* handle, int* err_out) {
    int res = libraw_unpack(handle);
    if (res != LIBRAW_SUCCESS) {
        *err_out = res;
        return NULL;
    }

    res = libraw_dcraw_process(handle);
    if (res != LIBRAW_SUCCESS) {
        *err_out = res;
        return NULL;
    }

    libraw_processed_image_t* img = libraw_dcraw_make_mem_image(handle, &res);
    if (res != LIBRAW_SUCCESS || img == NULL) {
        *err_out = (res != LIBRAW_SUCCESS) ? res : LIBRAW_UNSPECIFIED_ERROR;
        return NULL;
    }

    ImageData* data = calloc(1, sizeof(ImageData));
    if (data == NULL) {
        libraw_dcraw_clear_mem(img);
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
        return NULL;
    }

    data->width = img->width;
    data->height = img->height;
    data->channels = img->colors;
    data->colorspace = handle->params.output_color;

    size_t pixel_count = (size_t)img->width * img->height * img->colors;
    data->data = malloc(pixel_count * sizeof(uint16_t));
    
    if (data->data == NULL) {
        free(data);
        libraw_dcraw_clear_mem(img);
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
        return NULL;
    }

    if (img->bits == 16) {
        memcpy(data->data, img->data, pixel_count * sizeof(uint16_t));
        data->bit_depth = 16;
    } else {
        // Convert 8-bit to 16-bit
        for (size_t i = 0; i < pixel_count; i++) {
            data->data[i] = ((uint16_t)img->data[i]) << 8;
        }
        data->bit_depth = 16;
    }

    data->data_size = pixel_count * sizeof(uint16_t);
    libraw_dcraw_clear_mem(img);
    
    return data;
}

// ================
// == PUBLIC API ==
// ================

ImageMeta* read_raw_meta(const char* file_path, int* err_out) {
    *err_out = LIBRAW_SUCCESS;

    libraw_data_t* handle = open_raw_file(file_path, err_out);
    if (handle == NULL) {
        return NULL;
    }

    ImageMeta* meta = extract_meta(handle);
    libraw_close(handle);

    if (meta == NULL) {
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
    }

    return meta;
}

ImageData* read_raw_image(const char* file_path, int* err_out) {
    *err_out = LIBRAW_SUCCESS;

    libraw_data_t* handle = open_raw_file(file_path, err_out);
    if (handle == NULL) {
        return NULL;
    }

    ImageData* data = process_raw(handle, err_out);
    libraw_close(handle);

    return data;
}

ImageData* read_raw_full(const char* file_path, ImageMeta** meta_out, int* err_out) {
    *err_out = LIBRAW_SUCCESS;
    *meta_out = NULL;

    libraw_data_t* handle = open_raw_file(file_path, err_out);
    if (handle == NULL) {
        return NULL;
    }

    // Extract meta first (before processing modifies state)
    *meta_out = extract_meta(handle);

    ImageData* data = process_raw(handle, err_out);
    libraw_close(handle);

    // If image processing failed, clean up meta too
    if (data == NULL && *meta_out != NULL) {
        image_meta_free(*meta_out);
        *meta_out = NULL;
    }

    return data;
}

void image_data_free(ImageData* data) {
    if (data == NULL) {
        return;
    }
    free(data->data);
    free(data);
}

void image_meta_free(ImageMeta* meta) {
    if (meta == NULL) {
        return;
    }
    free(meta->lens_model);
    free(meta->camera_make);
    free(meta->camera_model);
    free(meta);
}