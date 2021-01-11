package exif

import (
	"fmt"
	"testing"
)

func TestReadingExifFile(t *testing.T) {
	file, err := OpenExifFile("../test-data/DSC_0352.jpg")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { file.Close() }()

	ifds, err := ReadExifTags(file)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	fmt.Println(ifds)
}
