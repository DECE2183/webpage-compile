# webpage-compile

A simple compiler for embedding web pages into microcontroller flash memory.

## Usage

```bash
webpage-compile [flags] [path to config]
```

If no path is specified, `webpage.yaml` in the current location will be used.

### Flags

| Flag                 | Description             |
| -------------------- | ----------------------- |
| `-h`, `--help`       | display help message    |
| `-c`, `--gen-config` | generate default config |

### Config

The default config looks like:

```yaml
input: .
output: webpage.bin
compress: false
omit-headers: false
extensions: {}
ignore: []
```

All paths in the config file are relative to the directory
in which the config file is located.

The 'extensions' field is a map where the key is the file extension
and the value is the MIME type for that file type.

Extensions example:

```yaml
extensions:
    - .bin: application/octet-stream
    - .epub: application/epub+zip
```

### Embedding

Embedding examples are located in the [examples](examples)
directory of this repository.

The resulting file consists of sequential structures,
which can be described using C-like syntax as follows:

```C
struct page {
    uint32_t path_len;
    uint8_t  path[path_len];
    uint32_t content_len;
    uint8_t  content[content_len];
};
```

`path` and `content` are null-terminated strings.

This file can be placed anywhere in the microcontroller's flash memory and then accessed by its address. All the http server needs to do is send the `content` of the structure with the `path` corresponding to the request URL.
