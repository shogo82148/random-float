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

func (r *Rand) float32src() float32 {
	var exp = bias32 - 1
	var frac uint32
	for {
		i := r.src.Int63()
		l := bits.Len64(uint64(i))
		exp -= 63 - l
		if exp <= 0 {
			frac = uint32(r.src.Int63())
			exp = 0
			break
		}
		if l > shift32 {
			frac = uint32(i >> (l - shift32 - 1))
			break
		} else if l > 0 {
			frac = uint32(i << (shift32 - l + 1))
			i = r.src.Int63() >> (63 - shift32 + l - 1)
			frac |= uint32(i)
			break
		}
	}
	return math.Float32frombits(uint32(exp)<<shift32 | frac&fracMask32)
}

func (r *Rand) float32s64() float32 {
	var exp = bias32 - 1
	var frac uint32
	for {
		i := r.s64.Uint64()
		l := bits.Len64(i)
		exp -= 64 - l
		if exp <= 0 {
			frac = uint32(r.s64.Uint64())
			exp = 0
			break
		}
		if l > shift32 {
			frac = uint32(i >> (l - shift32 - 1))
			break
		} else if l > 0 {
			frac = uint32(i << (shift32 - l + 1))
			i = r.s64.Uint64() >> (64 - shift32 + l - 1)
			frac |= uint32(i)
			break
		}
	}
	return math.Float32frombits(uint32(exp)<<shift32 | frac&fracMask32)
}

func (r *Rand) float64src() float64 {
	var exp = bias64 - 1
	var frac uint64
	for {
		i := r.src.Int63()
		l := bits.Len64(uint64(i))
		exp -= 63 - l
		if exp <= 0 {
			frac = uint64(r.src.Int63())
			exp = 0
			break
		}
		if l > shift64 {
			frac = uint64(i >> (l - shift64 - 1))
			break
		} else if l > 0 {
			frac = uint64(i << (shift64 - l + 1))
			i = r.src.Int63() >> (63 - shift64 + l - 1)
			frac |= uint64(i)
			break
		}
	}
	return math.Float64frombits(uint64(exp)<<shift64 | frac&fracMask64)
}

func (r *Rand) float64s64() float64 {
	var exp = bias64 - 1
	var frac uint64
	for {
		i := r.s64.Uint64()
		l := bits.Len64(i)
		exp -= 64 - l
		if exp <= 0 {
			frac = uint64(r.s64.Uint64())
			exp = 0
			break
		}
		if l > shift64 {
			frac = uint64(i >> (l - shift64 - 1))
			break
		} else if l > 0 {
			frac = uint64(i << (shift64 - l + 1))
			i = r.s64.Uint64() >> (64 - shift64 + l - 1)
			frac |= uint64(i)
			break
		}
	}
	return math.Float64frombits(uint64(exp)<<shift64 | frac&fracMask64)
}

func (r *Rand) Float32() float32 {
	if r.s64 != nil {
		return r.float32s64()
	} else {
		return r.float32src()
	}
}

func (r *Rand) Float64() float64 {
	if r.s64 != nil {
		return r.float64s64()
	} else {
		return r.float64src()
	}
}
