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
