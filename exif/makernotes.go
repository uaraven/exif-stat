package exif

import (
	"encoding/binary"
)

func nikonV3Detector(data []byte) bool {
	header := []byte{'N', 'i', 'k', 'o', 'n', 0x00, 0x02, 0x10, 0x00, 0x00}
	for i, v := range header {
		if v != data[i] {
			return false
		}
	}
	return true
}

func nikonV3VariantDetector(data []byte) bool {
	header := []byte{'N', 'i', 'k', 'o', 'n', 0x00, 0x02, 0x00, 0x00, 0x00}
	for i, v := range header {
		if v != data[i] {
			return false
		}
	}
	return true
}

func nikonV3Reader(file *File, entry IfdEntry) (*Ifd, error) {
	mainTiffHeaderOffset := file.TiffHeaderOffset
	mainOrder := file.Order
	defer func() {
		file.TiffHeaderOffset = mainTiffHeaderOffset
		file.Order = mainOrder
	}()

	offset := file.TiffHeaderOffset + int64(entry.Data) + 10 // + 10 bytes of nikon signature
	_, err := file.seek(offset)
	if err != nil {
		return nil, err
	}
	file.TiffHeaderOffset = offset
	file.Order = binary.BigEndian
	err = readTiffHeader(file)
	if err != nil {
		return nil, err
	}
	return readIfd(file, -1)
}

type makerNoteReader struct {
	CanRead func([]byte) bool
	Reader  func(*File, IfdEntry) (*Ifd, error)
}

var makerNoteReaders = []makerNoteReader{
	makerNoteReader{nikonV3Detector, nikonV3Reader},
	makerNoteReader{nikonV3VariantDetector, nikonV3Reader},
}

func readMakerNotes(file *File, entry IfdEntry) (*Ifd, error) {
	for _, reader := range makerNoteReaders {
		if reader.CanRead(entry.ValueBytes) {
			return reader.Reader(file, entry)
		}
	}
	return nil, nil
}
