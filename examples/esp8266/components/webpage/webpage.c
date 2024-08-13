#include "include/webpage.h"

#include <stdint.h>
#include <stdbool.h>
#include <string.h>

#include <esp_log.h>

extern const uint8_t webfs_start[] asm("_binary_webfs_bin_start");
extern const uint8_t webfs_end[] asm("_binary_webfs_bin_end");

static bool path_cmp(const char *path_a, const char *path_b)
{
    while (*path_a == '/') ++path_a;
    while (*path_b == '/') ++path_b;

    while (*path_a != 0 && *path_b != 0)
    {
        if (*path_a != *path_b) return false;
        ++path_a;
        ++path_b;
    }

    return *path_a == 0 && *path_b == 0;
}

webpage_t webpage_init(void)
{
    webfs_t *webfs = malloc(sizeof(webfs_t));
    memset(webfs, 0, sizeof(webfs_t));

    webpage_t page = {.webfs = webfs};
    uint32_t path_len;

    const uint8_t *mem = &webfs_start[0];
    while(true)
    {
        memcpy(&path_len, mem, sizeof(uint32_t));
        mem += sizeof(uint32_t);
        webfs->path = (const char *)mem;
        mem += path_len;

        memcpy(&webfs->content_len, mem, sizeof(uint32_t));
        mem += sizeof(uint32_t);
        webfs->content = (const char *)mem;
        mem += webfs->content_len;

        ESP_LOGI("webpage", "found page at '%s' with content length: %u", webfs->path, webfs->content_len);

        if (mem >= &webfs_end[0])
        {
            webfs->next = NULL;
            break;
        }

        webfs->next = malloc(sizeof(webfs_t));
        webfs = (webfs_t *)webfs->next;
    }

    return page;
}

void webpage_free(webpage_t *page)
{
    const webfs_t *webfs = page->webfs;
    const webfs_t *next_webfs;

    while (webfs != NULL)
    {
        next_webfs = webfs->next;
        free((void *)webfs);
        webfs = next_webfs;
    }
}

const char *webpage_find(webpage_t *page, const char *path, size_t *content_len)
{
    const webfs_t *webfs = page->webfs;
    while (webfs != NULL)
    {
        if (path_cmp(webfs->path, path))
        {
            *content_len = webfs->content_len;
            return webfs->content;
        }
        webfs = webfs->next;
    }
    *content_len = 0;
    return NULL;
}
