#include "webserver.h"

#include <esp_log.h>

#define TAG "HTTP"

esp_err_t http_get_handler(httpd_req_t *req)
{
    webserver_t *server = (webserver_t *)req->user_ctx;
    size_t content_len;

    ESP_LOGI(TAG, "GET request: '%s'", req->uri);
    const char *resp = webpage_find(server->webpage, req->uri, &content_len);
    if (resp == NULL)
    {
        httpd_resp_send_404(req);
        return ESP_FAIL;
    }

    httpd_send(req, resp, content_len);
    return ESP_OK;
}

esp_err_t http_get_index_handler(httpd_req_t *req)
{
    httpd_resp_set_type(req, "text/html");
    httpd_resp_set_status(req, "301 Moved Permanently");
    httpd_resp_set_hdr(req, "Location", "index.html");
    httpd_resp_send(req, NULL, 0);
    return ESP_OK;
}

httpd_uri_t http_get_index = {
    .uri = "/",
    .method = HTTP_GET,
    .handler = http_get_index_handler
};

webserver_t webserver_init(webpage_t *webpage)
{
    webserver_t server = {
        .handle = NULL,
        .webpage = webpage
    };
    return server;
}

void webserver_start(webserver_t *webserver)
{
    httpd_handle_t server = NULL;
    httpd_config_t config = HTTPD_DEFAULT_CONFIG();
    config.max_uri_handlers = 40;

    httpd_uri_t http_get = {
        .method = HTTP_GET,
        .handler = http_get_handler,
        .user_ctx = webserver
    };

    // Start the httpd server
    ESP_LOGI(TAG, "starting webserver on port: '%d'", config.server_port);
    if (httpd_start(&server, &config) != ESP_OK)
    {
        ESP_LOGI(TAG, "error starting server!");
        return;
    }

    // Set URI handlers
    httpd_register_uri_handler(server, &http_get_index);
    const webfs_t *webfs = webserver->webpage->webfs;
    while (webfs != NULL)
    {
        http_get.uri = webfs->path;
        ESP_LOGI(TAG, "registering URI handler for %s", webfs->path);
        httpd_register_uri_handler(server, &http_get);
        webfs = webfs->next;
    }

    webserver->handle = server;
    return;
}

void webserver_stop(webserver_t *webserver)
{
    if (webserver->handle == NULL) return;
    ESP_LOGI(TAG, "stopping webserver");
    httpd_stop(webserver->handle);
    webserver->handle = NULL;
}
