package main

import (
	"fmt"
	"strings"

	"github.com/dsoprea/go-exif"
	log "github.com/dsoprea/go-logging"
)

func gcd(a, b uint32) uint32 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func gcds(a, b int32) int32 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// Rational represents a rational value expressed as Numerator/Denominator
type Rational struct {
	Numerator   uint32
	Denominator uint32
}

// Normalize normalizes a rational
func (r *Rational) Normalize() Rational {
	if r.Denominator%r.Numerator == 0 || r.Numerator%r.Denominator == 0 {
		factor := gcd(r.Numerator, r.Denominator)
		return Rational{
			Numerator:   r.Numerator / factor,
			Denominator: r.Denominator / factor,
		}
	}
	return *r
}

// ToString converts a rational value to string
func (r *Rational) ToString() string {
	r1 := r.Normalize()
	if r1.Numerator > r.Denominator {
		var value = r1.Numerator / r1.Denominator
		r1 = Rational{
			Numerator:   r1.Numerator - value*r1.Denominator,
			Denominator: r1.Denominator,
		}
		if r1.Numerator == 0 {
			return fmt.Sprintf("%d", value)
		}
		return fmt.Sprintf("%d %d/%d", value, r1.Numerator, r1.Denominator)
	}
	return fmt.Sprintf("%d/%d", r1.Numerator, r1.Denominator)
}

// AsFloat converts a rational value to float
func (r *Rational) AsFloat() float64 {
	return float64(r.Numerator) / float64(r.Denominator)
}

// CompareTo compares a rational to another rational and return 1 if this rational is larger than other,
// -1 if it is smaller and 0 if they are equal
func (r *Rational) CompareTo(other Rational) int {
	r1 := r.Numerator * other.Denominator
	r2 := other.Numerator * r.Denominator

	if r1 > r2 {
		return 1
	} else if r1 < r2 {
		return -1
	} else {
		return 0
	}
}

// NewRational creates a new rational from numerator and denominator
func NewRational(numerator uint32, denominator uint32) Rational {
	return Rational{
		Numerator:   numerator,
		Denominator: denominator}
}

func newRational(r exif.Rational) Rational {
	return Rational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

// SignedRational represents a rational value expressed as Numerator/Denominator
type SignedRational struct {
	Numerator   int32
	Denominator int32
}

func abs(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// NewSignedRational creates a new rational from numerator and denominator
func NewSignedRational(numerator int32, denominator int32) SignedRational {
	return SignedRational{
		Numerator:   numerator,
		Denominator: denominator}
}

func newSignedRational(r exif.SignedRational) SignedRational {
	return SignedRational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

// Normalize normalizes a rational
func (r *SignedRational) Normalize() SignedRational {
	n := abs(r.Numerator)
	if r.Denominator%n == 0 || n%r.Denominator == 0 {
		factor := abs(gcds(n, r.Denominator))
		return SignedRational{
			Numerator:   r.Numerator / factor,
			Denominator: r.Denominator / factor,
		}
	}
	return *r
}

// ToString converts a signed rational value to string
// ToString converts a rational value to string
func (r *SignedRational) ToString() string {
	r1 := r.Normalize()
	if abs(r1.Numerator) > abs(r.Denominator) {
		var value = r1.Numerator / r1.Denominator
		r1 = SignedRational{
			Numerator:   r1.Numerator - value*r1.Denominator,
			Denominator: r1.Denominator,
		}
		if r1.Numerator == 0 {
			return fmt.Sprintf("%d", value)
		}
		return fmt.Sprintf("%d %d/%d", value, r1.Numerator, r1.Denominator)
	}
	return fmt.Sprintf("%d/%d", r1.Numerator, r1.Denominator)
}

// AsFloat converts a signed rational value to float
func (r *SignedRational) AsFloat() float64 {
	return float64(r.Numerator) / float64(r.Denominator)
}

// CompareTo compares a signed rational to another signed rational and return 1 if this rational is larger than other,
// -1 if it is smaller and 0 if they are equal
func (r *SignedRational) CompareTo(other SignedRational) int {
	r1 := r.Numerator * other.Denominator
	r2 := other.Numerator * r.Denominator

	if r1 > r2 {
		return 1
	} else if r1 < r2 {
		return -1
	} else {
		return 0
	}
}

// ExifInfo contains values of all the exif tag of interest
type ExifInfo struct {
	Make                 string
	Model                string
	CreateTime           string
	Iso                  uint16
	FNumber              Rational
	ExposureTime         Rational
	FocalLength          Rational
	FocalLength35        uint16
	Flash                string
	ExposureProgram      uint16
	ExposureCompensation SignedRational
	Megapixels           uint32
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
	sb.WriteString(fmt.Sprintf("Exposure program: %d\n", ei.ExposureProgram))
	return sb.String()
}

func csvHeader() string {
	return "Make,Model,CreateTime,Iso,FNumber,ShutterSpeed,FocalLength,FocalLength35,LensMake,LensModel,Flash,ExposureProgram"
}

func (ei *ExifInfo) asCsv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.Make))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.Model))
	sb.WriteString(fmt.Sprintf("\"%s\"", ei.CreateTime))
	sb.WriteString(fmt.Sprintf("\"%d\"", ei.Iso))
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

var extractors = map[uint16]tagValueExtractor{
	tagMake: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Make = tag.Value.(string)
	},
	tagModel: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.Model = tag.Value.(string)
	},
	tagDateTimeOriginal: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.CreateTime = tag.Value.(string)
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
		exifInfo.ExposureProgram = tag.Value.([]uint16)[0]
	},
	tagExposureCompensation: func(tag exif.ExifTag, exifInfo *ExifInfo) {
		exifInfo.ExposureCompensation = newSignedRational(tag.Value.([]exif.SignedRational)[0])
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
