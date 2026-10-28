package spacecurves

import (
	"math/rand/v2"
	"testing"
)

// mortonEncodeRef interleaves the low bits of each coordinate one bit at a
// time, as a reference for the bit-trick encoders.
func mortonEncodeRef(bits uint, p ...uint64) uint64 {
	var d uint64
	for i := range bits {
		for k, v := range p {
			d |= (v >> i & 1) << (i*uint(len(p)) + uint(k))
		}
	}
	return d
}

func TestMorton2D(t *testing.T) {
	const bits = 5
	for y := range uint64(1 << bits) {
		for x := range uint64(1 << bits) {
			morton := Morton2DEncode(bits, x, y)
			if x2, y2 := Morton2DDecode(bits, morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D(t *testing.T) {
	const bits = 5
	for y := range uint64(1 << bits) {
		for x := range uint64(1 << bits) {
			for z := range uint64(1 << bits) {
				morton := Morton3DEncode(bits, x, y, z)
				if x2, y2, z2 := Morton3DDecode(bits, morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMortonMatchesReference(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 2))
	for bits := range uint(33) {
		vals := edgeValues(bits, rnd)
		for i, x := range vals {
			y := vals[(i*7+3)%len(vals)]
			d := Morton2DEncode(bits, x, y)
			if want := mortonEncodeRef(bits, x, y); d != want {
				t.Fatalf("bits=%d: Morton2DEncode(%#x, %#x) = %#x, expected %#x", bits, x, y, d, want)
			}
			if x2, y2 := Morton2DDecode(bits, d); x2 != x || y2 != y {
				t.Fatalf("bits=%d: Morton2DDecode(%#x) = (%#x, %#x), expected (%#x, %#x)", bits, d, x2, y2, x, y)
			}
		}
	}
	for bits := range uint(22) {
		vals := edgeValues(bits, rnd)
		for i, x := range vals {
			y, z := vals[(i*7+3)%len(vals)], vals[(i*13+5)%len(vals)]
			d := Morton3DEncode(bits, x, y, z)
			if want := mortonEncodeRef(bits, x, y, z); d != want {
				t.Fatalf("bits=%d: Morton3DEncode(%#x, %#x, %#x) = %#x, expected %#x", bits, x, y, z, d, want)
			}
			if x2, y2, z2 := Morton3DDecode(bits, d); x2 != x || y2 != y || z2 != z {
				t.Fatalf("bits=%d: Morton3DDecode(%#x) = (%#x, %#x, %#x), expected (%#x, %#x, %#x)", bits, d, x2, y2, z2, x, y, z)
			}
		}
	}
}

func TestMortonToHilbert2D(t *testing.T) {
	const bits = 5
	for y := range uint64(1 << bits) {
		for x := range uint64(1 << bits) {
			morton := Morton2DEncode(bits, x, y)
			hilbert := MortonToHilbert2D(bits, morton)
			if want := Hilbert2DEncode(bits, x, y); hilbert != want {
				t.Fatalf("MortonToHilbert2D for (%d, %d) = %d, Hilbert2DEncode = %d", x, y, hilbert, want)
			}
			if got := HilbertToMorton2D(bits, hilbert); got != morton {
				t.Fatalf("HilbertToMorton2D(MortonToHilbert2D(%d)) = %d", morton, got)
			}
		}
	}
	rnd := rand.New(rand.NewPCG(3, 4))
	for _, morton := range edgeValues(64, rnd) {
		x, y := Morton2DDecode(32, morton)
		if got, want := MortonToHilbert2D(32, morton), Hilbert2DEncode(32, x, y); got != want {
			t.Fatalf("MortonToHilbert2D(%#x) = %#x, Hilbert2DEncode = %#x", morton, got, want)
		}
		if got := HilbertToMorton2D(32, MortonToHilbert2D(32, morton)); got != morton {
			t.Fatalf("HilbertToMorton2D(MortonToHilbert2D(%#x)) = %#x", morton, got)
		}
	}
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
	for d := uint64(1); d < n*n; d++ {
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

func BenchmarkMorton2DEncode_16(b *testing.B) {
	for b.Loop() {
		Morton2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkMorton2DDecode_16(b *testing.B) {
	for b.Loop() {
		Morton2DDecode(16, 0x12345678)
	}
}

func BenchmarkMorton3DEncode_10(b *testing.B) {
	for b.Loop() {
		Morton3DEncode(10, 0x2cd, 0x134, 0x3a5)
	}
}

func BenchmarkMorton3DDecode_10(b *testing.B) {
	for b.Loop() {
		Morton3DDecode(10, 0x12345678)
	}
}

func BenchmarkMortonToHilbert2D_16(b *testing.B) {
	m := Morton2DEncode(16, 0xabcd, 0x1234)
	for b.Loop() {
		MortonToHilbert2D(16, m)
	}
}

func TestMortonND(t *testing.T) {
	rnd := rand.New(rand.NewPCG(5, 6))
	for dims := 1; dims <= 8; dims++ {
		for bits := range 64/uint(dims) + 1 {
			p := make([]uint64, dims)
			for _, d := range edgeValues(uint(dims)*bits, rnd) {
				MortonDecode(bits, d, p)
				if got := MortonEncode(bits, p); got != d {
					t.Fatalf("dims=%d bits=%d: MortonEncode(MortonDecode(%#x)) = %#x", dims, bits, d, got)
				}
				if want := mortonEncodeRef(bits, p...); want != d {
					t.Fatalf("dims=%d bits=%d: reference encoding of %v = %#x, expected %#x", dims, bits, p, want, d)
				}
			}
		}
	}
}
