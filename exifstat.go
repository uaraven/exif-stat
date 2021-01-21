package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/uaraven/exif-stat/logger"
)

var (
	options = &struct {
		Args struct {
			FolderPath string
		} `positional-args:"yes" positional-arg-name:"folder-path" description:"Path to folder with image files" required:"yes"`
		OutputFile    string `short:"o" long:"output" description:"Name of the output CSV file" default:"exif-stats.csv"`
		Verbose       bool   `short:"v" long:"verbose" description:"Output more informationm, including warnings"`
		FastFile      bool   `long:"fast-io" description:"Use memory-mapped io. May be unstable with network paths"`
		WriteFileName bool   `short:"f" long:"file-name" description:"Include file name in the output"`
		ExtendFlash   bool   `long:"extend-flash" description:"Detailed flash status"`
	}{}
)

func csvHeader() string {
	var sb strings.Builder
	sb.WriteString("Make")
	sb.WriteString(",Model")
	sb.WriteString(",CreateTime")
	sb.WriteString(",Iso")
	sb.WriteString(",FNumber")
	sb.WriteString(",ExposureTime")
	sb.WriteString(",FocalLength")
	sb.WriteString(",FocalLength35")
	sb.WriteString(",ExpComp")
	sb.WriteString(",Flash")
	sb.WriteString(",ExposureProgram")
	sb.WriteString(",MPix")
	if options.WriteFileName {
		sb.WriteString(",FileName")
	}
	sb.WriteString("\n")
	return sb.String()
}

func (ei *ExifInfo) asCsv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\"%s\"", strings.TrimSpace(ei.Make)))
	sb.WriteString(fmt.Sprintf(",\"%s\"", strings.TrimSpace(ei.Model)))
	sb.WriteString(fmt.Sprintf(",\"%s\"", ei.CreateTime))
	sb.WriteString(fmt.Sprintf(",\"%d\"", ei.Iso))
	sb.WriteString(fmt.Sprintf(",\"%.1f\"", ei.FNumber.AsFloat()))
	sb.WriteString(fmt.Sprintf(",\"%s\"", ei.ExposureTime.Normalize().ToString()))
	sb.WriteString(fmt.Sprintf(",\"%.1f\"", ei.FocalLength.AsFloat()))
	sb.WriteString(fmt.Sprintf(",\"%d\"", ei.FocalLength35))
	sb.WriteString(fmt.Sprintf(",\"%s\"", ei.ExposureCompensation.Normalize().ToString()))
	sb.WriteString(fmt.Sprintf(",\"%s\"", ei.Flash))
	sb.WriteString(fmt.Sprintf(",\"%s\"", ei.ExposureProgram))
	mpix := float64(ei.Width*ei.Height) / 1000000.0
	sb.WriteString(fmt.Sprintf(",\"%.1f\"", mpix))
	if options.WriteFileName {
		sb.WriteString(fmt.Sprintf(",\"%s\"", ei.FileName))
	}
	sb.WriteString("\n")
	return sb.String()
}

func (ei *ExifInfo) isValidExif() bool {
	return len(ei.Make) > 0 && len(ei.Model) > 0 && ei.FNumber.Denominator != 0 && ei.ExposureTime.Denominator != 0 && ei.FocalLength.Denominator != 0 && len(ei.CreateTime) > 0
}

func parseExif(wg *sync.WaitGroup, paths chan string, exifs chan *ExifInfo) {
	defer close(exifs)
	defer wg.Done()
	for path := range paths {
		exif, err := ExtractExif(path, options.FastFile)
		if err == nil {
			exifs <- exif
		} else {
			logger.Verbose(1, fmt.Sprintf("\nFailed to extract EXIF from '%s': %s", path, err))
		}
	}
}

func writeCsv(wg *sync.WaitGroup, exifs chan *ExifInfo) {
	defer wg.Done()
	out, err := os.Create(options.OutputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.WriteString(csvHeader())
	for exif := range exifs {
		if exif.isValidExif() {
			out.WriteString(exif.asCsv())
		}
	}
}

func main() {
	_, err := flags.Parse(options)

	if err != nil {
		os.Exit(-1)
	}
	if options.Verbose {
		logger.SetVerbosityLevel(1)
	} else {
		logger.SetVerbosityLevel(0)
	}

	var wg sync.WaitGroup
	paths := make(chan string)
	exifs := make(chan *ExifInfo)

	wg.Add(3)
	go ListImages(options.Args.FolderPath, &wg, paths)
	go parseExif(&wg, paths, exifs)
	go writeCsv(&wg, exifs)

	wg.Wait()

	fmt.Println("\n100% Done\x1b[0K")
}
