package spacecurves

import "testing"

func TestMorton2D(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			morton := Morton2DEncode(n, x, y)
			if x2, y2 := Morton2DDecode(n, morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			for z := uint(0); z < n; z++ {
				morton := Morton3DEncode(5, x, y, z)
				if x2, y2, z2 := Morton3DDecode(n, morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMorton2D5bit(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			morton := Morton2DEncode5bit(x, y)
			if x2, y2 := Morton2DDecode5bit(morton); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMorton3D5bit(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			for z := uint(0); z < n; z++ {
				morton := Morton3DEncode5bit(x, y, z)
				if x2, y2, z2 := Morton3DDecode5bit(morton); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestMortonToHilbert2D(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			morton := Morton2DEncode5bit(x, y)
			hilbert := MortonToHilbert2D(morton, 10)
			if x2, y2 := Morton2DDecode5bit(HilbertToMorton2D(hilbert, 10)); x != x2 || y != y2 {
				t.Fatalf("Got (%d, %d) expected (%d, %d)", x2, y2, x, y)
			}
		}
	}
}

func TestMortonToHilbert3D(t *testing.T) {
	const n = 1 << 5
	for y := uint(0); y < n; y++ {
		for x := uint(0); x < n; x++ {
			for z := uint(0); z < n; z++ {
				morton := Morton3DEncode5bit(x, y, z)
				hilbert := MortonToHilbert2D(morton, 15)
				if x2, y2, z2 := Morton3DDecode5bit(HilbertToMorton2D(hilbert, 15)); x != x2 || y != y2 || z != z2 {
					t.Fatalf("Got (%d, %d, %d) expected (%d, %d, %d)", x2, y2, z2, x, y, z)
				}
			}
		}
	}
}

func TestDisplayMorton2D(t *testing.T) {
	const bits = 5
	const n = 1 << bits
	var out [(2 * n) * (2 * n)]byte
	for i := 0; i < len(out); i++ {
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
		case lastX > x && lastY > x:
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
	for y := 0; y < n*2; y++ {
		t.Logf("%s\n", string(out[y*n*2:y*n*2+n*2]))
	}
}

func BenchmarkMorton2DEncode_5(b *testing.B) {
	const n = 5
	for i := 0; i < b.N; i++ {
		for y := uint(0); y < 1<<n; y++ {
			for x := uint(0); x < 1<<n; x++ {
				_ = Morton2DEncode(n, x, y)
			}
		}
	}
}

func BenchmarkMorton2DEncode5bit(b *testing.B) {
	const n = 5
	for i := 0; i < b.N; i++ {
		for y := uint(0); y < 1<<n; y++ {
			for x := uint(0); x < 1<<n; x++ {
				_ = Morton2DEncode5bit(x, y)
			}
		}
	}
}
