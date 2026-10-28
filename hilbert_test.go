package spacecurves

import (
	"math/rand/v2"
	"testing"
)

func TestHilbert2D(t *testing.T) {
	const n = 8
	for y := range uint64(1 << n) {
		for x := range uint64(1 << n) {
			if x2, y2 := Hilbert2DDecode(n, Hilbert2DEncode(n, x, y)); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
	checkAdjacent2D(t, n, 0, 1<<(2*n), func(d uint64) (uint64, uint64) {
		return Hilbert2DDecode(n, d)
	})
}

func TestHilbertND(t *testing.T) {
	for dims := 1; dims <= 6; dims++ {
		for nbits := uint(0); nbits*uint(dims) <= 12; nbits++ {
			total := uint64(1) << (nbits * uint(dims))
			got := make([]uint64, dims)
			for d := range total {
				HilbertDecode(nbits, d, got)
				for _, v := range got {
					if v > lowMask(nbits) {
						t.Fatalf("dims=%d bits=%d: index %d decodes to out of range point %v", dims, nbits, d, got)
					}
				}
				if d2 := HilbertEncode(nbits, got); d2 != d {
					t.Fatalf("dims=%d bits=%d: HilbertEncode(HilbertDecode(%d)) = %d", dims, nbits, d, d2)
				}
			}
			checkAdjacentND(t, nbits, dims, 0, total)
		}
	}
}

func TestHilbertNDMaxBits(t *testing.T) {
	rnd := rand.New(rand.NewPCG(7, 8))
	for dims := 1; dims <= 8; dims++ {
		nbits := 64 / uint(dims)
		total := nbits * uint(dims)
		got := make([]uint64, dims)
		for _, d := range edgeValues(total, rnd) {
			HilbertDecode(nbits, d, got)
			if d2 := HilbertEncode(nbits, got); d2 != d {
				t.Fatalf("dims=%d bits=%d: HilbertEncode(HilbertDecode(%#x)) = %#x", dims, nbits, d, d2)
			}
		}
		last := lowMask(total)
		checkAdjacentND(t, nbits, dims, last-1<<10, last)
		// Indices where high bits carry over.
		for j := uint(1); j < total; j++ {
			for range 8 {
				d := (rnd.Uint64N(last>>j) + 1) << j
				checkAdjacentND(t, nbits, dims, d-1, d+1)
			}
		}
	}
}

func TestHilbertNDMatchesHilbert2D(t *testing.T) {
	const n = 6
	for y := range uint64(1 << n) {
		for x := range uint64(1 << n) {
			if got, want := HilbertEncode(n, []uint64{x, y}), Hilbert2DEncode(n, x, y); got != want {
				t.Fatalf("HilbertEncode(%d, %d) = %d, Hilbert2DEncode = %d", x, y, got, want)
			}
		}
	}
}

func TestHilbert3D(t *testing.T) {
	const n = 5
	got := make([]uint64, 3)
	for d := range uint64(1 << (3 * n)) {
		x, y, z := Hilbert3DDecode(n, d)
		HilbertDecode(n, d, got)
		if got[0] != x || got[1] != y || got[2] != z {
			t.Fatalf("Hilbert3DDecode(%d) = (%d, %d, %d), HilbertDecode = %v", d, x, y, z, got)
		}
		if d2 := Hilbert3DEncode(n, x, y, z); d2 != d {
			t.Fatalf("Hilbert3DEncode(Hilbert3DDecode(%d)) = %d", d, d2)
		}
	}
	checkAdjacent3D(t, 0, 1<<(3*n), func(d uint64) (uint64, uint64, uint64) {
		return Hilbert3DDecode(n, d)
	})
}

func TestHilbertNDDoesNotModifyPoint(t *testing.T) {
	p := []uint64{3, 5, 6 | 1<<40}
	HilbertEncode(3, p)
	if p[0] != 3 || p[1] != 5 || p[2] != 6|1<<40 {
		t.Fatalf("HilbertEncode modified point: %v", p)
	}
}

// checkAdjacentND verifies consecutive curve indices in [start, end) decode
// to points one unit step apart.
func checkAdjacentND(t *testing.T, nbits uint, dims int, start, end uint64) {
	t.Helper()
	last := make([]uint64, dims)
	p := make([]uint64, dims)
	HilbertDecode(nbits, start, last)
	for d := start + 1; d < end; d++ {
		HilbertDecode(nbits, d, p)
		var sum uint64
		for i := range p {
			sum += dist(last[i], p[i])
		}
		if sum != 1 {
			t.Fatalf("dims=%d bits=%d: indices %d and %d map to non-adjacent points %v and %v", dims, nbits, d-1, d, last, p)
		}
		last, p = p, last
	}
}

func TestDisplayHilbert2D(t *testing.T) {
	logCurve2D(t, 5, func(d uint64) (uint64, uint64) { return Hilbert2DDecode(5, d) })
}

func BenchmarkHilbert2DEncode_16(b *testing.B) {
	for b.Loop() {
		Hilbert2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkHilbert2DDecode_16(b *testing.B) {
	for b.Loop() {
		Hilbert2DDecode(16, 0x12345678)
	}
}

func BenchmarkHilbertEncode2D_16(b *testing.B) {
	p := []uint64{0xabcd, 0x1234}
	for b.Loop() {
		HilbertEncode(16, p)
	}
}

func BenchmarkHilbert3DEncode_10(b *testing.B) {
	for b.Loop() {
		Hilbert3DEncode(10, 0x2cd, 0x134, 0x3a5)
	}
}

func BenchmarkHilbert3DDecode_10(b *testing.B) {
	for b.Loop() {
		Hilbert3DDecode(10, 0x12345678)
	}
}

func BenchmarkHilbertEncode4D_16(b *testing.B) {
	p := []uint64{0xabcd, 0x1234, 0x5678, 0x9abc}
	for b.Loop() {
		HilbertEncode(16, p)
	}
}
