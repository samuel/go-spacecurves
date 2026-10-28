package spacecurves

import "testing"

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
