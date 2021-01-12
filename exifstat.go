package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var (
	options = &struct {
		FolderPath string `short:"s" long:"src-folder" description:"Path to folder with image files" required:"true"`
		OutputFile string `short:"o" long:"output" description:"Name of the output CSV file" default:"exif-stats.csv"`
		Verbose    bool   `short:"v" long:"verbose" description:"Print more stuff"`
	}{}
)

func csvHeader() string {
	return "Make,Model,CreateTime,Iso,FNumber,ExposureTime,FocalLength,FocalLength35,ExpComp,Flash,ExposureProgram,MPix\n"
}

func (ei *ExifInfo) asCsv() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\"%s\",", strings.TrimSpace(ei.Make)))
	sb.WriteString(fmt.Sprintf("\"%s\",", strings.TrimSpace(ei.Model)))
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.CreateTime))
	sb.WriteString(fmt.Sprintf("\"%d\",", ei.Iso))
	sb.WriteString(fmt.Sprintf("\"%.1f\",", ei.FNumber.AsFloat()))
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.ExposureTime.Normalize().ToString()))
	sb.WriteString(fmt.Sprintf("\"%.1f\",", ei.FocalLength.AsFloat()))
	sb.WriteString(fmt.Sprintf("\"%d\",", ei.FocalLength35))
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.ExposureCompensation.Normalize().ToString()))
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.Flash))
	sb.WriteString(fmt.Sprintf("\"%s\",", ei.ExposureProgram))
	mpix := float64(ei.Width*ei.Height) / 1000000.0
	sb.WriteString(fmt.Sprintf("\"%.1f\"", mpix))
	sb.WriteString("\n")
	return sb.String()
}

func main() {
	_, err := flags.Parse(options)

	if err != nil {
		os.Exit(-1)
	}
	fmt.Printf("Searching for images in '%s'\n", options.FolderPath)

	images, err := ListImages(options.FolderPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d image files\n", len(images))

	out, err := os.Create(options.OutputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.WriteString(csvHeader())

	fmt.Printf("Writing data to '%s'\n", options.OutputFile)
	total := len(images)
	for index, image := range images {
		exif, err := ExtractExif(image)
		if err != nil {
			fmt.Printf("Failed to extract EXIF from '%s': %s\n", image, err)
		} else {
			out.WriteString(exif.asCsv())
			if index%100 == 0 { // update status every 100 images
				fmt.Printf("%.1f%% %s\x1b[0K\r", float64(index)*100.0/float64(total), image)
			}
		}
	}
	fmt.Println("100% Done\x1b[0K")
}
