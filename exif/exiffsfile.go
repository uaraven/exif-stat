package exif

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/uaraven/exif-stat/logger"
)

type imageFileFs struct {
	Path             string
	File             *os.File
	Order            byte
	TiffHeaderOffset int64
}

// OpenExifFileIo opens file on a file system for reading
func OpenExifFileIo(filepath string) (file File, err error) {
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
	file = &imageFileFs{
		Path:  filepath,
		File:  f,
		Order: BigEndian,
	}

	return file, nil
}

func (file imageFileFs) ReadUint16() (uint16, error) {
	var word uint16
	err := binary.Read(file.File, file.getByteOrder(), &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file imageFileFs) ReadUint32() (uint32, error) {
	var word uint32
	err := binary.Read(file.File, file.getByteOrder(), &word)
	if err != nil {
		return 0, err
	}
	return word, nil
}

func (file imageFileFs) readBytes(size uint16) ([]byte, error) {
	buf := make([]byte, size)
	_, err := file.File.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (file imageFileFs) currentPosition() (int64, error) {
	return file.File.Seek(0, os.SEEK_CUR)
}

func (file imageFileFs) Seek(pos int64) (int64, error) {
	return file.File.Seek(pos, os.SEEK_SET)
}

func (file imageFileFs) seekRelative(pos int64) (int64, error) {
	return file.File.Seek(pos, os.SEEK_CUR)
}

func (file imageFileFs) Read(out interface{}) error {
	return binary.Read(file.File, file.getByteOrder(), out)
}

// Close closes the underlying file
func (file imageFileFs) Close() {
	file.File.Close()
}

func (file imageFileFs) GetFile() *os.File {
	return file.File
}

func (file imageFileFs) getByteOrder() binary.ByteOrder {
	if file.Order == BigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian

}

func (file *imageFileFs) SetOrder(newOrder byte) {
	file.Order = newOrder
}

func (file imageFileFs) GetOrder() byte {
	return file.Order
}

func (file imageFileFs) GetTiffHeaderOffset() int64 {
	return file.TiffHeaderOffset
}

func (file imageFileFs) GetPath() string {
	return file.Path
}

func (file *imageFileFs) SetTiffHeaderOffset(newOffset int64) {
	file.TiffHeaderOffset = newOffset
}
