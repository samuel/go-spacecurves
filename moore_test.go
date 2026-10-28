package spacecurves

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestMoore2D(t *testing.T) {
	for nbits := range uint(9) {
		total := uint(1) << (2 * nbits)
		for d := range total {
			x, y := Moore2DDecode(nbits, d)
			if x >= 1<<nbits || y >= 1<<nbits {
				t.Fatalf("bits=%d: index %d decodes to out of range point (%d, %d)", nbits, d, x, y)
			}
			if d2 := Moore2DEncode(nbits, x, y); d2 != d {
				t.Fatalf("bits=%d: Moore2DEncode(Moore2DDecode(%d)) = %d", nbits, d, d2)
			}
		}
		checkMooreLoop(t, nbits, 0, total)
	}
}

func TestMoore2DMaxBits(t *testing.T) {
	const nbits = uintBits / 2
	rnd := rand.New(rand.NewPCG(9, 10))
	for _, d := range edgeValues(2*nbits, rnd) {
		x, y := Moore2DDecode(nbits, d)
		if d2 := Moore2DEncode(nbits, x, y); d2 != d {
			t.Fatalf("Moore2DEncode(Moore2DDecode(%#x)) = %#x", d, d2)
		}
	}
	// Quadrant boundaries and the wraparound from the last index to 0.
	last := uint(math.MaxUint) >> (uintBits - 2*nbits)
	for q := range uint(4) {
		d := q << (2*nbits - 2)
		checkMooreLoop(t, nbits, d-1, d+1)
	}
	checkMooreLoop(t, nbits, last-1<<10, last)
}

func TestMoore2DMasksInput(t *testing.T) {
	const nbits = 4
	if got, want := Moore2DEncode(nbits, 0x13, 0x25), Moore2DEncode(nbits, 0x3, 0x5); got != want {
		t.Fatalf("Moore2DEncode with high bits = %d, expected %d", got, want)
	}
}

// checkMooreLoop verifies consecutive indices in [start, end] decode to
// adjacent points, wrapping modulo the curve length so the last index is
// checked against index 0.
func checkMooreLoop(t *testing.T, nbits, start, end uint) {
	t.Helper()
	if nbits == 0 {
		return
	}
	// Wraps to 0 when 2*nbits is the width of uint.
	mask := uint(1)<<(2*nbits-1)<<1 - 1
	start &= mask
	lx, ly := Moore2DDecode(nbits, start)
	for d := start; d != end&mask; {
		d = (d + 1) & mask
		x, y := Moore2DDecode(nbits, d)
		if dist(lx, x)+dist(ly, y) != 1 {
			t.Fatalf("bits=%d: indices %d and %d map to non-adjacent points (%d, %d) and (%d, %d)", nbits, (d-1)&mask, d, lx, ly, x, y)
		}
		lx, ly = x, y
	}
}

func TestDisplayMoore2D(t *testing.T) {
	logCurve2D(t, 5, func(d uint) (uint, uint) { return Moore2DDecode(5, d) })
}

func BenchmarkMoore2DEncode_16(b *testing.B) {
	for b.Loop() {
		Moore2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkMoore2DDecode_16(b *testing.B) {
	for b.Loop() {
		Moore2DDecode(16, 0x12345678)
	}
}
