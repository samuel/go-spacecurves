package spacecurves

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

func TestGray2DKnownValues(t *testing.T) {
	want := [][2]uint64{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	for d, p := range want {
		if x, y := Gray2DDecode(1, uint64(d)); x != p[0] || y != p[1] {
			t.Fatalf("Gray2DDecode(1, %d) = (%d, %d), expected (%d, %d)", d, x, y, p[0], p[1])
		}
	}
}

func TestGray2D(t *testing.T) {
	for nbits := range uint(9) {
		total := uint64(1) << (2 * nbits)
		lx, ly := Gray2DDecode(nbits, 0)
		for d := range total {
			x, y := Gray2DDecode(nbits, d)
			if x > lowMask(nbits) || y > lowMask(nbits) {
				t.Fatalf("bits=%d: index %d decodes to out of range point (%d, %d)", nbits, d, x, y)
			}
			if d2 := Gray2DEncode(nbits, x, y); d2 != d {
				t.Fatalf("bits=%d: Gray2DEncode(Gray2DDecode(%d)) = %d", nbits, d, d2)
			}
			if d > 0 && bits.OnesCount64(lx^x)+bits.OnesCount64(ly^y) != 1 {
				t.Fatalf("bits=%d: indices %d and %d map to (%d, %d) and (%d, %d), which differ in more than one bit", nbits, d-1, d, lx, ly, x, y)
			}
			lx, ly = x, y
		}
	}
}

func TestGray3D(t *testing.T) {
	for nbits := range uint(6) {
		total := uint64(1) << (3 * nbits)
		lx, ly, lz := Gray3DDecode(nbits, 0)
		for d := range total {
			x, y, z := Gray3DDecode(nbits, d)
			if x > lowMask(nbits) || y > lowMask(nbits) || z > lowMask(nbits) {
				t.Fatalf("bits=%d: index %d decodes to out of range point (%d, %d, %d)", nbits, d, x, y, z)
			}
			if d2 := Gray3DEncode(nbits, x, y, z); d2 != d {
				t.Fatalf("bits=%d: Gray3DEncode(Gray3DDecode(%d)) = %d", nbits, d, d2)
			}
			if d > 0 && bits.OnesCount64(lx^x)+bits.OnesCount64(ly^y)+bits.OnesCount64(lz^z) != 1 {
				t.Fatalf("bits=%d: indices %d and %d map to (%d, %d, %d) and (%d, %d, %d), which differ in more than one bit", nbits, d-1, d, lx, ly, lz, x, y, z)
			}
			lx, ly, lz = x, y, z
		}
	}
}

func TestGrayMaxBits(t *testing.T) {
	rnd := rand.New(rand.NewPCG(11, 12))
	const bits2 = 32
	for _, d := range edgeValues(2*bits2, rnd) {
		x, y := Gray2DDecode(bits2, d)
		if d2 := Gray2DEncode(bits2, x, y); d2 != d {
			t.Fatalf("Gray2DEncode(Gray2DDecode(%#x)) = %#x", d, d2)
		}
	}
	const bits3 = 21
	for _, d := range edgeValues(3*bits3, rnd) {
		x, y, z := Gray3DDecode(bits3, d)
		if d2 := Gray3DEncode(bits3, x, y, z); d2 != d {
			t.Fatalf("Gray3DEncode(Gray3DDecode(%#x)) = %#x", d, d2)
		}
	}
}

func BenchmarkGray2DEncode_16(b *testing.B) {
	for b.Loop() {
		Gray2DEncode(16, 0xabcd, 0x1234)
	}
}

func BenchmarkGray2DDecode_16(b *testing.B) {
	for b.Loop() {
		Gray2DDecode(16, 0x12345678)
	}
}

func BenchmarkGray3DEncode_10(b *testing.B) {
	for b.Loop() {
		Gray3DEncode(10, 0x2cd, 0x134, 0x3a5)
	}
}
