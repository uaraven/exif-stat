package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	imageWidth           = "0100"
	imageHeight          = "0101"
	cameraMake           = "010f"
	model                = "0110"
	orientation          = "0112"
	exposureTime         = "8769/829a"
	fNumber              = "8769/929d"
	iso                  = "8769/8827"
	createTime           = "8769/9004"
	focalLength          = "8769/920a"
	focalLength35        = "8769/a405"
	flash                = "8769/9209"
	exposureProgram      = "8769/8822"
	exposureCompensation = "8769/9204"
	nikonIso             = "8769/927c/0002"
)

const (
	// ExifDataMarker is an identifier of Exif Data marker
	exifDataMarker = 0xFFE1

	exifTagID       = 0x8769
	makerNotesTagID = 0x927c

	// TypeUnknown is an unknown Tag type
	TypeUnknown = 0
	// TypeUnsignedByte is byte
	TypeUnsignedByte = 1
	// TypeASCIItring is sequence of bytes as an ascii string
	TypeASCIItring = 2
	// TypeUnsignedShort is an uint16
	TypeUnsignedShort = 3
	// TypeUnsignedLong is an uint32
	TypeUnsignedLong = 4
	// TypeUnsignedRational is an {uint32, uint32}
	TypeUnsignedRational = 5
	// TypeSignedByte is an int8
	TypeSignedByte = 6
	// TypeUndefined is undefined type, treated as octets
	TypeUndefined = 7
	// TypeSignedShort is an int16
	TypeSignedShort = 8
	// TypeSignedLong is an int32
	TypeSignedLong = 9
	// TypeSignedRational is an {int32, int32}
	TypeSignedRational = 10
	// TypeSingleFloat is a float32
	TypeSingleFloat = 11
	// TypeDoubleFloat is a float64
	TypeDoubleFloat = 12
)

// limited set of known tag names that is used in exif-stat
var tagNames = map[string]string{
	imageWidth:           "Width",
	imageHeight:          "Height",
	cameraMake:           "Make",
	model:                "Model",
	orientation:          "Orientation",
	exposureTime:         "Exposure Time",
	fNumber:              "F-Number",
	iso:                  "ISO",
	createTime:           "Create Time",
	focalLength:          "Focal Length",
	focalLength35:        "Focal Length in 35mm",
	flash:                "Flash",
	exposureProgram:      "Exposure Program",
	exposureCompensation: "Exposure Compensation",
	nikonIso:             "ISO",
}

type tagReader func(file *File, count uint32) (interface{}, []byte, error)

var (
	dataFormatTypes = []tagReader{
		unsignedByteReader,     // unsigned byte
		asciiStringReader,      // ascii string
		unsignedShortReader,    // unsigned short
		unsignedLongReader,     // unsigned long
		unsignedRationalReader, // unsigned rational
		signedByteReader,       // signed byte
		undefinedReader,        // undefined
		signedShortReader,      // signed short
		signedLongReader,       // signed long
		signedRationalReader,   // signed rational
		float32Reader,          // single float
		float64Reader}          // double float
)

func readRawData(file *File, count uint32, bytesInElement uint32) ([]byte, error) {
	size := count * bytesInElement
	if size < 4 {
		size = 4
	}
	rawData := make([]byte, size) // we will buffer to hold at least 4 bytes
	var err error
	if count*bytesInElement > 4 {
		var offset uint32
		err = file.Read(&offset)
		if err != nil {
			return nil, err
		}
		pos, err := file.currentPosition()
		if err != nil {
			return nil, err
		}
		_, err = file.seek(file.TiffHeaderOffset + int64(offset))
		if err != nil {
			return nil, err
		}
		defer func() { file.seek(pos) }()
	}
	err = file.Read(rawData)
	return rawData, err
}

func asciiStringReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 1)
	if err != nil {
		return nil, nil, err
	}
	if rawData[len(rawData)-1] == 0 {
		rawData = rawData[:len(rawData)-1]
	}
	return string(rawData), rawData, nil
}

func unsignedByteReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 1)
	if err != nil {
		return nil, nil, err
	}
	return rawData, rawData, nil
}

func signedByteReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 1)
	if err != nil {
		return nil, nil, err
	}
	signedBytes := make([]int8, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &signedBytes)
	if err != nil {
		return nil, nil, err
	}
	return signedBytes, rawData, nil
}

func unsignedShortReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 2)
	if err != nil {
		return nil, nil, err
	}
	shorts := make([]uint16, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &shorts)
	if err != nil {
		return nil, nil, err
	}
	return shorts, rawData, nil
}

func undefinedReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 1)
	if err != nil {
		return nil, nil, err
	}
	return rawData, rawData, nil
}

func signedShortReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 2)
	if err != nil {
		return nil, nil, err
	}
	shorts := make([]int16, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &shorts)
	if err != nil {
		return nil, nil, err
	}
	return shorts, rawData, nil
}

func unsignedLongReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 4)
	if err != nil {
		return nil, nil, err
	}
	longs := make([]uint32, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &longs)
	if err != nil {
		return nil, nil, err
	}
	return longs, rawData, nil
}

func signedLongReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 4)
	if err != nil {
		return nil, nil, err
	}
	longs := make([]int32, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &longs)
	if err != nil {
		return nil, nil, err
	}
	return longs, rawData, nil
}

func unsignedRationalReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 8)
	if err != nil {
		return nil, nil, err
	}
	longs := make([]uint32, count*2)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &longs)
	if err != nil {
		return nil, nil, err
	}
	rationals := make([]Rational, count)
	for index := uint32(0); index < count; index += 2 {
		rationals[index/2] = NewRational(longs[index], longs[index+1])
	}
	return rationals, rawData, nil
}

func signedRationalReader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 8)
	if err != nil {
		return nil, nil, err
	}
	longs := make([]int32, count*2)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &longs)
	if err != nil {
		return nil, nil, err
	}
	rationals := make([]SignedRational, count)
	for index := uint32(0); index < count; index += 2 {
		rationals[index/2] = NewSignedRational(longs[index], longs[index+1])
	}
	return rationals, rawData, nil
}

func float64Reader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 8)
	if err != nil {
		return nil, nil, err
	}
	floats := make([]float64, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &floats)
	if err != nil {
		return nil, nil, err
	}
	return floats, rawData, nil
}

func float32Reader(file *File, count uint32) (interface{}, []byte, error) {
	rawData, err := readRawData(file, count, 4)
	if err != nil {
		return nil, nil, err
	}
	floats := make([]float32, count)
	err = binary.Read(bytes.NewReader(rawData), file.Order, &floats)
	if err != nil {
		return nil, nil, err
	}
	return floats, rawData, nil
}

func readMarker(file *File) (*marker, error) {
	markerID, err := file.readUint16()
	if err != nil {
		return nil, err
	}
	size, err := file.readUint16()
	if err != nil {
		return nil, err
	}
	pos, err := file.currentPosition()
	if err != nil {
		return nil, err
	}
	return &marker{
		Marker: markerID,
		Size:   size,
		Offset: pos,
	}, nil
}

func readIfdEntry(file *File) (*ifdEntry, error) {
	var tagNumber uint16
	err := file.Read(&tagNumber)
	if err != nil {
		return nil, err
	}

	var dataFormat uint16
	err = file.Read(&dataFormat)
	if err != nil {
		return nil, err
	}
	if dataFormat < 1 || dataFormat > 12 {
		return nil, fmt.Errorf("Unsupported data type format: %d of Tag ID %0x", dataFormat, tagNumber)
	}

	var numComponents uint32
	err = file.Read(&numComponents)
	if err != nil {
		return nil, err
	}

	// read data field and seek back so that dataformat reader can read it again
	var dataValue uint32
	err = file.Read(&dataValue)
	if err != nil {
		return nil, err
	}
	_, err = file.seekRelative(-4)
	if err != nil {
		return nil, err
	}

	value, rawData, err := dataFormatTypes[dataFormat-1](file, numComponents)
	if err != nil {
		return nil, err
	}
	return &ifdEntry{ComponentCount: numComponents, TagID: tagNumber, DataType: dataFormat, Data: dataValue, Value: value, ValueBytes: rawData}, nil
}

func readIfd(file *File, offset int64) (*ifd, error) {
	if offset > 0 {
		pos, err := file.currentPosition()
		if err != nil {
			return nil, err
		}
		defer func() { file.seek(pos) }()
		_, err = file.seek(file.TiffHeaderOffset + offset)
		if err != nil {
			return nil, err
		}
	}
	var err error
	var numEntries uint16
	err = file.Read(&numEntries)
	if err != nil {
		return nil, err
	}
	entries := make([]ifdEntry, 0)
	for index := uint16(0); index < numEntries; index++ {
		entry, err := readIfdEntry(file)
		if err != nil { // ignore the invalid entry
			return nil, err
		}
		entries = append(entries, *entry)
	}
	return &ifd{EntryCount: numEntries, IfdEntries: entries}, nil
}

