package utils

import (
	"fmt"

	"github.com/dsoprea/go-exif"
)

// ClearLine contains ASCII escape code to clear to the end of line
const ClearLine = "\x1b[0K"

func gcd(a, b uint32) uint32 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func gcds(a, b int32) int32 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// Rational represents a rational value expressed as Numerator/Denominator
type Rational struct {
	Numerator   uint32
	Denominator uint32
}

// Normalize normalizes a rational
func (r Rational) Normalize() Rational {
	if r.Numerator == 0 || r.Denominator == 0 {
		return r
	}
	if r.Denominator%r.Numerator == 0 || r.Numerator%r.Denominator == 0 {
		factor := gcd(r.Numerator, r.Denominator)
		return Rational{
			Numerator:   r.Numerator / factor,
			Denominator: r.Denominator / factor,
		}
	}
	return r
}

// ToString converts a rational value to string
func (r Rational) ToString() string {
	if r.Numerator == 0 {
		return "0"
	}
	r1 := r.Normalize()
	if r1.Numerator > r.Denominator {
		var value = r1.Numerator / r1.Denominator
		r1 = Rational{
			Numerator:   r1.Numerator - value*r1.Denominator,
			Denominator: r1.Denominator,
		}
		if r1.Numerator == 0 {
			return fmt.Sprintf("%d", value)
		}
		return fmt.Sprintf("%d %d/%d", value, r1.Numerator, r1.Denominator)
	}
	return fmt.Sprintf("%d/%d", r1.Numerator, r1.Denominator)
}

// AsFloat converts a rational value to float
func (r Rational) AsFloat() float64 {
	return float64(r.Numerator) / float64(r.Denominator)
}

// CompareTo compares a rational to another rational and return 1 if this rational is larger than other,
// -1 if it is smaller and 0 if they are equal
func (r Rational) CompareTo(other Rational) int {
	r1 := r.Numerator * other.Denominator
	r2 := other.Numerator * r.Denominator

	if r1 > r2 {
		return 1
	} else if r1 < r2 {
		return -1
	} else {
		return 0
	}
}

// NewRational creates a new rational from numerator and denominator
func NewRational(numerator uint32, denominator uint32) Rational {
	return Rational{
		Numerator:   numerator,
		Denominator: denominator}
}

func newRational(r exif.Rational) Rational {
	return Rational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

// SignedRational represents a rational value expressed as Numerator/Denominator
type SignedRational struct {
	Numerator   int32
	Denominator int32
}

func abs(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// NewSignedRational creates a new rational from numerator and denominator
func NewSignedRational(numerator int32, denominator int32) SignedRational {
	return SignedRational{
		Numerator:   numerator,
		Denominator: denominator}
}

func newSignedRational(r exif.SignedRational) SignedRational {
	return SignedRational{
		Numerator:   r.Numerator,
		Denominator: r.Denominator,
	}
}

// Normalize normalizes a rational
func (r SignedRational) Normalize() SignedRational {
	if r.Numerator == 0 || r.Denominator == 0 {
		return r
	}
	n := abs(r.Numerator)
	if r.Denominator%n == 0 || n%r.Denominator == 0 {
		factor := abs(gcds(n, r.Denominator))
		return SignedRational{
			Numerator:   r.Numerator / factor,
			Denominator: r.Denominator / factor,
		}
	}
	return r
}

// ToString converts a signed rational value to string
// ToString converts a rational value to string
func (r SignedRational) ToString() string {
	if r.Numerator == 0 {
		return "0"
	}
	r1 := r.Normalize()
	if abs(r1.Numerator) > abs(r.Denominator) {
		var value = r1.Numerator / r1.Denominator
		r1 = SignedRational{
			Numerator:   r1.Numerator - value*r1.Denominator,
			Denominator: r1.Denominator,
		}
		if r1.Numerator == 0 {
			return fmt.Sprintf("%d", value)
		}
		return fmt.Sprintf("%d %d/%d", value, r1.Numerator, r1.Denominator)
	}
	return fmt.Sprintf("%d/%d", r1.Numerator, r1.Denominator)
}

// AsFloat converts a signed rational value to float
func (r SignedRational) AsFloat() float64 {
	return float64(r.Numerator) / float64(r.Denominator)
}

// CompareTo compares a signed rational to another signed rational and return 1 if this rational is larger than other,
// -1 if it is smaller and 0 if they are equal
func (r SignedRational) CompareTo(other SignedRational) int {
	r1 := r.Numerator * other.Denominator
	r2 := other.Numerator * r.Denominator

	if r1 > r2 {
		return 1
	} else if r1 < r2 {
		return -1
	} else {
		return 0
	}
}

// Shorten the string if it's longer than predefined number of character by replacing part of it with ellipsis
func Shorten(text string) string {
	if len(text) > 60 {
		return "…" + text[len(text)-61:]
	}
	return text
}
