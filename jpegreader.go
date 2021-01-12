package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/uaraven/exif-stat/exif"
	"github.com/uaraven/exif-stat/logger"
)

// ExifInfo contains values of all the exif tag of interest
type ExifInfo struct {
	Make                 string
	Model                string
	CreateTime           string
	Iso                  uint16
	FNumber              exif.Rational
	ExposureTime         exif.Rational
	FocalLength          exif.Rational
	FocalLength35        uint16
	Flash                string
	ExposureProgram      string
	ExposureCompensation exif.SignedRational
	Width                uint32
	Height               uint32
	FileName             string
}

func (ei *ExifInfo) toString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Make: %s\n", ei.Make))
	sb.WriteString(fmt.Sprintf("Model: %s\n", ei.Model))
	sb.WriteString(fmt.Sprintf("CreateTime: %s\n", ei.CreateTime))
	sb.WriteString(fmt.Sprintf("Iso: %d\n", ei.Iso))
	sb.WriteString(fmt.Sprintf("FNumber: %s\n", ei.FNumber.ToString()))
	sb.WriteString(fmt.Sprintf("Exposure time: %s\n", ei.ExposureTime.ToString()))
	sb.WriteString(fmt.Sprintf("Focal length: %f\n", ei.FocalLength.AsFloat()))
	sb.WriteString(fmt.Sprintf("Focal length in 35mm: %d\n", ei.FocalLength35))
	sb.WriteString(fmt.Sprintf("Flash: %s\n", ei.Flash))
	sb.WriteString(fmt.Sprintf("Exposure program: %s\n", ei.ExposureProgram))
	return sb.String()
}

func (ei *ExifInfo) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Make":                 ei.Make,
		"Model":                ei.Model,
		"CreateTime":           ei.CreateTime,
		"Iso":                  ei.Iso,
		"FNumber":              ei.FNumber.AsFloat(),
		"ExposureTime":         ei.ExposureTime.ToString(),
		"FocalLength":          ei.FocalLength.AsFloat(),
		"FocalLength35":        ei.FocalLength35,
		"Flash":                ei.Flash,
		"ExposureProgram":      ei.ExposureProgram,
		"ExposureCompensation": ei.ExposureCompensation.Normalize().ToString(),
		"Width":                ei.Width,
		"Height":               ei.Height,
	}
}

type tagValueExtractor = func(tag exif.Tag, exifInfo *ExifInfo)

const (
	tagMake                 = "010f"
	tagModel                = "0110"
	tagDateTimeOriginal     = "8769/9003"
	tagIso                  = "8769/8827"
	tagFNumber              = "8769/829d"
	tagExposureTime         = "8769/829a"
	tagFocalLength          = "8769/920a"
	tagFocalLength35        = "8769/a405"
	tagFlash                = "8769/9209"
	tagExposureProgram      = "8769/8822"
	tagExposureCompensation = "8769/9204"
	tagMeteringMode         = "8769/9207"
	tagMaxAperture          = "8769/9205"
	tagOrientation          = "0112"
	tagImageWidth           = "8769/a002"
	tagImageHeight          = "8769/a003"
	tagExposureMore         = "8769/a402"
	tagNikonIso             = "8769/927c/0002"
)

var exifFlashValues = map[uint]string{
	0x0:  "No Flash",
	0x1:  "Fired",
	0x5:  "Fired, Return not detected",
	0x7:  "Fired, Return detected",
	0x8:  "On, Did not fire",
	0x9:  "On, Fired",
	0xd:  "On, Return not detected",
	0xf:  "On, Return detected",
	0x10: "Off, Did not fire",
	0x14: "Off, Did not fire, Return not detected",
	0x18: "Auto, Did not fire",
	0x19: "Auto, Fired",
	0x1d: "Auto, Fired, Return not detected",
	0x1f: "Auto, Fired, Return detected",
	0x20: "No flash function",
	0x30: "Off, No flash function",
	0x41: "Fired, Red-eye reduction",
	0x45: "Fired, Red-eye reduction, Return not detected",
	0x47: "Fired, Red-eye reduction, Return detected",
	0x49: "On, Red-eye reduction",
	0x4d: "On, Red-eye reduction, Return not detected",
	0x4f: "On, Red-eye reduction, Return detected",
	0x50: "Off, Red-eye reduction",
	0x58: "Auto, Did not fire, Red-eye reduction",
	0x59: "Auto, Fired, Red-eye reduction",
	0x5d: "Auto, Fired, Red-eye reduction, Return not detected",
	0x5f: "Auto, Fired, Red-eye reduction, Return detected",
}

var exifExposurePrograms = map[uint]string{
	0: "Not Defined",
	1: "Manual",
	2: "Program AE",
	3: "Aperture-priority AE",
	4: "Shutter speed priority AE",
	5: "Creative (Slow speed)",
	6: "Action (High speed)",
	7: "Portrait",
	8: "Landscape",
	9: "Bulb",
}

