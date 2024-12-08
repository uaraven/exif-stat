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

type imageFileMmap struct {
	Path             string
	File             *os.File
	Data             mmap.MMap
	Order            byte
	Reader           *bytes.Reader
	TiffHeaderOffset int64
}

// OpenExifFileMMap opens file for reading by mmapping it
func OpenExifFileMMap(filepath string) (file File, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Failed to open file %s, Error: %v", filepath, r))
			file = nil
			err = fmt.Errorf("Failed to open file %s, Error: %v", filepath, r)
		}
	}()
	f, err := os.Open(filepath)
	defer func() {
		if err != nil && f != nil {
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
	file = &imageFileMmap{
		Path:   filepath,
		File:   f,
		Data:   data,
		Reader: reader,
		Order:  BigEndian,
	}
	word, err := file.ReadUint16()
	if word != 0xFFD8 {
		return nil, fmt.Errorf("Not a JPEG file? %s", filepath)
	}

	return file, nil
}

func (file imageFileMmap) ReadUint16() (uint16, error) {
	var word uint16
	err := binary.Read(file.Reader, file.getByteOrder(), &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file imageFileMmap) ReadUint32() (uint32, error) {
	var word uint32
	err := binary.Read(file.Reader, file.getByteOrder(), &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file imageFileMmap) readBytes(size uint16) ([]byte, error) {
	buf := make([]byte, size)
	_, err := file.Reader.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (file imageFileMmap) currentPosition() (int64, error) {
	return file.Reader.Seek(0, io.SeekCurrent)
}

func (file imageFileMmap) Seek(pos int64) (int64, error) {
	return file.Reader.Seek(pos, io.SeekStart)
}

func (file imageFileMmap) seekRelative(pos int64) (int64, error) {
	return file.Reader.Seek(pos, io.SeekCurrent)
}

func (file imageFileMmap) Read(out interface{}) error {
	return binary.Read(file.Reader, file.getByteOrder(), out)
}

// Close closes the underlying file
func (file imageFileMmap) Close() {
	file.Data.Unmap()
	file.File.Close()
}

func (file imageFileMmap) GetFile() *os.File {
	return file.File
}

func (file imageFileMmap) getByteOrder() binary.ByteOrder {
	if file.Order == BigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian

}

func (file *imageFileMmap) SetOrder(newOrder byte) {
	file.Order = newOrder
}

func (file imageFileMmap) GetOrder() byte {
	return file.Order
}

func (file imageFileMmap) GetTiffHeaderOffset() int64 {
	return file.TiffHeaderOffset
}

func (file imageFileMmap) GetPath() string {
	return file.Path
}

func (file *imageFileMmap) SetTiffHeaderOffset(newOffset int64) {
	file.TiffHeaderOffset = newOffset
}
