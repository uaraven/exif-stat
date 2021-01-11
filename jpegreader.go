package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dsoprea/go-exif"
	log "github.com/dsoprea/go-logging"
)

// ExifInfo contains values of all the exif tag of interest
type ExifInfo struct {
	Make                 string
	Model                string
	CreateTime           time.Time
	Iso                  uint16
	FNumber              Rational
	ExposureTime         Rational
	FocalLength          Rational
	FocalLength35        uint16
	Flash                string
	ExposureProgram      string
	ExposureCompensation SignedRational
	Width                uint32
	Height               uint32
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

func csvHeader() string {
	return "Make,Model,CreateTime,Iso,FNumber,ExposureTime,FocalLength,FocalLength35,ExpComp,Flash,ExposureProgram,MPix"
}

func (ei *ExifInfo) asCsv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.Make))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.Model))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.CreateTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("\"%d\"", ei.Iso))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.FNumber.ToString()))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.ExposureTime.ToString()))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.FocalLength.ToString()))
	sb.WriteString(fmt.Sprintf("\"%d\"", ei.FocalLength35))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.ExposureCompensation.ToString()))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.ExposureCompensation.ToString()))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.Flash))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.ExposureProgram))
	mpix := float64(ei.Width*ei.Height) / 1000000.0
	sb.WriteString(fmt.Sprintf("\"%.1f\"", mpix))
	return sb.String()
}

type tagValueExtractor = func(tag exif.ExifTag, exifInfo *ExifInfo)

const (
	tagMake                 = 0x010f
	tagModel                = 0x0110
	tagDateTimeOriginal     = 0x9003
	tagIso                  = 0x8827
	tagFNumber              = 0x829d
	tagExposureTime         = 0x829a
	tagFocalLength          = 0x920a
	tagFocalLength35        = 0xa405
	tagFlash                = 0x9209
	tagExposureProgram      = 0x8822
	tagExposureCompensation = 0x9204
	tagMeteringMode         = 0x9207
	tagMaxAperture          = 0x9205
	tagOrientation          = 0x0112
	tagImageWidth           = 0xa002
	tagImageHeight          = 0xa003
	tagExposureMore         = 0xa402
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

var extractors = map[uint16]tagValueExtractor{
	tagMake: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Make = tag.Value.(string)
	},
	tagModel: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Model = tag.Value.(string)
	},
	tagDateTimeOriginal: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		tm, err := exif.ParseExifFullTimestamp(tag.Value.(string))
		if err == nil {
			exifInfo.CreateTime = tm
		} else {
			exifInfo.CreateTime = time.Unix(0, 0)
		}
	},
	tagIso: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Iso = tag.Value.([]uint16)[0]
	},
	tagFNumber: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.FNumber = newRational(tag.Value.([]exif.Rational)[0])
	},
	tagExposureTime: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.ExposureTime = newRational(tag.Value.([]exif.Rational)[0])
	},
	tagFocalLength: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.FocalLength = newRational(tag.Value.([]exif.Rational)[0])
	},
	tagFocalLength35: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.FocalLength35 = tag.Value.([]uint16)[0]
	},
	tagFlash: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		val, ok := exifFlashValues[uint(tag.Value.([]uint16)[0])]
		if ok {
			exifInfo.Flash = val
		}
	},
	tagExposureProgram: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		val, ok := exifExposurePrograms[uint(tag.Value.([]uint16)[0])]
		if ok {
			exifInfo.ExposureProgram = val
		}
	},
	tagExposureCompensation: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.ExposureCompensation = newSignedRational(tag.Value.([]exif.SignedRational)[0])
	},
	tagImageWidth: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Width = tag.Value.([]uint32)[0]
	},
	tagImageHeight: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Height = tag.Value.([]uint32)[0]
	},
}

func retrieveFlatExifData(exifData []byte) (exifTags []exif.ExifTag, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, exifData)
	log.PanicIf(err)

	q := []*exif.Ifd{index.RootIfd}

	exifTags = make([]exif.ExifTag, 0)

	for len(q) > 0 {
		var ifd *exif.Ifd
		ifd, q = q[0], q[1:]

		ti := exif.NewTagIndex()
		for _, ite := range ifd.Entries {
			tagName := ""

			it, err := ti.Get(ifd.IfdPath, ite.TagId)
			if err != nil {
				// If it's a non-standard tag, just leave the name blank.
				if log.Is(err, exif.ErrTagNotFound) != true {
					// log.PanicIf(err)
					tagName = "Unknown"
				}
			} else {
				tagName = it.Name
			}

			value, err := ifd.TagValue(ite)
			if err != nil {
				if err == exif.ErrUnhandledUnknownTypedTag {
					value = exif.UnparseableUnknownTagValuePlaceholder
				} else {
					value = "Unknown"
					// continue
					// log.Panic(err)
				}
			}

			valueBytes, err := ifd.TagValueBytes(ite)
			if err != nil && err != exif.ErrUnhandledUnknownTypedTag {
				//log.Panic(err)
			}

			et := exif.ExifTag{
				IfdPath:      ifd.IfdPath,
				TagId:        ite.TagId,
				TagName:      tagName,
				TagTypeId:    ite.TagType,
				TagTypeName:  exif.TypeNames[ite.TagType],
				Value:        value,
				ValueBytes:   valueBytes,
				ChildIfdPath: ite.ChildIfdPath,
			}

			exifTags = append(exifTags, et)
		}

		for _, childIfd := range ifd.Children {
			q = append(q, childIfd)
		}

		if ifd.NextIfd != nil {
			q = append(q, ifd.NextIfd)
		}
	}

	return exifTags, nil
}

// ExtractExif parses image file with a given path and extracts exif information
func ExtractExif(imageFilePath string) (*ExifInfo, error) {
	rawExif, err := exif.SearchFileAndExtractExif(imageFilePath)
	if err != nil {
		return nil, err
	}

	entries, err := retrieveFlatExifData(rawExif)
	if err != nil {
		if err == exif.ErrTagTypeNotValid || err == exif.ErrTagNotStandard {
			return nil, err
		}
	}

	var exifInfo ExifInfo

	for _, entry := range entries {
		extractor, ok := extractors[entry.TagId]
		if ok {
			extractor(entry, &exifInfo)
		}

	}

	return &exifInfo, nil
}
