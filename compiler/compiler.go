package compiler

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var defIgnoreMap = map[string]bool{
	".git": true,
}

func Compile(rootDir string, outputPath string, compress, addHeaders bool, extensions map[string]string, ignore []string) error {
	extMap := defExtensionsMap
	for key, val := range extensions {
		key = strings.ToLower(key)
		if !strings.HasPrefix(key, ".") {
			key = "." + key
		}
		extMap[key] = val
	}

	ignoreMap := defIgnoreMap
	for _, item := range ignore {
		ignoreMap[item] = true
	}

	buff := &bytes.Buffer{}
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := path[len(rootDir):]

		if d.IsDir() {
			if ignoreMap[relPath] {
				return filepath.SkipDir
			}
			return nil
		}

		if ignoreMap[relPath] {
			return nil
		}

		ext := filepath.Ext(relPath)
		header, ok := extMap[strings.ToLower(ext)]
		if !ok {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		return writeFile(buff, file, relPath, int(info.Size()), compress, addHeaders, header)
	})

	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, buff.Bytes(), 0755)
	return err
}

func writeFile(wr io.Writer, r io.Reader, path string, fileSize int, compress, addHeaders bool, header string) error {
	var err error
	response := &bytes.Buffer{}

	// Write HTTP response

	if addHeaders {
		_, err = response.Write([]byte(HEADER_OK + "\r\n"))
		if err != nil {
			return err
		}
		_, err = response.Write([]byte(HEADER_CONTENT_TYPE + header + "\r\n"))
		if err != nil {
			return err
		}
	}

	if compress {
		if addHeaders {
			_, err = response.Write([]byte(HEADER_GZIP_ENCODING + "\r\n"))
			if err != nil {
				return err
			}
		}

		compressed := &bytes.Buffer{}
		compressor := gzip.NewWriter(compressed)

		_, err = io.Copy(compressor, r)
		if err != nil {
			return err
		}

		err = compressor.Close()
		if err != nil {
			return err
		}

		fileSize = compressed.Len()
		r = compressed
	}

	if addHeaders {
		_, err = response.Write([]byte(HEADER_CONTENT_LENGTH + strconv.Itoa(fileSize) + "\r\n\r\n"))
		if err != nil {
			return err
		}
	}

	_, err = io.Copy(response, r)
	if err != nil {
		return err
	}

	_, err = response.Write([]byte{0})
	if err != nil {
		return err
	}

	// Write file info

	err = binary.Write(wr, binary.LittleEndian, uint32(len(path))+1)
	if err != nil {
		return err
	}

	_, err = wr.Write([]byte(path))
	if err != nil {
		return err
	}

	_, err = wr.Write([]byte{0})
	if err != nil {
		return err
	}

	err = binary.Write(wr, binary.LittleEndian, uint32(response.Len()))
	if err != nil {
		return err
	}

	// Write respose to output file
	_, err = io.Copy(wr, response)
	return err
}
