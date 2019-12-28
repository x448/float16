# float16 library in Go
float16 provides [IEEE 754 half-precision floating-point numbers](https://en.wikipedia.org/wiki/Half-precision_floating-point_format) (binary16) with IEEE 754 default rounding for conversions.  

All possible 4+ billion conversions between float16 and float32 are verified to be correct.

Conversions between float16 and float32:

* use IEEE 754-2008 "Round-to-Nearest RoundTiesToEven" when converting float32 to float16
* use lossless conversion when converting float16 to float32
* use zero allocations and take less than 3 ns/op (in pure Go) on a desktop CPU

This library is used by [fxamacker/cbor](https://github.com/fxamacker/cbor) v1.4+ and is ready for production use on supported platforms.

## Status
This float16 library produces correct results for all conversions between float16 and float32.

Current status:

* 100% of unit tests pass:
  * short mode (`go test -short`) tests around 65763 conversions in 0.005s.  
  * normal mode (`go test`) tests all possible 4+ billion conversions in about 45s.  
* 100% code coverage with both short mode and normal mode.  
* conversions use zero allocations and take less than 3 ns/op (in pure Go) on a desktop.
* tested on amd64 but it should work on all little-endian platforms supported by Go.
 
Roadmap: 
* add functions for fast batch conversions
* speed up unit test when verifying all possible 4+ billion conversions
* test on additional platforms
 
## float16 to float32 conversion
Conversions from float16 to float32 are lossless conversions.  All 65536 possible float16 to float32 conversions (in pure Go) are confirmed to be correct.  

Unit tests take a fraction of a second to check all 65536 expected values for float16 to float32 conversions.

## float32 to float16 conversion
Conversions from float32 to float16 use IEEE 754 default rounding ("Round-to-Nearest RoundTiesToEven").  All 4294967296 possible float32 to float16 conversions (in pure Go) are confirmed to be correct.  

Unit tests in normal mode take about 35-55 seconds to check all 4+ billion expected values for float32 to float16 conversions.  

Unit tests in short mode use a small subset (65763) of expected values and finish in under 1 second while still reaching 100% code coverage.

## API and Float16 type
Float16 is a Go type with a simple API.  The only state for each Float16 is a uint16.
```
package float16 // import "github.com/x448/float16"

// Exported types
type Float16 uint16

// Exported functions
Fromfloat32(f32 float32) Float16     // returns Float16 converted from f32 using IEEE 754 default rounding
Frombits(u16 uint16) Float16         // returns Float16 by casting uint16 to Float16
NaN() Float16                        // returns IEEE 754 half-precision not-a-number
Inf(sign int) Float16                // returns IEEE 754 half-precision infinity according to sign

// Exported methods
(f Float16) Float32() float32        // returns float32 converted from f16 using lossless conversion
(f Float16) IsNaN() bool             // returns true if f is not-a-number (NaN)
(f Float16) IsInf(sign int) bool     // returns true if f is infinite according to sign (-1=NegInf, 0=Both, 1=PosInf)
(f Float16) IsFinite() bool          // returns true if f is not infinite or NaN
(f Float16) IsNormal() bool          // returns true if f is not zero, infinite, subnormal, or NaN.
(f Float16) Signbit() bool           // returns true if f is negative or negative zero
(f Float16) String() string          // returns the string representation of f to satisfy fmt.Stringer interface
```

## Usage
```
// Convert float32 to float16
pi := float32(math.Pi)
pi16 := float16.Fromfloat32(pi)

// Convert float16 to float32
pi32 := pi16.Float32()
```

## Benchmarks
Conversions (in pure Go) are around 3 ns/op for float16 to Float32 as well as Float32 to float16.

```
All functions have zero allocations except float16.String().

FromFloat32pi-2  2.59ns ± 0%    // speed using Fromfloat32() to convert float32 of math.Pi to Float16
ToFloat32pi-2    2.69ns ± 0%    // speed using Float32() to convert float16 of math.Pi to float32
Frombits-2       0.36ns ± 8%    // speed using Frombits() to cast uint16 to Float16
```

## System Requirements
* Tested on Go 1.12 and Go 1.13 but it should also work with older versions.
* Tested on amd64 but it should also work on all little-endian platforms supported by Go.

## Special Thanks
Special thanks to Kathryn Long (starkat99) for creating [half-rs](https://github.com/starkat99/half-rs), a very nice rust implementation of float16. 


