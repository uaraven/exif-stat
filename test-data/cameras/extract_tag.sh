#!/bin/sh

# Run this script and pass a path to an image file to it 
# It will spit out all the EXIF tag values that's neede by exif-stat

exiftool -H -G1 -IFD0:Make -IFD0:Model -ExifIFD:ISO -Nikon:ISO -ExifIFD:FNumber -ExifIFD:ExposureTime\
 -ExifIFD:ExposureProgram -ExifIFD:ExifImageWidth -ExifIFD:ExifImageHeight -ExifIFD:FocalLength -ExifIFD:FocalLengthIn35mmFormat\
 -ExifIFD:ExposureCompensation -ExifIFD:Flash -ExifIFD:CreateDate $1