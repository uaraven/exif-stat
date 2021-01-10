package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	options = &struct {
		FolderPath string `short:"s" long:"src-folder" description:"Path to folder with image files" required:"true"`
		OutputFile string `short:"o" long:"output" description:"Name of the output CSV file" default:"exif-stats.csv"`
		Verbose    bool   `short:"v" long:"verbose" description:"Print more stuff"`
	}{}
)

func main() {
	_, err := flags.Parse(options)

	if err != nil {
		os.Exit(-1)
	}
	fmt.Printf("Parsing images in %s folder", options.FolderPath)

}
