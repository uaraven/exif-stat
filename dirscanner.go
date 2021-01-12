package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uaraven/exif-stat/utils"
)

var supportedFiles = map[string]bool{
	".jpg":  true,
	".jpeg": true,
}

func isSupportedFile(path string) bool {
	_, ok := supportedFiles[strings.ToLower(filepath.Ext(path))]
	return ok
}

// ListImages lists all the supported images in given path. includes images in subdirectories
func ListImages(path string) ([]string, error) {
	var imageFiles []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Printf("%s%s\r", utils.Shorten(path), utils.ClearLine)
		}
		if !info.IsDir() && isSupportedFile(path) && filepath.Base(path)[0] != '.' {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})
	fmt.Printf("\r%s\n", utils.ClearLine)
	if err != nil {
		return imageFiles, err
	}
	return imageFiles, nil
}
