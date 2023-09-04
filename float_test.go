package randomfloat

import (
	"math/rand"
	"testing"
)

func TestFloat32(t *testing.T) {
	r := New(rand.NewSource(42))
	for i := 0; i < 10000; i++ {
		f := r.Float32()
		t.Logf("%x", f)
	}
}

func BenchmarkFloat32(b *testing.B) {
	r := New(rand.NewSource(42))
	for i := 0; i < b.N; i++ {
		r.Float32()
	}
}
