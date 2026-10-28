package spacecurves

import "testing"

// curve2D and curve3D list the 2D and 3D curves with their maximum bits.
var curve2D = []struct {
	name   string
	encode func(uint, uint64, uint64) uint64
	decode func(uint, uint64) (uint64, uint64)
}{
	{"Morton", Morton2DEncode, Morton2DDecode},
	{"Hilbert", Hilbert2DEncode, Hilbert2DDecode},
	{"Moore", Moore2DEncode, Moore2DDecode},
	{"Gray", Gray2DEncode, Gray2DDecode},
	{"HilbertND", func(bits uint, x, y uint64) uint64 { return HilbertEncode(bits, []uint64{x, y}) }, func(bits uint, d uint64) (uint64, uint64) {
		var p [2]uint64
		HilbertDecode(bits, d, p[:])
		return p[0], p[1]
	}},
}

var curve3D = []struct {
	name   string
	encode func(uint, uint64, uint64, uint64) uint64
	decode func(uint, uint64) (uint64, uint64, uint64)
}{
	{"Morton", Morton3DEncode, Morton3DDecode},
	{"Hilbert", Hilbert3DEncode, Hilbert3DDecode},
	{"Gray", Gray3DEncode, Gray3DDecode},
}

func FuzzCurve2D(f *testing.F) {
	f.Add(uint8(32), uint64(0x89abcdef), uint64(0x01234567), uint64(0xfedcba9876543210))
	f.Add(uint8(0), uint64(1), uint64(1), uint64(1))
	f.Fuzz(func(t *testing.T, b uint8, x, y, d uint64) {
		bits := uint(b) % 33
		m := lowMask(bits)
		for _, c := range curve2D {
			if x2, y2 := c.decode(bits, c.encode(bits, x, y)); x2 != x&m || y2 != y&m {
				t.Fatalf("%s bits=%d: decode(encode(%#x, %#x)) = (%#x, %#x)", c.name, bits, x, y, x2, y2)
			}
			if x2, y2 := c.decode(bits, d); c.encode(bits, x2, y2) != d&lowMask(2*bits) {
				t.Fatalf("%s bits=%d: encode(decode(%#x)) = %#x", c.name, bits, d, c.encode(bits, x2, y2))
			}
		}
		if got, want := MortonToHilbert2D(bits, Morton2DEncode(bits, x, y)), Hilbert2DEncode(bits, x, y); got != want {
			t.Fatalf("bits=%d: MortonToHilbert2D = %#x, Hilbert2DEncode = %#x", bits, got, want)
		}
		if got := HilbertToMorton2D(bits, MortonToHilbert2D(bits, d)); got != d&lowMask(2*bits) {
			t.Fatalf("bits=%d: HilbertToMorton2D(MortonToHilbert2D(%#x)) = %#x", bits, d, got)
		}
	})
}

func FuzzCurve3D(f *testing.F) {
	f.Add(uint8(21), uint64(0x1abcde), uint64(0x012345), uint64(0x1fffff), uint64(0x7edcba9876543210))
	f.Fuzz(func(t *testing.T, b uint8, x, y, z, d uint64) {
		bits := uint(b) % 22
		m := lowMask(bits)
		for _, c := range curve3D {
			if x2, y2, z2 := c.decode(bits, c.encode(bits, x, y, z)); x2 != x&m || y2 != y&m || z2 != z&m {
				t.Fatalf("%s bits=%d: decode(encode(%#x, %#x, %#x)) = (%#x, %#x, %#x)", c.name, bits, x, y, z, x2, y2, z2)
			}
			if x2, y2, z2 := c.decode(bits, d); c.encode(bits, x2, y2, z2) != d&lowMask(3*bits) {
				t.Fatalf("%s bits=%d: encode(decode(%#x)) = %#x", c.name, bits, d, c.encode(bits, x2, y2, z2))
			}
		}
	})
}

