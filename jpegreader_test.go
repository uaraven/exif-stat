package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/uaraven/exif-stat/utils"
)

func verify(t *testing.T, exif *ExifInfo,
	creationTime string, make string, model string,
	iso uint16, fstop utils.Rational, exposureTime utils.Rational,
	focalLength utils.Rational, focalLength35 uint16,
	exposureComp utils.SignedRational, flash string, exposureProgram string) {

	if exif.CreateTime != creationTime {
		t.Errorf("Invalid creation time: '%s', expected '%s'", exif.CreateTime, creationTime)
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

}

type CameraTest struct {
	Camera string                 `json:"Camera"`
	Image  string                 `json:"Image"`
	Exif   map[string]interface{} `json:"Exif"`
}

func compareExifMaps(t *testing.T, camera string, expected map[string]interface{}, actual map[string]interface{}) {
	for k, ve := range expected {
		va, ok := actual[k]
		if !ok {
			t.Fatalf("Camera %s. Actual exif does not contain %s", camera, k)
		}
		if fmt.Sprintf("%v", va) != fmt.Sprintf("%v", ve) {
			t.Fatalf("Camera '%s', Tag '%s' Actual value '%v' != '%v'", camera, k, va, ve)
		}
	}
}

func TestCameras(t *testing.T) {
	cameraJSON, err := ioutil.ReadFile("test-data/cameras/cameras.json")
	if err != nil {
		t.Logf("Failed to load cameras.json")
		t.FailNow()
	}
	var cameras []CameraTest
	err = json.Unmarshal([]byte(cameraJSON), &cameras)
	if err != nil {
		t.Logf("Failed to parse cameras.json")
		t.FailNow()
	}
	for _, camera := range cameras {
		filepath := "test-data/cameras/" + camera.Image
		exifInfo, err := ExtractExif(filepath)
		if err != nil {
			t.Errorf("Failed to read exif from %s: %v", filepath, err)
		}
		compareExifMaps(t, camera.Camera, camera.Exif, exifInfo.toMap())
	}
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
