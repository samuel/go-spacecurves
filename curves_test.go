package spacecurves

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// edgeValues returns 0, the max value, alternating bit patterns, and random
// values, all within the given bit width.
func edgeValues(bits uint, rnd *rand.Rand) []uint64 {
	mask := lowMask(bits)
	vals := []uint64{0, 1 & mask, mask, (mask - 1) & mask, 0x5555555555555555 & mask, 0xaaaaaaaaaaaaaaaa & mask}
	for range 64 {
		vals = append(vals, rnd.Uint64()&mask)
	}
	return vals
}

func dist(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

// checkAdjacent2D verifies consecutive curve indices in [start, end) decode
// to points one unit step apart and within the bit width.
func checkAdjacent2D(t *testing.T, bits uint, start, end uint64, decode func(uint64) (uint64, uint64)) {
	t.Helper()
	lx, ly := decode(start)
	for d := start + 1; d < end; d++ {
		x, y := decode(d)
		if x > lowMask(bits) || y > lowMask(bits) {
			t.Fatalf("Index %d maps to out of range point (%d, %d)", d, x, y)
		}
		if dist(lx, x)+dist(ly, y) != 1 {
			t.Fatalf("Indices %d and %d map to non-adjacent points (%d, %d) and (%d, %d)", d-1, d, lx, ly, x, y)
		}
		lx, ly = x, y
	}
}

// checkAdjacent3D verifies consecutive curve indices in [start, end) decode
// to points one unit step apart.
func checkAdjacent3D(t *testing.T, start, end uint64, decode func(uint64) (uint64, uint64, uint64)) {
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

// logCurve2D logs an ASCII rendering of a 2D curve made of unit steps.
func logCurve2D(t *testing.T, bits uint, decode func(uint64) (uint64, uint64)) {
	t.Helper()
	n := uint64(1) << bits
	out := make([]byte, (2*n)*(2*n))
	for i := range out {
		out[i] = ' '
	}
	lastX, lastY := decode(0)
	out[lastY*n*4+lastX*2] = '+'
	for d := uint64(1); d < n*n; d++ {
		x, y := decode(d)
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

// TestGolden pins curve outputs, since stored indices depend on them.
func TestGolden(t *testing.T) {
	order2D := func(decode func(uint, uint64) (uint64, uint64)) []uint64 {
		var out []uint64
		for d := range uint64(16) {
			x, y := decode(2, d)
			out = append(out, x, y)
		}
		return out
	}
	if got, want := order2D(Hilbert2DDecode), []uint64{0, 0, 1, 0, 1, 1, 0, 1, 0, 2, 0, 3, 1, 3, 1, 2, 2, 2, 2, 3, 3, 3, 3, 2, 3, 1, 2, 1, 2, 0, 3, 0}; !slices.Equal(got, want) {
		t.Errorf("Hilbert2D order = %v, expected %v", got, want)
	}
	if got, want := order2D(Moore2DDecode), []uint64{1, 0, 0, 0, 0, 1, 1, 1, 1, 2, 0, 2, 0, 3, 1, 3, 2, 3, 3, 3, 3, 2, 2, 2, 2, 1, 3, 1, 3, 0, 2, 0}; !slices.Equal(got, want) {
		t.Errorf("Moore2D order = %v, expected %v", got, want)
	}
	var order3D []uint64
	for d := range uint64(8) {
		x, y, z := Hilbert3DDecode(1, d)
		order3D = append(order3D, x, y, z)
	}
	if want := []uint64{0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0}; !slices.Equal(order3D, want) {
		t.Errorf("Hilbert3D order = %v, expected %v", order3D, want)
	}

	const x2, y2 = 0x89abcdef, 0x01234567
	const x3, y3, z3 = 0x1abcde, 0x012345, 0x1fffff
	for _, tc := range []struct {
		name      string
		got, want uint64
	}{
		{"Morton2DEncode", Morton2DEncode(32, x2, y2), 0x40434c4f70737c7f},
		{"Morton3DEncode", Morton3DEncode(21, x3, y3, z3), 0x5b2eb3db76be5bee},
		{"Hilbert2DEncode", Hilbert2DEncode(32, x2, y2), 0xeac2e2604ae268c0},
		{"Hilbert3DEncode", Hilbert3DEncode(21, x3, y3, z3), 0x6c5654ed92c6a472},
		{"HilbertEncode", HilbertEncode(16, []uint64{0x1234, 0xabcd, 0x0f0f, 0xffff}), 0x6d2612e9aaee5a2d},
		{"Moore2DEncode", Moore2DEncode(32, x2, y2), 0xff179d359f9d3d15},
		{"Gray2DEncode", Gray2DEncode(32, x2, y2), 0x7f82778a5fa257aa},
		{"Gray3DEncode", Gray3DEncode(21, x3, y3, z3), 0x6dcb22925b2b92b4},
		{"MortonToHilbert2D", MortonToHilbert2D(32, Morton2DEncode(32, x2, y2)), 0xeac2e2604ae268c0},
		{"MortonEncode", MortonEncode(16, []uint64{0x1234, 0xabcd, 0x0f0f, 0xffff}), 0xa8a9ecfeaa99efce},
		{"GrayEncode", GrayEncode(16, []uint64{0x1234, 0xabcd, 0x0f0f, 0xffff}), 0xcf3148ab33114a8b},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %#x, expected %#x", tc.name, tc.got, tc.want)
		}
	}
}

func TestTooManyBitsPanics(t *testing.T) {
	for name, f := range map[string]func(){
		"Morton2DEncode":    func() { Morton2DEncode(33, 0, 0) },
		"Morton2DDecode":    func() { Morton2DDecode(33, 0) },
		"Morton3DEncode":    func() { Morton3DEncode(22, 0, 0, 0) },
		"Morton3DDecode":    func() { Morton3DDecode(22, 0) },
		"MortonToHilbert2D": func() { MortonToHilbert2D(33, 0) },
		"HilbertToMorton2D": func() { HilbertToMorton2D(33, 0) },
		"Hilbert2DEncode":   func() { Hilbert2DEncode(33, 0, 0) },
		"Hilbert2DDecode":   func() { Hilbert2DDecode(33, 0) },
		"Hilbert3DEncode":   func() { Hilbert3DEncode(22, 0, 0, 0) },
		"Hilbert3DDecode":   func() { Hilbert3DDecode(22, 0) },
		"HilbertEncode":     func() { HilbertEncode(22, make([]uint64, 3)) },
		"HilbertDecode":     func() { HilbertDecode(22, 0, make([]uint64, 3)) },
		"Moore2DEncode":     func() { Moore2DEncode(33, 0, 0) },
		"Moore2DDecode":     func() { Moore2DDecode(33, 0) },
		"Gray2DEncode":      func() { Gray2DEncode(33, 0, 0) },
		"Gray2DDecode":      func() { Gray2DDecode(33, 0) },
		"Gray3DEncode":      func() { Gray3DEncode(22, 0, 0, 0) },
		"Gray3DDecode":      func() { Gray3DDecode(22, 0) },
		"Morton2DNextInBox": func() { Morton2DNextInBox(33, 0, 0, 0, 0, 0) },
		"Morton3DPrevInBox": func() { Morton3DPrevInBox(22, 0, 0, 0, 0, 0, 0, 0) },
		"Hilbert2DRanges":   func() { Hilbert2DRanges(33, 0, 0, 0, 0, 0) },
		"Hilbert3DRanges":   func() { Hilbert3DRanges(22, 0, 0, 0, 0, 0, 0, 0) },
		"HilbertRanges":     func() { HilbertRanges(22, make([]uint64, 3), make([]uint64, 3), 0) },
		"MortonEncode":      func() { MortonEncode(13, make([]uint64, 5)) },
		"MortonDecode":      func() { MortonDecode(13, 0, make([]uint64, 5)) },
		"GrayEncode":        func() { GrayEncode(13, make([]uint64, 5)) },
		"GrayDecode":        func() { GrayDecode(13, 0, make([]uint64, 5)) },
		"Morton2DRanges":    func() { Morton2DRanges(33, 0, 0, 0, 0, 0) },
		"Morton3DRanges":    func() { Morton3DRanges(22, 0, 0, 0, 0, 0, 0, 0) },
		"MortonRanges":      func() { MortonRanges(13, make([]uint64, 5), make([]uint64, 5), 0) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", name)
				}
			}()
			f()
		}()
	}
}

// TestIgnoresHighBits checks coordinates and indices are masked to the bit
// width.
func TestIgnoresHighBits(t *testing.T) {
	const bits = 4
	const hi = 0xabc << 20
	type tc struct {
		name      string
		got, want uint64
	}
	x1, y1 := Morton2DDecode(bits, hi|0x5a)
	x2, y2 := Hilbert2DDecode(bits, hi|0x5a)
	x3, y3 := Moore2DDecode(bits, hi|0x5a)
	x4, y4 := Gray2DDecode(bits, hi|0x5a)
	x5, y5, z5 := Hilbert3DDecode(bits, hi|0x5a)
	x6, y6, z6 := Morton3DDecode(bits, hi|0x5a)
	x7, y7, z7 := Gray3DDecode(bits, hi|0x5a)
	p := make([]uint64, 2)
	HilbertDecode(bits, hi|0x5a, p)
	m4, m4Want := make([]uint64, 4), make([]uint64, 4)
	MortonDecode(bits, hi<<20|0x5a5a, m4)
	MortonDecode(bits, 0x5a5a, m4Want)
	g4, g4Want := make([]uint64, 4), make([]uint64, 4)
	GrayDecode(bits, hi<<20|0x5a5a, g4)
	GrayDecode(bits, 0x5a5a, g4Want)
	if !slices.Equal(m4, m4Want) || !slices.Equal(g4, g4Want) {
		t.Errorf("MortonDecode/GrayDecode with high bits = %v/%v, expected %v/%v", m4, g4, m4Want, g4Want)
	}
	cases := []tc{
		{"Morton2DEncode", Morton2DEncode(bits, hi|3, hi|5), Morton2DEncode(bits, 3, 5)},
		{"Morton3DEncode", Morton3DEncode(bits, hi|3, hi|5, hi|6), Morton3DEncode(bits, 3, 5, 6)},
		{"Hilbert2DEncode", Hilbert2DEncode(bits, hi|3, hi|5), Hilbert2DEncode(bits, 3, 5)},
		{"Hilbert3DEncode", Hilbert3DEncode(bits, hi|3, hi|5, hi|6), Hilbert3DEncode(bits, 3, 5, 6)},
		{"HilbertEncode", HilbertEncode(bits, []uint64{hi | 3, hi | 5}), HilbertEncode(bits, []uint64{3, 5})},
		{"MortonEncode", MortonEncode(bits, []uint64{hi | 3, hi | 5, hi | 6, hi | 9}), MortonEncode(bits, []uint64{3, 5, 6, 9})},
		{"GrayEncode", GrayEncode(bits, []uint64{hi | 3, hi | 5, hi | 6, hi | 9}), GrayEncode(bits, []uint64{3, 5, 6, 9})},
		{"Moore2DEncode", Moore2DEncode(bits, hi|3, hi|5), Moore2DEncode(bits, 3, 5)},
		{"Gray2DEncode", Gray2DEncode(bits, hi|3, hi|5), Gray2DEncode(bits, 3, 5)},
		{"Gray3DEncode", Gray3DEncode(bits, hi|3, hi|5, hi|6), Gray3DEncode(bits, 3, 5, 6)},
		{"MortonToHilbert2D", MortonToHilbert2D(bits, hi|0x5a), MortonToHilbert2D(bits, 0x5a)},
		{"HilbertToMorton2D", HilbertToMorton2D(bits, hi|0x5a), HilbertToMorton2D(bits, 0x5a)},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %#x, expected %#x", c.name, c.got, c.want)
		}
	}
	for _, c := range []struct {
		name string
		got  [3]uint64
		want [3]uint64
	}{
		{"Morton2DDecode", [3]uint64{x1, y1}, pack2(Morton2DDecode(bits, 0x5a))},
		{"Hilbert2DDecode", [3]uint64{x2, y2}, pack2(Hilbert2DDecode(bits, 0x5a))},
		{"Moore2DDecode", [3]uint64{x3, y3}, pack2(Moore2DDecode(bits, 0x5a))},
		{"Gray2DDecode", [3]uint64{x4, y4}, pack2(Gray2DDecode(bits, 0x5a))},
		{"HilbertDecode", [3]uint64{p[0], p[1]}, pack2(Hilbert2DDecode(bits, 0x5a))},
		{"Hilbert3DDecode", [3]uint64{x5, y5, z5}, pack3(Hilbert3DDecode(bits, 0x5a))},
		{"Morton3DDecode", [3]uint64{x6, y6, z6}, pack3(Morton3DDecode(bits, 0x5a))},
		{"Gray3DDecode", [3]uint64{x7, y7, z7}, pack3(Gray3DDecode(bits, 0x5a))},
	} {
		if c.got != c.want {
			t.Errorf("%s = %v, expected %v", c.name, c.got, c.want)
		}
	}
}

func pack2(x, y uint64) [3]uint64 { return [3]uint64{x, y} }

func pack3(x, y, z uint64) [3]uint64 { return [3]uint64{x, y, z} }
