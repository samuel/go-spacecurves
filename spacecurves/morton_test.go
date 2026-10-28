package spacecurves

import (
	"math"
	"math/rand/v2"
	"testing"
)

// edgeValues returns 0, the max value, alternating bit patterns, and random
// values, all within the given bit width.
func edgeValues(bits uint, rnd *rand.Rand) []uint {
	mask := uint(1)<<bits - 1
	vals := []uint{0, 1, mask, mask - 1, uint(0x5555555555555555&math.MaxUint) & mask, uint(0xaaaaaaaaaaaaaaaa&math.MaxUint) & mask}
	for range 64 {
		vals = append(vals, uint(rnd.Uint64())&mask)
	}
	return vals
}

func TestMorton2D(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			morton := Morton2DEncode(bits, x, y)
			if x2, y2 := Morton2DDecode(bits, morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			for z := range uint(1 << bits) {
				morton := Morton3DEncode(bits, x, y, z)
				if x2, y2, z2 := Morton3DDecode(bits, morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMorton2D5bit(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			morton := Morton2DEncode5bit(x, y)
			if want := Morton2DEncode(bits, x, y); morton != want {
				t.Fatalf("Morton2DEncode5bit(%d, %d) = %d, expected %d", x, y, morton, want)
			}
			if x2, y2 := Morton2DDecode5bit(morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D5bit(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			for z := range uint(1 << bits) {
				morton := Morton3DEncode5bit(x, y, z)
				if want := Morton3DEncode(bits, x, y, z); morton != want {
					t.Fatalf("Morton3DEncode5bit(%d, %d, %d) = %d, expected %d", x, y, z, morton, want)
				}
				if x2, y2, z2 := Morton3DDecode5bit(morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMorton2D16bit(t *testing.T) {
	const bits = 16
	rnd := rand.New(rand.NewPCG(1, 2))
	vals := edgeValues(bits, rnd)
	for _, y := range vals {
		for _, x := range vals {
			morton := Morton2DEncode16bit(x, y)
			if want := Morton2DEncode(bits, x, y); morton != want {
				t.Fatalf("Morton2DEncode16bit(%d, %d) = %d, expected %d", x, y, morton, want)
			}
			if x2, y2 := Morton2DDecode16bit(morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D10bit(t *testing.T) {
	const bits = 10
	rnd := rand.New(rand.NewPCG(3, 4))
	vals := edgeValues(bits, rnd)
	for _, y := range vals {
		for _, x := range vals {
			for _, z := range vals {
				morton := Morton3DEncode10bit(x, y, z)
				if want := Morton3DEncode(bits, x, y, z); morton != want {
					t.Fatalf("Morton3DEncode10bit(%d, %d, %d) = %d, expected %d", x, y, z, morton, want)
				}
				if x2, y2, z2 := Morton3DDecode10bit(morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMortonToHilbert2D(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			morton := Morton2DEncode(bits, x, y)
			hilbert := MortonToHilbert2D(bits, morton)
			if got := HilbertToMorton2D(bits, hilbert); got != morton {
				t.Fatalf("HilbertToMorton2D(MortonToHilbert2D(%d)) = %d", morton, got)
			}
		}
	}
	checkAdjacent2D(t, bits, 0, 1<<(2*bits), func(d uint) (uint, uint) {
		return Morton2DDecode(bits, HilbertToMorton2D(bits, d))
	})
}

func TestMortonToHilbert3D(t *testing.T) {
	const bits = 5
	for y := range uint(1 << bits) {
		for x := range uint(1 << bits) {
			for z := range uint(1 << bits) {
				morton := Morton3DEncode(bits, x, y, z)
				hilbert := MortonToHilbert3D(bits, morton)
				if got := HilbertToMorton3D(bits, hilbert); got != morton {
					t.Fatalf("HilbertToMorton3D(MortonToHilbert3D(%d)) = %d", morton, got)
				}
			}
		}
	}
	checkAdjacent3D(t, 0, 1<<(3*bits), func(d uint) (uint, uint, uint) {
		return Morton3DDecode(bits, HilbertToMorton3D(bits, d))
	})
}

func TestMortonToHilbert3DHighBits(t *testing.T) {
	// Maximum bits per axis: 21 for 64-bit uint, 10 for 32-bit.
	const bits = 10 + 11*(math.MaxUint>>32&1)
	rnd := rand.New(rand.NewPCG(5, 6))
	for _, morton := range edgeValues(3*bits, rnd) {
		if got := HilbertToMorton3D(bits, MortonToHilbert3D(bits, morton)); got != morton {
			t.Fatalf("HilbertToMorton3D(MortonToHilbert3D(%#x)) = %#x", morton, got)
		}
	}
	const n = 1 << 12
	decode := func(d uint) (uint, uint, uint) {
		return Morton3DDecode(bits, HilbertToMorton3D(bits, d))
	}
	checkAdjacent3D(t, 1<<(3*bits)-n, 1<<(3*bits), decode)
	for range 16 {
		start := uint(rnd.Uint64N(1<<(3*bits) - n))
		checkAdjacent3D(t, start, start+n, decode)
	}
	// Indices where high bits carry over.
	for j := uint(1); j < 3*bits; j++ {
		for range 16 {
			d := uint(rnd.Uint64N(1<<(3*bits-j)-1)+1) << j
			checkAdjacent3D(t, d-1, d+1, decode)
		}
	}
}

// checkAdjacent3D verifies consecutive curve indices in [start, end) decode
// to points one unit step apart.
func checkAdjacent3D(t *testing.T, start, end uint, decode func(uint) (uint, uint, uint)) {
	t.Helper()
	lx, ly, lz := decode(start)
	for d := start + 1; d < end; d++ {
		x, y, z := decode(d)
		if dist(lx, x)+dist(ly, y)+dist(lz, z) != 1 {
			t.Fatalf("Indices %d and %d map to non-adjacent points (%d, %d, %d) and (%d, %d, %d)", d-1, d, lx, ly, lz, x, y, z)
		}
		lx, ly, lz = x, y, z
	}
}

func dist(a, b uint) uint {
	if a > b {
		return a - b
	}
	return b - a
}

func TestDisplayMorton2D(t *testing.T) {
	const bits = 5
	const n = 1 << bits
	var out [(2 * n) * (2 * n)]byte
	for i := range len(out) {
		out[i] = ' '
	}
	out[0] = '+'
	lastX, lastY := Morton2DDecode(bits, 0)
	for d := uint(1); d < n*n; d++ {
		x, y := Morton2DDecode(bits, d)
		o := y*n*4 + x*2
		out[o] = '+'
		switch {
		case lastX < x && lastY < y:
			out[o-1-n*2] = '\\'
		case lastX < x && lastY > y:
			out[o-1+n*2] = '/'
		case lastX > x && lastY < y:
			out[o+1-n*2] = '/'
		case lastX > x && lastY > y:
			out[o+1+n*2] = '\\'
		case lastX < x:
			out[o-1] = '-'
		case lastX > x:
			out[o+1] = '-'
		case lastY < y:
			out[o-n*2] = '|'
		case lastY > y:
			out[o+n*2] = '|'
		}
		lastX, lastY = x, y
	}
	for y := range n * 2 {
		t.Logf("%s\n", string(out[y*n*2:y*n*2+n*2]))
	}
}

func BenchmarkMorton2DEncode_5(b *testing.B) {
	const n = 5
	for b.Loop() {
		for y := range uint(1 << n) {
			for x := range uint(1 << n) {
				Morton2DEncode(n, x, y)
			}
		}
	}
}

func BenchmarkMorton2DEncode5bit(b *testing.B) {
	const n = 5
	for b.Loop() {
		for y := range uint(1 << n) {
			for x := range uint(1 << n) {
				Morton2DEncode5bit(x, y)
			}
		}
	}
}

func BenchmarkMorton2DEncode_16(b *testing.B) {
	for b.Loop() {
		Morton2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkMorton2DEncode16bit(b *testing.B) {
	for b.Loop() {
		Morton2DEncode16bit(0xabcd, 0x1234)
	}
}

func BenchmarkMorton3DEncode_10(b *testing.B) {
	for b.Loop() {
		Morton3DEncode(10, 0x2cd, 0x134, 0x3a5)
	}
}

func BenchmarkMorton3DEncode10bit(b *testing.B) {
	for b.Loop() {
		Morton3DEncode10bit(0x2cd, 0x134, 0x3a5)
	}
}

func BenchmarkMortonToHilbert2D_16(b *testing.B) {
	m := Morton2DEncode16bit(0xabcd, 0x1234)
	for b.Loop() {
		MortonToHilbert2D(16, m)
	}
}

func BenchmarkMortonToHilbert3D_10(b *testing.B) {
	m := Morton3DEncode10bit(0x2cd, 0x134, 0x3a5)
	for b.Loop() {
		MortonToHilbert3D(10, m)
	}
}
