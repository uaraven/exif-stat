package main

import "github.com/uaraven/exif-stat/exif"

const ExifMagic = 0xFFD8

type ExifFileReader interface {
	// ReadExif parses image file with a given path and extracts exif information
	ReadExif() (exifInfo *ExifInfo, err error)
	IsFileSupported() bool
}

func readExifFromFile(f exif.File, fileName string) (*ExifInfo, error) {
	exifInfo := &ExifInfo{
		FileName: fileName,
	}

	tags, err := exif.ReadExifTags(f)
	if err != nil {
		return nil, err
	}

	tagMap := exif.TagsAsMap(tags)

	for path, extractor := range extractors {
		tag, ok := tagMap[path]
		if ok {
			extractor(tag, exifInfo)
		}
	}
	if _, ok := tagMap[tagIso]; !ok { // no standard ISO tag
		if tag, ok := tagMap[tagNikonIso]; ok { // but there is Nikon-specific ISO tag
			extractNikonIso(tag, exifInfo)
		}
	}

	exifInfo = postProcessExif(exifInfo)
	return exifInfo, nil
}

func readTiffFromFile(f exif.File, fileName string) (*ExifInfo, error) {
	exifInfo := &ExifInfo{
		FileName: fileName,
	}

	tags, err := exif.ReadTiffTags(f)
	if err != nil {
		return nil, err
	}

	tagMap := exif.TagsAsMap(tags)

	for path, extractor := range extractors {
		tag, ok := tagMap[path]
		if ok {
			extractor(tag, exifInfo)
		}
	}
	if _, ok := tagMap[tagIso]; !ok { // no standard ISO tag
		if tag, ok := tagMap[tagNikonIso]; ok { // but there is Nikon-specific ISO tag
			extractNikonIso(tag, exifInfo)
		}
	}

	exifInfo = postProcessExif(exifInfo)
	return exifInfo, nil
}
