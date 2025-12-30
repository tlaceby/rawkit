#include "rawkit.h"
#include "helpers.h"

bool file_exists (const char* path) {
    FILE* fp;
    if ((fp = fopen(path, 'r'))) {
        fclose(fp);
        return true;
    }

    return false;
}

const char* get_libraw_version() {
    return libraw_version();
}