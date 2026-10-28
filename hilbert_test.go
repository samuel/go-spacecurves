package spacecurves

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestHilbert2D(t *testing.T) {
	const n = 8
	for y := range uint(1 << n) {
		for x := range uint(1 << n) {
			if x2, y2 := Hilbert2DDecode(n, Hilbert2DEncode(n, x, y)); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
	checkAdjacent2D(t, n, 0, 1<<(2*n), func(d uint) (uint, uint) {
		return Hilbert2DDecode(n, d)
	})
}

// MortonToHilbert2D produces the transpose of the Hilbert2DEncode curve.
func TestHilbert2DTransposeOfMortonToHilbert2D(t *testing.T) {
	const n = 6
	for y := range uint(1 << n) {
		for x := range uint(1 << n) {
			got := MortonToHilbert2D(n, Morton2DEncode(n, y, x))
			if want := Hilbert2DEncode(n, x, y); got != want {
				t.Fatalf("MortonToHilbert2D for (%d, %d) = %d, Hilbert2DEncode = %d", x, y, got, want)
			}
		}
	}
}

// checkAdjacent2D verifies consecutive curve indices in [start, end) decode
// to points one unit step apart and within the bit width.
func checkAdjacent2D(t *testing.T, bits, start, end uint, decode func(uint) (uint, uint)) {
	t.Helper()
	lx, ly := decode(start)
	for d := start + 1; d < end; d++ {
		x, y := decode(d)
		if x >= 1<<bits || y >= 1<<bits {
			t.Fatalf("Index %d maps to out of range point (%d, %d)", d, x, y)
		}
		if dist(lx, x)+dist(ly, y) != 1 {
			t.Fatalf("Indices %d and %d map to non-adjacent points (%d, %d) and (%d, %d)", d-1, d, lx, ly, x, y)
		}
		lx, ly = x, y
	}
}

func TestDisplayHilbert2D(t *testing.T) {
	const bits = 5
	const n = 1 << bits
	var out [(2 * n) * (2 * n)]byte
	for i := range len(out) {
		out[i] = ' '
	}
	out[0] = '+'
	lastX, lastY := Hilbert2DDecode(bits, 0)
	for d := uint(1); d < n*n; d++ {
		x, y := Hilbert2DDecode(bits, d)
		o := y*n*4 + x*2
		out[o] = '+'
		if lastX < x {
			out[o-1] = '-'
		} else if lastY < y {
			out[o-n*2] = '|'
		} else if lastX > x {
			out[o+1] = '-'
		} else if lastY > y {
			out[o+n*2] = '|'
		}
		lastX, lastY = x, y
	}
	for y := range n * 2 {
		t.Logf("%s\n", string(out[y*n*2:y*n*2+n*2]))
	}
}

func BenchmarkHilbert2DEncode_5(b *testing.B) {
	const n = 5
	for b.Loop() {
		for y := range uint(1 << n) {
			for x := range uint(1 << n) {
				Hilbert2DEncode(n, x, y)
			}
		}
	}
}

func BenchmarkHilbert2DDecode_16(b *testing.B) {
	for b.Loop() {
		Hilbert2DDecode(16, 0x12345678)
	}
}

func TestHilbertND(t *testing.T) {
	for dims := 1; dims <= 6; dims++ {
		for nbits := uint(0); nbits*uint(dims) <= 12; nbits++ {
			total := uint(1) << (nbits * uint(dims))
			got := make([]uint, dims)
			for d := range total {
				HilbertDecode(nbits, d, got)
				for _, v := range got {
					if v >= 1<<nbits {
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
		nbits := uintBits / uint(dims)
		total := nbits * uint(dims)
		got := make([]uint, dims)
		for _, d := range edgeValues(total, rnd) {
			HilbertDecode(nbits, d, got)
			if d2 := HilbertEncode(nbits, got); d2 != d {
				t.Fatalf("dims=%d bits=%d: HilbertEncode(HilbertDecode(%#x)) = %#x", dims, nbits, d, d2)
			}
		}
		last := uint(math.MaxUint) >> (uintBits - total)
		checkAdjacentND(t, nbits, dims, last-1<<10, last)
		// Indices where high bits carry over.
		for j := uint(1); j < total; j++ {
			for range 8 {
				d := uint(rnd.Uint64N(uint64(last>>j))+1) << j
				checkAdjacentND(t, nbits, dims, d-1, d+1)
			}
		}
	}
}

func TestHilbertNDMatchesHilbert2D(t *testing.T) {
	const n = 6
	for y := range uint(1 << n) {
		for x := range uint(1 << n) {
			if got, want := HilbertEncode(n, []uint{x, y}), Hilbert2DEncode(n, x, y); got != want {
				t.Fatalf("HilbertEncode(%d, %d) = %d, Hilbert2DEncode = %d", x, y, got, want)
			}
		}
	}
}

func TestHilbert3D(t *testing.T) {
	const n = 5
	for d := range uint(1 << (3 * n)) {
		x, y, z := Hilbert3DDecode(n, d)
		got := make([]uint, 3)
		HilbertDecode(n, d, got)
		if got[0] != x || got[1] != y || got[2] != z {
			t.Fatalf("Hilbert3DDecode(%d) = (%d, %d, %d), HilbertDecode = %v", d, x, y, z, got)
		}
		if d2 := Hilbert3DEncode(n, x, y, z); d2 != d {
			t.Fatalf("Hilbert3DEncode(Hilbert3DDecode(%d)) = %d", d, d2)
		}
	}
	checkAdjacent3D(t, 0, 1<<(3*n), func(d uint) (uint, uint, uint) {
		return Hilbert3DDecode(n, d)
	})
}

func TestHilbertNDPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	HilbertEncode(uintBits/3+1, make([]uint, 3))
}

func TestHilbertNDDoesNotModifyPoint(t *testing.T) {
	p := []uint{3, 5, 6}
	HilbertEncode(3, p)
	if p[0] != 3 || p[1] != 5 || p[2] != 6 {
		t.Fatalf("HilbertEncode modified point: %v", p)
	}
}

// checkAdjacentND verifies consecutive curve indices in [start, end) decode
// to points one unit step apart.
func checkAdjacentND(t *testing.T, nbits uint, dims int, start, end uint) {
	t.Helper()
	last := make([]uint, dims)
	p := make([]uint, dims)
	HilbertDecode(nbits, start, last)
	for d := start + 1; d < end; d++ {
		HilbertDecode(nbits, d, p)
		var sum uint
		for i := range p {
			sum += dist(last[i], p[i])
		}
		if sum != 1 {
			t.Fatalf("dims=%d bits=%d: indices %d and %d map to non-adjacent points %v and %v", dims, nbits, d-1, d, last, p)
		}
		last, p = p, last
	}
}

func BenchmarkHilbert2DEncode_16(b *testing.B) {
	for b.Loop() {
		Hilbert2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkHilbertEncode2D_16(b *testing.B) {
	p := []uint{0xabcd, 0x1234}
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

func BenchmarkMortonToHilbert3DEncode_10(b *testing.B) {
	for b.Loop() {
		MortonToHilbert3D(10, Morton3DEncode10bit(0x2cd, 0x134, 0x3a5))
	}
}

func BenchmarkHilbertEncode4D_16(b *testing.B) {
	p := []uint{0xabcd, 0x1234, 0x5678, 0x9abc}
	for b.Loop() {
		HilbertEncode(16, p)
	}
}
