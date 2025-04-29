package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GetReaderForFileName(filename string, fastIo bool) (ExifFileReader, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	var reader ExifFileReader
	if ext == ".jpg" || ext == ".jpeg" {
		reader = NewJpegReader(filename, fastIo)
	} else if ext == ".raf" {
		reader = NewRafReader(filename, fastIo)
	} else if ext == ".arw" {
		reader = NewArwReader(filename, fastIo)
	} else {
		return nil, fmt.Errorf("unsupported file: %s", filename)
	}
	return reader, nil
}
