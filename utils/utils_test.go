package utils

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

	s := Rational{0, 2}.ToString()
	if s != "0" {
		t.Errorf("ToString failed, expected: 0, actual: %s", s)
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

	s := SignedRational{-0, 2}.ToString()
	if s != "0" {
		t.Errorf("ToString failed, expected: 0, actual: %s", s)
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
