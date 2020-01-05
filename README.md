# Float16 (Binary16) in Go/Golang
[![Build Status](https://travis-ci.org/cbor-go/float16.svg?branch=master)](https://travis-ci.org/cbor-go/float16)
[![codecov](https://codecov.io/gh/cbor-go/float16/branch/master/graph/badge.svg?v=4)](https://codecov.io/gh/cbor-go/float16)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbor-go/float16)](https://goreportcard.com/report/github.com/cbor-go/float16)
[![Release](https://img.shields.io/github/release/cbor-go/float16.svg?style=flat-square)](https://github.com/cbor-go/float16/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/cbor-go/float16/master/LICENSE)

`float16` package provides [IEEE 754 half-precision floating-point format (binary16)](https://en.wikipedia.org/wiki/Half-precision_floating-point_format) with IEEE 754 default rounding for conversions. IEEE 754-2008 refers to this 16-bit floating-point format as binary16.

IEEE 754 default rounding ("Round-to-Nearest RoundTiesToEven") is considered the most accurate and statistically unbiased estimate of the true result.

All possible 4+ billion floating-point conversions with this library are verified to be correct.

This library uses the lowercase word "float16" to refer to IEEE 754 binary16. And uses capitalized "Float16" to export a Go data type representing float16.

## Features
Current features include:

* float16 to float32 conversions use lossless conversion.
* float32 to float16 conversions use IEEE 754-2008 "Round-to-Nearest RoundTiesToEven".
* conversions use __zero allocs__ and are about __2.65 ns/op__ (in pure Go) on a desktop amd64.
* unit tests provide 100% code coverage and check all possible 4+ billion conversions.
* other functions include: IsInf(), IsNaN(), IsNormal(), PrecisionFromfloat32(), String(), etc.
* all functions in this library use zero allocs except String().

## Status
This library is used by [fxamacker/cbor](https://github.com/fxamacker/cbor) and is ready for production use on supported platforms. The version number < 1.0 indicates more functions and options are planned but not yet published.

Current status:

* core API is done and breaking API changes are unlikely.
* 100% of unit tests pass:
  * short mode (`go test -short`) tests around 65765 conversions in 0.005s.  
  * normal mode (`go test`) tests all possible 4+ billion conversions in about 75s.  
* 100% code coverage with both short mode and normal mode.  
* tested on amd64 but it should work on all little-endian platforms supported by Go.
 
Roadmap:

* add functions for fast batch conversions.
* speed up unit test when verifying all possible 4+ billion conversions.
* test on additional platforms.
 
## Float16 to Float32 Conversion
Conversions from float16 to float32 are lossless conversions.  All 65536 possible float16 to float32 conversions (in pure Go) are confirmed to be correct.  

Unit tests take a fraction of a second to check all 65536 expected values for float16 to float32 conversions.

## Float32 to Float16 Conversion
Conversions from float32 to float16 use IEEE 754 default rounding ("Round-to-Nearest RoundTiesToEven").  All 4294967296 possible float32 to float16 conversions (in pure Go) are confirmed to be correct.  

Unit tests in normal mode take about 60-90 seconds to check all 4+ billion expected values for float32 to float16 conversions as well as PrecisionFromfloat32() for each.

Unit tests in short mode use a small subset (65765) of expected values and finish in under 0.01 second while still reaching 100% code coverage.

## Usage
Install with `go get github.com/cbor-go/float16`.
```
// Convert float32 to float16
pi := float32(math.Pi)
pi16 := float16.Fromfloat32(pi)

// Convert float16 to float32
pi32 := pi16.Float32()

// Only convert if there's no data loss (useful for CBOR encoders)
// PrecisionFromfloat32() is faster than the overhead of calling a function
if float16.Precision(pi) == float16.PrecisionExact {
	pi16 := float16.Fromfloat32(pi)
}
```

## Float16 Type and API
Float16 (capitalized) is a Go type with uint16 as the underlying state.  There are 4 exported functions and 7 exported  methods.
```
package float16 // import "github.com/cbor-go/float16"

// Exported types
type Float16 uint16

// Exported functions
Fromfloat32(f32 float32) Float16    // Float16 number converted from f32 using IEEE 754 default rounding
Frombits(b16 uint16) Float16        // Float16 number corresponding to b16 (IEEE 754 binary16 rep.)
NaN() Float16                       // Float16 of IEEE 754 binary16 not-a-number
Inf(sign int) Float16               // Float16 of IEEE 754 binary16 infinity according to sign

PrecisionFromfloat32(f32 float32) Precision  // quickly indicates exact, inexact, overflow, underflow
                                             // (inline and < 1 ns/op)
// Exported methods
(f Float16) Float32() float32       // float32 number converted from f16 using lossless conversion
(f Float16) Bits() uint16           // the IEEE 754 binary16 representation of f
(f Float16) IsNaN() bool            // true if f is not-a-number (NaN)
(f Float16) IsQuietNaN() bool       // true if f is a quiet not-a-number (NaN)
(f Float16) IsInf(sign int) bool    // true if f is infinite based on sign (-1=NegInf, 0=any, 1=PosInf)
(f Float16) IsFinite() bool         // true if f is not infinite or NaN
(f Float16) IsNormal() bool         // true if f is not zero, infinite, subnormal, or NaN.
(f Float16) Signbit() bool          // true if f is negative or negative zero
(f Float16) String() string         // string representation of f to satisfy fmt.Stringer interface
```
See [API](https://godoc.org/github.com/cbor-go/float16) at godoc.org for more info.

## Benchmarks
Conversions (in pure Go) are around 2.65 ns/op for float16 to Float32 as well as Float32 to float16 on amd64. And speeds can vary depending on input value.

Frombits is included as a canary to catch overoptimized benchmarks. It should be faster than all other functions except PrecisionFromfloat32.
```
All functions have zero allocations except float16.String().

FromFloat32pi-2  2.59ns ± 0%    // speed using Fromfloat32() to convert a float32 of math.Pi to Float16
ToFloat32pi-2    2.69ns ± 0%    // speed using Float32() to convert a float16 of math.Pi to float32
Frombits-2       0.29ns ± 5%    // speed using Frombits() to cast a uint16 to Float16

PrecisionFromFloat32-2  0.29ns ± 1%  // speed using PrecisionFromfloat32() to check for overflows, etc.
```

## System Requirements
* Tested on Go 1.11, 1.12, and 1.13 but it should also work with older versions.
* Tested on amd64 but it should also work on all little-endian platforms supported by Go.

## Special Thanks
Special thanks to Kathryn Long (starkat99) for creating [half-rs](https://github.com/starkat99/half-rs), a very nice rust implementation of float16.

## License
Copyright (c) 2019 Montgomery Edwards⁴⁴⁸ and Faye Amacker

Licensed under [MIT License](LICENSE)
