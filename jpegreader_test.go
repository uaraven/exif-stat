package main

import (
	"testing"
	"time"

	"github.com/uaraven/exif-stat/utils"
)

func verify(t *testing.T, exif *ExifInfo,
	creationTime string, make string, model string,
	iso uint16, fstop utils.Rational, exposureTime utils.Rational,
	focalLength utils.Rational, focalLength35 uint16,
	exposureComp utils.SignedRational, flash string, exposureProgram string) {

	if exif.CreateTime.Format(time.RFC3339) != creationTime {
		t.Errorf("Invalid creation time: '%s', expected '%s'", exif.CreateTime.Format(time.RFC3339), creationTime)
	}

	if exif.Make != make {
		t.Errorf("Invalid make: '%s', expected '%s'", exif.Make, make)
	}

	if exif.Model != model {
		t.Errorf("Invalid model: '%s', expected '%s'", exif.Model, model)
	}

	if exif.Iso != iso {
		t.Errorf("Invalid ISO: '%d', expected:'%d'", exif.Iso, iso)
	}

	if exif.FNumber.CompareTo(fstop) != 0 {
		t.Errorf("Invalid F-Number: '%s', expected:'%s'", exif.FNumber.ToString(), fstop.ToString())
	}

	if exif.ExposureTime.CompareTo(exposureTime) != 0 {
		t.Errorf("Invalid ExposureTime: '%s', expected:'%s'", exif.ExposureTime.ToString(), exposureTime.ToString())
	}

	if exif.FocalLength.CompareTo(focalLength) != 0 {
		t.Errorf("Invalid Focal Length: '%s', expected:'%s'", exif.FocalLength.ToString(), focalLength.ToString())
	}

	if exif.FocalLength35 != focalLength35 {
		t.Errorf("Invalid Focal Length 35mm: '%d', expected:'%d'", exif.FocalLength35, focalLength35)
	}

	if exif.ExposureCompensation.CompareTo(exposureComp) != 0 {
		t.Errorf("Invalid Exposure Compensation: '%s', expected:'%s'", exif.ExposureCompensation.ToString(), exposureComp.ToString())
	}

	if exif.Flash != flash {
		t.Errorf("Invalid Flash: '%s', expected:'%s'", exif.Flash, flash)
	}

	if exif.ExposureProgram != exposureProgram {
		t.Errorf("Invalid ExposureProgram: '%s', expected:'%s'", exif.ExposureProgram, exposureProgram)
	}

	// CreateTime           string
}

func TestExtractExifNikonD50(t *testing.T) {
	exifInfo, err := ExtractExif("test-data/scan/DSC_0352.jpg")

	if err != nil {
		t.Errorf("ExtractExif returned error: %v", err)
	}

	verify(t, exifInfo,
		"2006-04-16T15:11:33Z",
		"NIKON CORPORATION", "NIKON D50", 0, utils.NewRational(10, 1), utils.NewRational(1, 400), utils.NewRational(18, 1), 27, utils.NewSignedRational(-1, 3), "No Flash", "Program AE")

}

func TestExtractExifGX1(t *testing.T) {
	exifInfo, err := ExtractExif("test-data/scan/subdir/P1020630.jpg")

	if err != nil {
		t.Errorf("ExtractExif returned error: %v", err)
	}

	verify(t, exifInfo,
		"2019-06-08T08:19:37Z",
		"Panasonic", "DMC-GX85", 250, utils.NewRational(8, 1), utils.NewRational(1, 800), utils.NewRational(140, 1), 280, utils.NewSignedRational(-33, 100), "Off, Did not fire", "Aperture-priority AE")

}
