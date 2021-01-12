package main

var cropFactor = map[string]float64{
	"Canon EOS R6": 1.0,
}

// Consult camera DB to find the crop factor and calculate 35mm equivalent focal length, if missing
func postProcessFocalLength35(exif *ExifInfo) *ExifInfo {
	crop, ok := cropFactor[exif.Model]
	if ok {
		exif.FocalLength35 = uint16(exif.FocalLength.AsFloat() * crop)
	}
	return exif
}

// perform post-processing of exif based on camera knowledge. For example populate 35mm focus length in Canon cameras that do no report it
func postProcessExif(exif *ExifInfo) *ExifInfo {
	if exif.FocalLength35 == 0 {
		return postProcessFocalLength35(exif)
	}

	return exif
}
