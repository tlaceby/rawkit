#include "rawkit.h"
#include "memory.h"
#include "string.h"
#include "helpers.h"

ImageProcessor* image_processor_create(const char* file_path, int* err_out) {
    // ensure string is not empty
    if (file_path == NULL || strcmp(file_path, "") == 0) {
        *err_out = LIBRAW_UNSPECIFIED_ERROR;
        return NULL;
    }

    if (file_exists(file_path) == false) {
        *err_out = LIBRAW_FILE_UNSUPPORTED;
        return NULL;
    }

    libraw_data_t* handle = malloc(sizeof(libraw_data_t));
    if (handle == NULL) {
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
        return NULL;
    }

    int res = libraw_open_file(handle, file_path);
    if (res != LIBRAW_SUCCESS) {
        free(handle);
        *err_out = res;
        return NULL;
    }

    ImageProcessor* processor = malloc(sizeof(ImageProcessor));
    if (processor == NULL) {
        free(handle);
        *err_out = LIBRAW_UNSUFFICIENT_MEMORY;
        return NULL;
    }

    processor->handle = handle;
    processor->file_path = file_path;
    return processor;
}

void image_processor_cleanup(ImageProcessor* processor) {
    if (processor == NULL) {
        return;
    }

    if (processor->handle != NULL) {
        libraw_close(processor->handle);
        free(processor->file_path);
    }
}