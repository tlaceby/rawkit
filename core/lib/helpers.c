#include "helpers.h"
#include <stdio.h>

bool file_exists(const char* path) {
    if (path == NULL) {
        return false;
    }
    
    FILE* fp = fopen(path, "r");
    if (fp) {
        fclose(fp);
        return true;
    }

    return false;
}