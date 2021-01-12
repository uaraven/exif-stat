package exif

import (
	"fmt"
	"strconv"
	"strings"
)

type marker struct {
	Marker uint16
	Size   uint16
	Offset int64
}

// IfdEntry is an exif tag
type ifdEntry struct {
	// Index of IFD
	IfdIndex int
	// Tag ID
	TagID uint16
	// Data Type
	DataType uint16
	// Number of components
	ComponentCount uint32
	// Raw data value (can be data or offset)
	Data uint32
	// Parsed and converted data
	Value interface{}
	// Actual data as slice of bytes
	ValueBytes []byte
}

// Ifd represents image format descriptor
type ifd struct {
	Index      int
	EntryCount uint16
	IfdEntries []ifdEntry
}

// ToString returns a string representation of IfdEntry
func (ie ifdEntry) ToString() string {
	return fmt.Sprintf("ID=%x Value=%v Bytes=%v", ie.TagID, ie.Value, ie.ValueBytes)
}

// Tag is a simplified representation of an Exif Tag
type Tag struct {
	ID       uint16
	IDPath   []uint16
	DataType int
	Value    interface{}
	RawData  []byte
}

// PathName creates path-line name from tag Id and parent ids
func (tag Tag) PathName() string {
	var sb strings.Builder
	for _, p := range tag.IDPath {
		sb.WriteString(fmt.Sprintf("%04x", p))
		sb.WriteString("/")
	}
	sb.WriteString(fmt.Sprintf("%04x", tag.ID))
	return sb.String()
}

// ToString returns a string representation of IfdEntry
func (tag Tag) ToString() string {
	return fmt.Sprintf("Path=%s ID=%x Value=%v", tag.PathName(), tag.ID, tag.Value)
}

// Tags is a slice of all Exif tags
type Tags []Tag

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
	factor := gcd(r.Numerator, r.Denominator)
	if factor != 0 {
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
	if r1.Denominator == 1 {
		return strconv.FormatUint(uint64(r1.Numerator), 10)
	}
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

// Normalize normalizes a rational
func (r SignedRational) Normalize() SignedRational {
	if r.Numerator == 0 || r.Denominator == 0 {
		return r
	}
	n := abs(r.Numerator)
	commonDivisor := gcds(n, r.Denominator)
	// if r.Denominator%n == 0 || n%r.Denominator == 0 {
	if commonDivisor != 0 {
		return SignedRational{
			Numerator:   r.Numerator / commonDivisor,
			Denominator: r.Denominator / commonDivisor,
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
	if r1.Denominator == 1 {
		return strconv.FormatInt(int64(r1.Numerator), 10)
	}
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
