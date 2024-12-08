package main

import (
	"fmt"
	"strings"

	"github.com/uaraven/exif-stat/exif"
	"github.com/uaraven/exif-stat/logger"
)

type JpegReader struct {
	fileName string
	useMmap  bool
}

var _ ExifFileReader = &JpegReader{}

func NewJpegReader(fileName string, useMmap bool) *JpegReader {
	return &JpegReader{fileName, useMmap}
}

func (jp JpegReader) ReadExif() (exifInfo *ExifInfo, err error) {
	exifInfo = &ExifInfo{
		FileName: jp.fileName,
	}
	defer func() {
		state := recover()
		if state != nil {
			logger.Verbose(2, fmt.Sprintf("failure while reading %s: %v", jp.fileName, state))
			exifInfo = nil
			err = fmt.Errorf("failure while reading %s: %v", jp.fileName, state)
		}
	}()
	var f exif.File
	if jp.useMmap {
		f, err = exif.OpenExifFileMMap(jp.fileName)
	} else {
		f, err = exif.OpenExifFileIo(jp.fileName)
	}
	if err != nil {
		return nil, err
	}
	defer func() { f.Close() }()

	word, err := f.ReadUint16()
	if word != ExifMagic {
		return nil, fmt.Errorf("not a JPEG file? %s", jp.fileName)
	}

	exifInfo, err = readExifFromFile(f, jp.fileName)
	return
}

func (jp JpegReader) IsFileSupported() bool {
	ff := strings.ToLower(jp.fileName)
	return strings.HasSuffix(ff, ".jpg") || strings.HasSuffix(ff, ".jpeg")
}
