#pragma once

#include <esp_http_server.h>

#include "webpage.h"

typedef struct {
    httpd_handle_t handle;
    webpage_t *webpage;
} webserver_t;

webserver_t webserver_init(webpage_t *webpage);
void webserver_start(webserver_t *webserver);
void webserver_stop(webserver_t *webserver);
