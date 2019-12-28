// Copyright 2019 Montgomery Edwards⁴⁴⁸ and Faye Amacker
//
// Special thanks to Kathryn Long for her Rust implementation
// of float16 at github.com/starkat99/half-rs (MIT license)

package float16

import (
	"math"
	"strconv"
)

// Float16 represents IEEE 754 half-precision floating-point numbers (binary16).
type Float16 uint16

// Frombits returns a Float16 value by casting the specified u16 to Float16.
func Frombits(u16 uint16) Float16 {
	return Float16(u16)
}

// Fromfloat32 returns a Float16 value converted from f32. Conversion uses
// IEEE default rounding (nearest int, with ties to even).
func Fromfloat32(f32 float32) Float16 {
	return Float16(f32bitsToF16bits(math.Float32bits(f32)))
}

// NaN returns a Float16 with Not-a-Number (NaN) value.
func NaN() Float16 {
	return Float16(0x7c00 | 0x03ff)
}

// Inf returns a Float16 with an infinity value with the specified sign.
// A sign >= returns positive infinity.
// A sign < 0 returns negative infinity.
func Inf(sign int) Float16 {
	if sign >= 0 {
		return Float16(0x7c00)
	}
	return Float16(0x8000 | 0x7c00)
}

// Float32 simply wraps Float32.  It returns a float32
// converted from f (Float16). This is a lossless conversion.
func (f Float16) Float32() float32 {
	u32 := f16bitsToF32bits(uint16(f))
	return math.Float32frombits(u32)
}

// IsNaN reports whether f is an IEEE 754 “not-a-number” value.
func (f Float16) IsNaN() bool {
	return (f&0x7c00 == 0x7c00) && (f&0x03ff != 0)
}

// IsInf reports whether f is an infinity (inf).
// A sign > 0 reports whether f is positive inf.
// A sign < 0 reports whether f is negative inf.
// A sign == 0 reports whether f is either inf.
func (f Float16) IsInf(sign int) bool {
	return ((f == 0x7c00) && sign >= 0) ||
		(f == 0xfc00 && sign <= 0)
}

// IsFinite returns true if f is neither infinite nor NaN.
func (f Float16) IsFinite() bool {
	return (uint16(f) & uint16(0x7c00)) != uint16(0x7c00)
}

// IsNormal returns true if f is neither zero, infinite, subnormal, or NaN.
func (f Float16) IsNormal() bool {
	exp := uint16(f) & uint16(0x7c00)
	return (exp != uint16(0x7c00)) && (exp != 0)
}

// Signbit reports whether f is negative or negative zero.
func (f Float16) Signbit() bool {
	return (uint16(f) & uint16(0x8000)) != 0
}

// String satisfies the fmt.Stringer interface.
func (f Float16) String() string {
	return strconv.FormatFloat(float64(f.Float32()), 'f', -1, 32)
}

// f16bitsToF32bits returns uint32 (float32 bits) converted from specified uint16.
func f16bitsToF32bits(in uint16) uint32 {
	// All 65536 conversions with this were confirmed to be correct
	// by Montgomery Edwards⁴⁴⁸ (github.com/x448).

	sign := uint32(in&0x8000) << 16 // sign for 32-bit
	exp := uint32(in&0x7c00) >> 10  // exponenent for 16-bit
	coef := uint32(in&0x03ff) << 13 // significand for 32-bit

	if exp == 0x1f {
		if coef == 0 {
			// infinity
			return sign | 0x7f800000 | coef
		}
		// NaN
		return sign | 0x7fc00000 | coef
	}

	if exp == 0 {
		if coef == 0 {
			// zero
			return sign
		}

		// normalize subnormal numbers
		exp++
		for coef&0x7f800000 == 0 {
			coef <<= 1
			exp--
		}
		coef &= 0x007fffff
	}

	return sign | ((exp + (0x7f - 0xf)) << 23) | coef
}

// f32bitsToF16bits returns uint16 (Float16 bits) converted from the specified float32.
// Conversion rounds to nearest integer with ties to even.
func f32bitsToF16bits(u32 uint32) uint16 {
	// Translated from Rust to Go by Montgomery Edwards⁴⁴⁸ (github.com/x448).
	// All 4294967296 conversions with this were confirmed to be correct by x448.
	// Original Rust implementation is by Kathryn Long (github.com/starkat99) with MIT license.

	sign := u32 & 0x80000000
	exp := u32 & 0x7f800000
	coef := u32 & 0x007fffff

	if exp == 0x7f800000 {
		// NaN or Infinity
		nanBit := uint32(0)
		if coef != 0 {
			nanBit = uint32(0x0200)
		}
		return uint16((sign >> 16) | uint32(0x7c00) | nanBit | (coef >> 13))
	}

	halfSign := sign >> 16

	unbiasedExp := int32(exp>>23) - 127
	halfExp := unbiasedExp + 15

	if halfExp >= 0x1f {
		return uint16(halfSign | uint32(0x7c00))
	}

	if halfExp <= 0 {
		if 14-halfExp > 24 {
			return uint16(halfSign)
		}
		coef := coef | uint32(0x00800000)
		halfCoef := coef >> uint32(14-halfExp)
		roundBit := uint32(1) << uint32(13-halfExp)
		if (coef&roundBit) != 0 && (coef&(3*roundBit-1)) != 0 {
			halfCoef++
		}
		return uint16(halfSign | halfCoef)
	}

	uHalfExp := uint32(halfExp) << 10
	halfCoef := coef >> 13
	roundBit := uint32(0x00001000)
	if (coef&roundBit) != 0 && (coef&(3*roundBit-1)) != 0 {
		return uint16((halfSign | uHalfExp | halfCoef) + 1)
	}
	return uint16(halfSign | uHalfExp | halfCoef)
}
