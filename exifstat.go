package main

import (
	"fmt"
	"os"
	"strings"

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
		WriteFileName bool   `short:"f" long:"file-name" description:"Include file name in the output"`
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

	logger.Verbose(0, fmt.Sprintf("Searching for images in '%s'", options.Args.FolderPath))

	images, err := ListImages(options.Args.FolderPath)
	if err != nil {
		panic(err)
	}
	logger.Verbose(0, fmt.Sprintf("Found %d image files", len(images)))

	out, err := os.Create(options.OutputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.WriteString(csvHeader())

	logger.Info(fmt.Sprintf("Writing data to '%s'", options.OutputFile))
	total := len(images)
	if total == 0 {
		total = 1 // to avoid div by zero later
	}
	for index, image := range images {
		exif, err := ExtractExif(image)
		if err != nil {
			logger.Verbose(1, fmt.Sprintf("\nFailed to extract EXIF from '%s': %s", image, err))
		} else {
			if exif.isValidExif() {
				out.WriteString(exif.asCsv())
			}
			if index%100 == 0 { // update status every 100 images
				fmt.Printf("%.1f%% %s\x1b[0K\r", float64(index)*100.0/float64(total), image)
			}
		}
	}
	fmt.Println("100% Done\x1b[0K")
}
