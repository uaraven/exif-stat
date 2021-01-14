package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

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
		exifInfo, err := ExtractExif(filepath, false)
		if err != nil {
			t.Errorf("Failed to read exif from %s: %v", filepath, err)
		}
		compareExifMaps(t, camera.Camera, camera.Exif, exifInfo.toMap())
	}
}
