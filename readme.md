# Exif-stat

Extracts some metadata from iamge files and stores it in CSV file for additional analysis.

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
 - Exposure compensation in EV
 - Flash
 - ExposureProgram (PASM, etc.)
 
## Limitations

Only standard EXIF tags are parsed, no vendor-specific tags (yet)

