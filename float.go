package randomfloat

import (
	"math"
	"math/bits"
	"math/rand"
)

const (
	mask32     = 0xff       // mask for exponent
	shift32    = 32 - 8 - 1 // shift for exponent
	bias32     = 127        // bias for exponent
	signMask32 = 1 << 31    // mask for sign bit
	fracMask32 = 1<<shift32 - 1
)

const (
	mask64     = 0x7ff       // mask for exponent
	shift64    = 64 - 11 - 1 // shift for exponent
	bias64     = 1023        // bias for exponent
	signMask64 = 1 << 63     // mask for sign bit
	fracMask64 = 1<<shift64 - 1
)

type Rand struct {
	src rand.Source
	s64 rand.Source64 // non-nil if src is source64
}

// New returns a new Rand that uses random values from src
// to generate other random values.
func New(src rand.Source) *Rand {
	s64, _ := src.(rand.Source64)
	return &Rand{src: src, s64: s64}
}

type bitsN interface {
	uint64 | uint32
}

type intN interface {
	int64 | uint64
}

type floatN interface {
	float64 | float32
}

type floatNFromBits[B bitsN, F floatN] func(B) F

type srcN[I intN] func() I

func randFloat[I intN, B bitsN, F floatN](src srcN[I], bias, shift, num int, mask B, float floatNFromBits[B, F]) F {
	var exp = bias - 1
	var frac B
	for {
		i := src()
		l := bits.Len64(uint64(i))
		exp -= num - l
		if exp <= 0 {
			frac = B(src())
			exp = 0
			break
		}
		if l > shift {
			frac = B(i >> (l - shift - 1))
			break
		} else if l > 0 {
			frac = B(i << (shift - l + 1))
			i = src() >> (num - shift + l - 1)
			frac |= B(i)
			break
		}
	}
	return float(B(exp)<<shift | frac&mask)
}

func (r *Rand) Float32() float32 {
	if r.s64 != nil {
		return randFloat(r.s64.Uint64, bias32, shift32, 64, fracMask32, math.Float32frombits)
	} else {
		return randFloat(r.src.Int63, bias32, shift32, 63, fracMask32, math.Float32frombits)
	}
}

func (r *Rand) Float64() float64 {
	if r.s64 != nil {
		return randFloat(r.s64.Uint64, bias64, shift64, 64, fracMask64, math.Float64frombits)
	} else {
		return randFloat(r.src.Int63, bias64, shift64, 63, fracMask64, math.Float64frombits)
	}
}
