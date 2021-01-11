package main

import (
	"testing"
)

// TestListImages tests ligst images
func TestListImages(t *testing.T) {
	images, err := ListImages("test-data")
	if err != nil {
		t.Errorf("ListImages returned error %s", err)
	}
	var expected = map[string]bool{
		"test-data/DSC_0352.jpg":                       true,
		"test-data/P1020297.JPG":                       true,
		"test-data/_DSC0958.jpg":                       true,
		"test-data/subdir/DSC_3455.JPG":                true,
		"test-data/subdir/P1020630.jpg":                true,
		"test-data/subdir/triple-nested/DSC_9068.jpg":  true,
		"test-data/subdir/triple-nested/P1030129.jpeg": true,
	}

	for _, image := range images {
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
