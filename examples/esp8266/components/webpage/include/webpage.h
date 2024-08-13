#pragma once

#include <stdlib.h>

typedef struct webfs {
    const char         *path;
    const char         *content;
    size_t              content_len;
    const struct webfs *next;
} webfs_t;

typedef struct webpage {
    const webfs_t *webfs;
} webpage_t;

webpage_t webpage_init(void);
void webpage_free(webpage_t *page);

const char *webpage_find(webpage_t *page, const char *path, size_t *content_len);
