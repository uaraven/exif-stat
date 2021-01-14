package exif

import (
	"testing"
)

func TestReadingExifFileMMap(t *testing.T) {
	file, err := OpenExifFileMMap("../test-data/scan/DSC_0352.jpg")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { file.Close() }()

	ifds, err := ReadExifTags(file)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	if len(ifds) != 88 {
		t.Fatalf("Failed to find all 88 tags, found: %d", len(ifds))
	}

	tagMap := TagsAsMap(ifds)

	if tag, ok := tagMap[nikonIso]; ok {
		if tag.Value.([]uint16)[1] != 400 {
			t.Fatalf("Invalid Nikon-specific ISO expected:400, actual: %v", tag.Value)
		}
	} else {
		t.Fatalf("Failed to find nikon-specific ISO tag")
	}

	if tag, ok := tagMap[focalLength35]; ok {
		if tag.Value.([]uint16)[0] != 27 {
			t.Fatalf("Invalid 35mm focal length. Expected:27, actual: %v", tag.Value)
		}
	} else {
		t.Fatalf("Failed to find 35mm focal length tag")
	}
}

func TestReadingFileWithoutExifMMap(t *testing.T) {
	file, err := OpenExifFileMMap("../test-data/fail/stack.jpg")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { file.Close() }()

	ifds, err := ReadExifTags(file)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	if len(ifds) != 0 {
		t.Fatalf("Should have found no tags")
	}
}

func TestReadingNonExistingFileMMap(t *testing.T) {
	file, err := OpenExifFileMMap("../test-data/fail/no-file.jpg")
	if err == nil {
		t.Fatalf("Should have failed to open file")
		return
	}
	if file != nil {
		file.Close()
	}
}

func TestReadingEmptyFileMMap(t *testing.T) {
	file, err := OpenExifFileMMap("../test-data/fail/empty.jpg")
	if err == nil {
		t.Fatalf("Should have failed to open file")
		return
	}
	if file != nil {
		t.Fatalf("Should have not returned the file")
		file.Close()
	}

}
