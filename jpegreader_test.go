package main

import (
	"testing"
)

func TestExtractExif(t *testing.T) {
	exifInfo, err := ExtractExif("test-data/DSC_0352.jpg")

	if err != nil {
		t.Errorf("ExtractExif returned error: %v", err)
	}

	if exifInfo.Make != "NIKON CORPORATION" {
		t.Errorf("Invalid make: '%s', expected 'NIKON CORPORATION'", exifInfo.Make)
	}

	if exifInfo.Model != "NIKON D50" {
		t.Errorf("Invalid model: '%s', expected 'NIKON D50'", exifInfo.Model)
	}
}
