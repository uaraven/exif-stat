package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
func ListImages(path string, wg *sync.WaitGroup, paths chan string) {
	defer close(paths)
	defer wg.Done()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Printf("%s%s\r", utils.Shorten(path), utils.ClearLine)
		} else if !info.IsDir() && isSupportedFile(path) && filepath.Base(path)[0] != '.' {
			paths <- path
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
