package spacecurves

import "testing"

func TestHilbert2D(t *testing.T) {
	const n = 8
	for y := uint(0); y < 1<<n; y++ {
		for x := uint(0); x < 1<<n; x++ {
			if x2, y2 := Hilbert2DDecode(n, Hilbert2DEncode(n, x, y)); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestDisplayHilbert2D(t *testing.T) {
	const bits = 5
	const n = 1 << bits
	var out [(2 * n) * (2 * n)]byte
	for i := 0; i < len(out); i++ {
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
	for y := 0; y < n*2; y++ {
		t.Logf("%s\n", string(out[y*n*2:y*n*2+n*2]))
	}
}

func BenchmarkHilbert2DEncode_5(b *testing.B) {
	const n = 5
	for i := 0; i < b.N; i++ {
		for y := uint(0); y < 1<<n; y++ {
			for x := uint(0); x < 1<<n; x++ {
				Hilbert2DEncode(n, x, y)
			}
		}
	}
}