func entryToTag(parents []uint16, entry ifdEntry) Tag {
	return Tag{
		ID:       entry.TagID,
		IDPath:   parents,
		DataType: int(entry.DataType),
		Value:    entry.Value,
		RawData:  entry.ValueBytes,
	}
}

func entriesToTags(parentIDs []uint16, file *File, entries []ifdEntry) (Tags, error) {
	tags := make([]Tag, 0)
	for _, entry := range entries {
		if entry.TagID == exifTagID {
			exifTagEntries, err := readIfd(file, int64(entry.Value.([]uint32)[0]))
			if err != nil {
				return nil, err
			}
			parents := append(parentIDs, entry.TagID)
			exifTags, err := entriesToTags(parents, file, exifTagEntries.IfdEntries)
			if err != nil {
				return nil, err
			}
			for _, tag := range exifTags {
				tags = append(tags, tag)
			}
		} else if entry.TagID == makerNotesTagID {
			exifTagEntries, err := readMakerNotes(file, entry)
			if err != nil {
				return nil, err
			}
			if exifTagEntries != nil {
				parents := append(parentIDs, entry.TagID)
				exifTags, err := entriesToTags(parents, file, exifTagEntries.IfdEntries)
				if err != nil {
					return nil, err
				}
				for _, tag := range exifTags {
					tags = append(tags, tag)
				}
			}
		} else {
			tag := entryToTag(parentIDs, entry)
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func readExifHeader(file *File, marker *marker) (*ifd, error) {
	// check headers
	file.seek(marker.Offset)
	// examine exif header
	var exifMagic uint32
	err := file.Read(&exifMagic)
	if err != nil {
		return nil, err
	}
	if exifMagic != 0x45786966 {
		return nil, fmt.Errorf("Exif data does not contain valid exif marker")
	}
	var word uint16
	err = file.Read(&word)
	if err != nil {
		return nil, err
	}
	if word != 0 {
		return nil, fmt.Errorf("Exif data does not contain valid exif marker")
	}
	tiffHeaderOffset, err := file.currentPosition()
	if err != nil {
		return nil, err
	}

	err = readTiffHeader(file)
	if err != nil {
		return nil, err
	}

	file.TiffHeaderOffset = tiffHeaderOffset

	// we're at IFD0 and can start reading IFD
	return readIfd(file, -1)
}

func readTiffHeader(file *File) error {
	// examine TIFF header

	var wword uint32
	err := file.Read(&wword)
	if err != nil {
		return err
	}
	if wword == 0x49492A00 {
		file.Order = binary.LittleEndian
	} else if wword == 0x4d4d002A {
		file.Order = binary.BigEndian
	} else {
		return fmt.Errorf("Invalid byte order in TIFF header %x", wword)
	}
	err = file.Read(&wword)
	if err != nil {
		return err
	}
	_, err = file.seekRelative(int64(wword - 8)) // relative offset from the start of the TIFF header
	if err != nil {
		return err
	}
	return nil
}

// ReadExifTags parses file, extracts Ifds from it and parses ifds for all tags
func ReadExifTags(file *File) (Tags, error) {
	// find exif marker in the file
	var marker *marker
	var err error
	for {
		marker, err = readMarker(file)
		if err != nil {
			return nil, err
		}
		if marker.Marker == exifDataMarker {
			break
		} else {
			file.seekRelative(int64(marker.Size - 2))
		}
	}
	if marker.Marker != exifDataMarker {
		return nil, fmt.Errorf("Cannot find exif data in %s", file.Path)
	}
	ifd, err := readExifHeader(file, marker)
	if err != nil {
		return nil, err
	}
	parent := make([]uint16, 0)
	return entriesToTags(parent, file, ifd.IfdEntries)
}

// TagsAsMap converts list of tags into a map of tag path -> tag
func TagsAsMap(tags Tags) map[string]Tag {
	result := make(map[string]Tag, 0)
	for _, tag := range tags {
		result[tag.PathName()] = tag
	}
	return result
}
