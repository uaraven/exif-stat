# Exif-stat

Extracts some metadata from image files and stores it in CSV file for additional analysis.

## What it does?

It scans all JPEG files in a given folder (including all subfolders recursively) and tries to read EXIF data. Some of that data is then written to a file.

It is expected that scanned JPEG files are pictures from digital cameras.

Following information is extracted:

 - Camera Make
 - Camera Model
 - Photo Creation time
 - Iso number
 - Aperture F-number,
 - Exposure time
 - Focal length
 - Equivalent focal length for 35mm
 - Exposure compensation
 - Flash
 - ExposureProgram (PASM, etc.)
 
## Supported EXIF data

Only standard EXIF tags are parsed. Of the vendor-specific tags only some Nikon tags are parsed to retrieve ISO value when it is not present in Exif IFD.

## Tested cameras

| Make      | Model    | Notes                                                |
|:---------:|:---------|:-----------------------------------------------------|
| Nikon     | D50      | No ISO in Exif IFD, retrieved from Nikon maker notes |
| Nikon     | D90      |                                                      |
| Nikon     | D7000    |                                                      |
| Nikon     | D750     | Exif IFD does not contain image width or height tags |
| Panasonic | DMC-GX1  |                                                      |
| Panasonic | DMC-GX85 |                                                      |
| Fujifilm  | X-S10    |                                                      |