var extractors = map[string]tagValueExtractor{
	tagMake: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.Make = tag.Value.(string)
	},
	tagModel: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.Model = tag.Value.(string)
	},
	tagDateTimeOriginal: func(tag exif.Tag, exifInfo *ExifInfo) {
		tm, err := parseExifFullTimestamp(tag.Value.(string))
		if err == nil {
			exifInfo.CreateTime = tm.Format(time.RFC3339)
		} else {
			exifInfo.CreateTime = ""
		}
	},
	tagIso: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.Iso = tag.Value.([]uint16)[0]
	},
	tagFNumber: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.FNumber = tag.Value.([]exif.Rational)[0]
	},
	tagExposureTime: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.ExposureTime = tag.Value.([]exif.Rational)[0]
	},
	tagFocalLength: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.FocalLength = tag.Value.([]exif.Rational)[0]
	},
	tagFocalLength35: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.FocalLength35 = tag.Value.([]uint16)[0]
	},
	tagFlash: func(tag exif.Tag, exifInfo *ExifInfo) {
		val, ok := exifFlashValues[uint(tag.Value.([]uint16)[0])]
		if ok {
			exifInfo.Flash = val
		}
	},
	tagExposureProgram: func(tag exif.Tag, exifInfo *ExifInfo) {
		val, ok := exifExposurePrograms[uint(tag.Value.([]uint16)[0])]
		if ok {
			exifInfo.ExposureProgram = val
		}
	},
	tagExposureCompensation: func(tag exif.Tag, exifInfo *ExifInfo) {
		exifInfo.ExposureCompensation = tag.Value.([]exif.SignedRational)[0]
	},
	tagImageWidth: func(tag exif.Tag, exifInfo *ExifInfo) {
		switch tag.DataType {
		case exif.TypeUnsignedLong:
			exifInfo.Width = tag.Value.([]uint32)[0]
		case exif.TypeSignedLong:
			exifInfo.Width = uint32(tag.Value.([]int32)[0])
		default:
			exifInfo.Width = uint32(tag.Value.([]uint16)[0])
		}
	},
	tagImageHeight: func(tag exif.Tag, exifInfo *ExifInfo) {
		switch tag.DataType {
		case exif.TypeUnsignedLong:
			exifInfo.Height = tag.Value.([]uint32)[0]
		case exif.TypeSignedLong:
			exifInfo.Height = uint32(tag.Value.([]int32)[0])
		default:
			exifInfo.Height = uint32(tag.Value.([]uint16)[0])
		}
	},
}

func extractNikonIso(tag exif.Tag, exifInfo *ExifInfo) {
	if tag.DataType != exif.TypeUnsignedShort { // sometimes there is TypeUndefined and all zeroes here
		logger.Verbose(2, fmt.Sprintf("\nUnexpected data type in tag %v in exifinfo: %v", tag, exifInfo))
		return
	}
	exifInfo.Iso = tag.Value.([]uint16)[1]
}

func parseExifFullTimestamp(timestamp string) (*time.Time, error) {
	parts := strings.Split(timestamp, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid timestamp %s", timestamp)
	}
	datestampValue, timestampValue := parts[0], parts[1]

	// Normalize the separators.
	datestampValue = strings.ReplaceAll(datestampValue, "-", ":")
	timestampValue = strings.ReplaceAll(timestampValue, "-", ":")

	dateParts := strings.Split(datestampValue, ":")

	year, err := strconv.ParseUint(dateParts[0], 10, 16)
	if err != nil {
		return nil, err
	}

	month, err := strconv.ParseUint(dateParts[1], 10, 8)
	if err != nil {
		return nil, err
	}

	day, err := strconv.ParseUint(dateParts[2], 10, 8)
	if err != nil {
		return nil, err
	}

	timeParts := strings.Split(timestampValue, ":")

	hour, err := strconv.ParseUint(timeParts[0], 10, 8)
	if err != nil {
		return nil, err
	}

	minute, err := strconv.ParseUint(timeParts[1], 10, 8)
	if err != nil {
		return nil, err
	}

	second, err := strconv.ParseUint(timeParts[2], 10, 8)
	if err != nil {
		return nil, err
	}

	res := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)
	return &res, nil
}

// ExtractExif parses image file with a given path and extracts exif information
func ExtractExif(imageFilePath string) (*ExifInfo, error) {

	file, err := exif.OpenExifFile(imageFilePath)
	if err != nil {
		return nil, err
	}
	defer func() { file.Close() }()

	tags, err := exif.ReadExifTags(file)
	if err != nil {
		return nil, err
	}

	tagMap := exif.TagsAsMap(tags)
	exifInfo := &ExifInfo{
		FileName: imageFilePath,
	}

	for path, extractor := range extractors {
		tag, ok := tagMap[path]
		if ok {
			extractor(tag, exifInfo)
		}
	}
	if _, ok := tagMap[tagIso]; !ok { // no standard ISO tag
		if tag, ok := tagMap[tagNikonIso]; ok { // but there is Nikon-specific ISO tag
			extractNikonIso(tag, exifInfo)
		}
	}

	return exifInfo, nil
}
