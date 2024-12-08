package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/uaraven/exif-stat/exif"
	"github.com/uaraven/exif-stat/logger"
)

// RafReader reads Exif tags from Fujifilm RAW file
// See https://libopenraw.freedesktop.org/formats/raf/ for format specs
type RafReader struct {
	fileName string
	useMmap  bool
}

func NewRafReader(fileName string, useMmap bool) *RafReader {
	return &RafReader{fileName: fileName, useMmap: useMmap}
}

var _ ExifFileReader = &RafReader{}

func (rr RafReader) ReadExif() (exifInfo *ExifInfo, err error) {
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
	_, err = f.Seek(16 + 4 + 8 + 32 + 4 + 20)
	if err != nil {
		return nil, err
	}
	builtInJpegOffset, err := f.ReadUint32()
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(int64(builtInJpegOffset))
	if err != nil {
		return nil, err
	}
	word, err := f.ReadUint16()
	if word != ExifMagic {
		return nil, fmt.Errorf("unexpected file data: %s", rr.fileName)
	}

	exifInfo, err = readExifFromFile(f, rr.fileName)
	return
}

func (rr RafReader) IsFileSupported() bool {
	isRaf := strings.HasSuffix(strings.ToLower(rr.fileName), ".raf")
	if !isRaf {
		return false
	}
	f, err := os.Open(rr.fileName)
	if err != nil {
		return false
	}
	defer f.Close()
	magic := make([]byte, 16)
	_, err = f.Read(magic)
	if err != nil || string(magic) != "FUJIFILMCCD-RAW" {
		return false
	}
	return true
}
