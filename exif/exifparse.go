package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/uaraven/exif-stat/utils"
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

	// UnknownType is an unknown Tag type
	UnknownType = 0
	// UnsignedByte is byte
	UnsignedByte = 1
	// ASCIItring is sequence of bytes as an ascii string
	ASCIItring = 2
	// UnsignedShort is an uint16
	UnsignedShort = 3
	// UnsignedLong is an uint32
	UnsignedLong = 4
	// UnsignedRational is an {uint32, uint32}
	UnsignedRational = 5
	// SignedByte is an int8
	SignedByte = 6
	// Undefined is undefined type, treated as octets
	Undefined = 7
	// SignedShort is an int16
	SignedShort = 8
	// SignedLong is an int32
	SignedLong = 9
	// SignedRational is an {int32, int32}
	SignedRational = 10
	// SingleFloat is a float32
	SingleFloat = 11
	// DoubleFloat is a float64
	DoubleFloat = 12
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

// TagDataType contains data of the exif tag
type TagDataType struct {
	DataLength int
	Reader     func(file *File, count uint32) (interface{}, []byte, error)
}

var (
	dataFormatTypes = []TagDataType{
		TagDataType{1, unsignedByteReader},     // unsigned byte
		TagDataType{1, asciiStringReader},      // ascii string
		TagDataType{2, unsignedShortReader},    // unsigned short
		TagDataType{4, unsignedLongReader},     // unsigned long
		TagDataType{8, unsignedRationalReader}, // unsigned rational
		TagDataType{1, signedByteReader},       // signed byte
		TagDataType{0, undefinedReader},        // undefined
		TagDataType{2, signedShortReader},      // signed short
		TagDataType{4, signedLongReader},       // signed long
		TagDataType{8, signedRationalReader},   // signed rational
		TagDataType{4, float32Reader},          // single float
		TagDataType{8, float64Reader}}          // double float
)

// Marker contains data of TIFF marker
type Marker struct {
	Marker uint16
	Size   uint16
	Offset int64
}

// IfdEntry is an exif tag
type IfdEntry struct {
	// Tag ID
	TagID uint16
	// Data Type
	DataType uint16
	// Number of components
	ComponentCount uint32
	// Raw data value (can be data or offset)
	Data uint32
	// Parsed and converted data
	Value interface{}
	// Actual data as slice of bytes
	ValueBytes []byte
}

// ToString returns a string representation of IfdEntry
func (ie IfdEntry) ToString() string {
	return fmt.Sprintf("ID=%x Value=%v Bytes=%v", ie.TagID, ie.Value, ie.ValueBytes)
}

// Tag is a simplified representation of an Exif Tag
type Tag struct {
	ID       uint16
	IDPath   []uint16
	DataType int
	Value    interface{}
}

// PathName creates path-line name from tag Id and parent ids
func (tag Tag) PathName() string {
	var sb strings.Builder
	for _, p := range tag.IDPath {
		sb.WriteString(fmt.Sprintf("%04x", p))
		sb.WriteString("/")
	}
	sb.WriteString(fmt.Sprintf("%04x", tag.ID))
	return sb.String()
}

// ToString returns a string representation of IfdEntry
func (tag Tag) ToString() string {
	return fmt.Sprintf("Path=%s ID=%x Value=%v", tag.PathName(), tag.ID, tag.Value)
}

// Tags is a slice of all Exif tags
type Tags []Tag

// Ifd represents image format descriptor
type Ifd struct {
	EntryCount uint16
	IfdEntries []IfdEntry
}

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
	rationals := make([]utils.Rational, count)
	for index := uint32(0); index < count; index += 2 {
		rationals[index/2] = utils.NewRational(longs[index], longs[index+1])
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
	rationals := make([]utils.SignedRational, count)
	for index := uint32(0); index < count; index += 2 {
		rationals[index/2] = utils.NewSignedRational(longs[index], longs[index+1])
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

func readMarker(file *File) (*Marker, error) {
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
	return &Marker{
		Marker: markerID,
		Size:   size,
		Offset: pos,
	}, nil
}

func readIfdEntry(file *File) (*IfdEntry, error) {
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

	value, rawData, err := dataFormatTypes[dataFormat-1].Reader(file, numComponents)
	if err != nil {
		return nil, err
	}
	return &IfdEntry{ComponentCount: numComponents, TagID: tagNumber, DataType: dataFormat, Data: dataValue, Value: value, ValueBytes: rawData}, nil
}

func readIfd(file *File, offset int64) (*Ifd, error) {
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
	entries := make([]IfdEntry, 0)
	for index := uint16(0); index < numEntries; index++ {
		entry, err := readIfdEntry(file)
		if err != nil { // ignore the invalid entry
			return nil, err
		}
		entries = append(entries, *entry)
	}
	return &Ifd{EntryCount: numEntries, IfdEntries: entries}, nil
}

func entryToTag(parents []uint16, entry IfdEntry) Tag {
	return Tag{
		ID:       entry.TagID,
		IDPath:   parents,
		DataType: int(entry.DataType),
		Value:    entry.Value,
	}
}

func entriesToTags(parentIDs []uint16, file *File, entries []IfdEntry) (Tags, error) {
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

func readExifHeader(file *File, marker *Marker) (*Ifd, error) {
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
	var marker *Marker
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