func FuzzHilbertND(f *testing.F) {
	f.Add(uint8(4), uint8(16), uint64(0x0123456789abcdef))
	f.Fuzz(func(t *testing.T, dimsB, bitsB uint8, d uint64) {
		dims := int(dimsB)%8 + 1
		bits := uint(bitsB) % (64/uint(dims) + 1)
		p := make([]uint64, dims)
		HilbertDecode(bits, d, p)
		if got := HilbertEncode(bits, p); got != d&lowMask(uint(dims)*bits) {
			t.Fatalf("dims=%d bits=%d: HilbertEncode(HilbertDecode(%#x)) = %#x", dims, bits, d, got)
		}
		// The next index must be adjacent unless d is the last one.
		if next := d&lowMask(uint(dims)*bits) + 1; next <= lowMask(uint(dims)*bits) && next != 0 {
			q := make([]uint64, dims)
			HilbertDecode(bits, next, q)
			var sum uint64
			for i := range p {
				sum += dist(p[i], q[i])
			}
			if sum != 1 {
				t.Fatalf("dims=%d bits=%d: indices %#x and %#x are not adjacent: %v, %v", dims, bits, next-1, next, p, q)
			}
		}
	})
}

// fuzzBox returns a small grid and a box in it, possibly extending past the
// grid or inverted.
func fuzzBox(b uint8, lo, hi uint16) (uint, [2]uint64, [2]uint64) {
	bits := uint(b)%4 + 1
	n := uint64(1)<<bits + 2
	return bits, [2]uint64{uint64(lo&0xff) % n, uint64(lo>>8) % n}, [2]uint64{uint64(hi&0xff) % n, uint64(hi>>8) % n}
}

func FuzzMorton2DInBox(f *testing.F) {
	f.Add(uint8(3), uint16(0x0201), uint16(0x0605), uint16(17))
	f.Fuzz(func(t *testing.T, b uint8, lo, hi, d16 uint16) {
		bits, mn, mx := fuzzBox(b, lo, hi)
		total := uint64(1) << (2 * bits)
		d := uint64(d16) % (total + 2)
		in := func(e uint64) bool {
			x, y := Morton2DDecode(bits, e)
			return inAxis(x, mn[0], mx[0]) && inAxis(y, mn[1], mx[1])
		}
		wantNext, wantNextOK := uint64(0), false
		for e := d; e < total; e++ {
			if in(e) {
				wantNext, wantNextOK = e, true
				break
			}
		}
		wantPrev, wantPrevOK := uint64(0), false
		for e := int64(min(d, total-1)); e >= 0; e-- {
			if in(uint64(e)) {
				wantPrev, wantPrevOK = uint64(e), true
				break
			}
		}
		if got, ok := Morton2DNextInBox(bits, d, mn[0], mn[1], mx[0], mx[1]); got != wantNext || ok != wantNextOK {
			t.Fatalf("Morton2DNextInBox(%d, %v..%v) = %d, %v, expected %d, %v", d, mn, mx, got, ok, wantNext, wantNextOK)
		}
		if got, ok := Morton2DPrevInBox(bits, d, mn[0], mn[1], mx[0], mx[1]); got != wantPrev || ok != wantPrevOK {
			t.Fatalf("Morton2DPrevInBox(%d, %v..%v) = %d, %v, expected %d, %v", d, mn, mx, got, ok, wantPrev, wantPrevOK)
		}
	})
}

func FuzzHilbert2DRanges(f *testing.F) {
	f.Add(uint8(3), uint16(0x0201), uint16(0x0605), uint8(0))
	f.Add(uint8(4), uint16(0x0003), uint16(0x0e0c), uint8(3))
	f.Fuzz(func(t *testing.T, b uint8, lo, hi uint16, capB uint8) {
		bits, mn, mx := fuzzBox(b, lo, hi)
		maxRanges := int(capB % 8)
		ranges := Hilbert2DRanges(bits, mn[0], mn[1], mx[0], mx[1], maxRanges)
		checkRanges(t, ranges, 1<<(2*bits), maxRanges, func(d uint64) bool {
			x, y := Hilbert2DDecode(bits, d)
			return inAxis(x, mn[0], mx[0]) && inAxis(y, mn[1], mx[1])
		})
	})
}
