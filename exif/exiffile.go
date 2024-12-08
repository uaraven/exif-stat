package exif

import (
	"encoding/binary"
	"os"
)

const (
	// BigEndian representation to read multibyte data from file
	BigEndian = 'M'
	// LittleEndian representation to read multibyte data from file
	LittleEndian = 'I'
)

// File contains data required to read exif information from a file
type File interface {
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	readBytes(uint16) ([]byte, error)
	currentPosition() (int64, error)
	Seek(pos int64) (int64, error)
	seekRelative(int64) (int64, error)
	Read(interface{}) error
	Close()
	GetPath() string
	GetFile() *os.File
	getByteOrder() binary.ByteOrder
	GetOrder() byte
	SetOrder(byte)
	GetTiffHeaderOffset() int64
	SetTiffHeaderOffset(int64)
}
