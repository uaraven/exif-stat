package main

import (
	"testing"
)

func TestRational(t *testing.T) {
	r := Rational{2, 6}
	rn := r.Normalize()
	if rn.Numerator != 1 || rn.Denominator != 3 {
		t.Errorf("Normalization failed, expected: 1/3, actual: %s", rn.ToString())
	}
	if rn.ToString() != "1/3" {
		t.Errorf("ToString failed, expected: 1/3, actual: %s", rn.ToString())
	}

	rn = Rational{6, 2}
	rn = rn.Normalize()
	if rn.Numerator != 3 || rn.Denominator != 1 {
		t.Errorf("Normalization failed, expected: 3/1, actual: %s", rn.ToString())
	}

	if rn.ToString() != "3" {
		t.Errorf("ToString failed, expected: 3, actual: %s", rn.ToString())
	}

	r1 := Rational{1, 3}
	r2 := Rational{1, 4}
	if r1.CompareTo(r2) <= 0 {
		t.Errorf("CompareTo failed, %s must be > %s", r1.ToString(), r2.ToString())
	}

	r1 = Rational{1, 2}
	r2 = Rational{4, 8}
	if r1.CompareTo(r2) != 0 {
		t.Errorf("CompareTo failed, %s must be == %s", r1.ToString(), r2.ToString())
	}

	r1 = Rational{1, 5}
	r2 = Rational{2, 5}
	if r1.CompareTo(r2) >= 0 {
		t.Errorf("CompareTo failed, %s must be < %s", r1.ToString(), r2.ToString())
	}
}

func TestSignedRational(t *testing.T) {
	r := SignedRational{-2, 6}
	rn := r.Normalize()
	if rn.Numerator != -1 || rn.Denominator != 3 {
		t.Errorf("Normalization failed, expected: -1/3, actual: %s", rn.ToString())
	}
	if rn.ToString() != "-1/3" {
		t.Errorf("ToString failed, expected: -1/3, actual: %s", rn.ToString())
	}

	rn = SignedRational{-6, 2}
	rn = rn.Normalize()
	if rn.Numerator != -3 || rn.Denominator != 1 {
		t.Errorf("Normalization failed, expected: -3/1, actual: %s", rn.ToString())
	}

	if rn.ToString() != "-3" {
		t.Errorf("ToString failed, expected: -3, actual: %s", rn.ToString())
	}

	r1 := SignedRational{-1, 3}
	r2 := SignedRational{-1, 4}
	if r1.CompareTo(r2) >= 0 {
		t.Errorf("CompareTo failed, %s must be > %s", r1.ToString(), r2.ToString())
	}

	r1 = SignedRational{-1, 2}
	r2 = SignedRational{-4, 8}
	if r1.CompareTo(r2) != 0 {
		t.Errorf("CompareTo failed, %s must be == %s", r1.ToString(), r2.ToString())
	}

	r1 = SignedRational{-1, 5}
	r2 = SignedRational{-2, 5}
	if r1.CompareTo(r2) <= 0 {
		t.Errorf("CompareTo failed, %s must be < %s", r1.ToString(), r2.ToString())
	}
}

func verify(t *testing.T, exif *ExifInfo, make string, model string, fstop Rational) {
	if exif.Make != make {
		t.Errorf("Invalid make: '%s', expected '%s'", exif.Make, make)
	}

	if exif.Model != model {
		t.Errorf("Invalid model: '%s', expected '%s'", exif.Model, model)
	}
}

func TestExtractExifNikonD50(t *testing.T) {
	exifInfo, err := ExtractExif("test-data/DSC_0352.jpg")

	if err != nil {
		t.Errorf("ExtractExif returned error: %v", err)
	}

	verify(t, exifInfo, "NIKON CORPORATION", "NIKON D50", NewRational(10, 1))

}
