package main

import (
	"testing"
)

// TestListImages tests ligst images
func TestListImages(t *testing.T) {
	paths := make(chan string)
	err := ListImages("test-data/scan", paths)
	if err != nil {
		t.Errorf("ListImages returned error %s", err)
	}
	var expected = map[string]bool{
		"test-data/scan/DSC_0352.jpg":                       true,
		"test-data/scan/P1020297.JPG":                       true,
		"test-data/scan/_DSC0958.jpg":                       true,
		"test-data/scan/subdir/DSC_3455.JPG":                true,
		"test-data/scan/subdir/P1020630.jpg":                true,
		"test-data/scan/subdir/triple-nested/DSC_9068.jpg":  true,
		"test-data/scan/subdir/triple-nested/P1030129.jpeg": true,
	}

	for image := range paths {
		_, ok := expected[image]
		if !ok {
			t.Errorf("Image '%s' was expected but not found", image)
		}
		delete(expected, image)
	}

	if len(expected) != 0 {
		keys := make([]string, 0, len(expected))
		for k := range expected {
			keys = append(keys, k)
		}
		t.Errorf("Some of the images were expected but not present: %v", keys)
	}
}
