package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/mmap-go"
	"github.com/uaraven/exif-stat/logger"
)

// File contains data required to read exif information from a file
type File struct {
	Path             string
	File             *os.File
	Data             mmap.MMap
	Order            binary.ByteOrder
	Reader           *bytes.Reader
	TiffHeaderOffset int64
}

// OpenExifFile opens file for reading
func OpenExifFile(filepath string) (*File, error) {
	f, err := os.Open(filepath)
	defer func() {
		if err != nil {
			f.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	var file = &File{
		Path:   filepath,
		File:   f,
		Data:   data,
		Reader: reader,
		Order:  binary.BigEndian,
	}
	word, err := file.readUint16()
	if word != 0xFFD8 {
		return nil, fmt.Errorf("Not a JPEG file? %s", filepath)
	}

	return file, nil
}

func (file *File) readUint16() (uint16, error) {
	var word uint16
	err := binary.Read(file.Reader, file.Order, &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file *File) readUint32() (uint32, error) {
	var word uint32
	err := binary.Read(file.Reader, file.Order, &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file *File) readBytes(size uint16) ([]byte, error) {
	buf := make([]byte, size)
	_, err := file.Reader.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (file *File) currentPosition() (int64, error) {
	return file.Reader.Seek(0, io.SeekCurrent)
}

func (file *File) seek(pos int64) (int64, error) {
	return file.Reader.Seek(pos, io.SeekStart)
}

func (file *File) seekRelative(pos int64) (int64, error) {
	return file.Reader.Seek(pos, io.SeekCurrent)
}

func (file *File) Read(out interface{}) error {
	return binary.Read(file.Reader, file.Order, out)
}

// Close closes the underlying file
func (file *File) Close() {
	logger.Verbose(2, fmt.Sprintf("Closing file %v", file.Path))
	file.Data.Unmap()
	file.File.Close()
}
