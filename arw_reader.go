package main

import (
	"fmt"
	"strings"

	"github.com/uaraven/exif-stat/exif"
	"github.com/uaraven/exif-stat/logger"
)

// ArwReader reads Exif tags from Sony ARW RAW file
type ArwReader struct {
	fileName string
	useMmap  bool
}

func NewArwReader(fileName string, useMmap bool) *ArwReader {
	return &ArwReader{fileName: fileName, useMmap: useMmap}
}

var _ ExifFileReader = &ArwReader{}

func (rr ArwReader) ReadExif() (exifInfo *ExifInfo, err error) {
	defer func() {
		state := recover()
		if state != nil {
			logger.Verbose(2, fmt.Sprintf("failure while reading %s: %v", rr.fileName, state))
			exifInfo = nil
			err = fmt.Errorf("failure while reading %s: %v", rr.fileName, state)
		}
	}()
	var f exif.File
	if rr.useMmap {
		f, err = exif.OpenExifFileMMap(rr.fileName)
	} else {
		f, err = exif.OpenExifFileIo(rr.fileName)
	}
	if err != nil {
		return nil, err
	}
	defer func() { f.Close() }()

	exifInfo, err = readTiffFromFile(f, rr.fileName)
	return
}

func (rr ArwReader) IsFileSupported() bool {
	isRaf := strings.HasSuffix(strings.ToLower(rr.fileName), ".arw")
	if !isRaf {
		return false
	}
	return true
}
