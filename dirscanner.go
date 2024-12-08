package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ryanuber/go-glob"
	"github.com/uaraven/exif-stat/utils"
)

var supportedFiles = map[string]bool{
	".jpg":  true,
	".jpeg": true,
}

func isSupportedFile(path string, mask []string) bool {
	var ok bool
	if mask == nil {
		_, ok = supportedFiles[strings.ToLower(filepath.Ext(path))]
	} else {
		for _, m := range mask {
			ok = glob.Glob(m, strings.ToLower(path))
			if ok {
				return ok
			}
		}
	}
	return ok
}

// ListImages lists all the supported images in given path. includes images in subdirectories
func ListImages(path string, mask string, wg *sync.WaitGroup, paths chan string) {
	defer close(paths)
	defer wg.Done()

	var includeMask []string

	if mask != "" {
		includeMask = strings.Split(strings.ToLower(mask), ",")
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Printf("%s%s\r", utils.Shorten(path), utils.ClearLine)
		} else if !info.IsDir() && isSupportedFile(path, includeMask) && filepath.Base(path)[0] != '.' {
			paths <- path
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
