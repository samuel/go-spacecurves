package spacecurves

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// randomAxis returns a random inclusive axis range that sometimes extends
// past the grid or is inverted.
func randomAxis(rnd *rand.Rand, bits uint) (uint, uint) {
	n := uint64(1)<<bits + 2
	return uint(rnd.Uint64N(n)), uint(rnd.Uint64N(n))
}

func inAxis(v, lo, hi uint) bool { return lo <= v && v <= hi }

func TestMorton2DInBox(t *testing.T) {
	const bits = 4
	const total = 1 << (2 * bits)
	rnd := rand.New(rand.NewPCG(13, 14))
	for range 300 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		var in [total]bool
		for d := range uint(total) {
			x, y := Morton2DDecode(bits, d)
			in[d] = inAxis(x, minX, maxX) && inAxis(y, minY, maxY)
		}
		for d := range uint(total + 2) {
			wantNext, wantNextOK := uint(0), false
			for e := d; e < total; e++ {
				if in[e] {
					wantNext, wantNextOK = e, true
					break
				}
			}
			wantPrev, wantPrevOK := uint(0), false
			for e := int(min(d, total-1)); e >= 0; e-- {
				if in[e] {
					wantPrev, wantPrevOK = uint(e), true
					break
				}
			}
			if got, ok := Morton2DNextInBox(bits, d, minX, minY, maxX, maxY); got != wantNext || ok != wantNextOK {
				t.Fatalf("Morton2DNextInBox(%d, box x %d..%d y %d..%d) = %d, %v, expected %d, %v", d, minX, maxX, minY, maxY, got, ok, wantNext, wantNextOK)
			}
			if got, ok := Morton2DPrevInBox(bits, d, minX, minY, maxX, maxY); got != wantPrev || ok != wantPrevOK {
				t.Fatalf("Morton2DPrevInBox(%d, box x %d..%d y %d..%d) = %d, %v, expected %d, %v", d, minX, maxX, minY, maxY, got, ok, wantPrev, wantPrevOK)
			}
		}
	}
}

func TestMorton3DInBox(t *testing.T) {
	const bits = 3
	const total = 1 << (3 * bits)
	rnd := rand.New(rand.NewPCG(15, 16))
	for range 200 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		minZ, maxZ := randomAxis(rnd, bits)
		var in [total]bool
		for d := range uint(total) {
			x, y, z := Morton3DDecode(bits, d)
			in[d] = inAxis(x, minX, maxX) && inAxis(y, minY, maxY) && inAxis(z, minZ, maxZ)
		}
		for d := range uint(total + 2) {
			wantNext, wantNextOK := uint(0), false
			for e := d; e < total; e++ {
				if in[e] {
					wantNext, wantNextOK = e, true
					break
				}
			}
			wantPrev, wantPrevOK := uint(0), false
			for e := int(min(d, total-1)); e >= 0; e-- {
				if in[e] {
					wantPrev, wantPrevOK = uint(e), true
					break
				}
			}
			if got, ok := Morton3DNextInBox(bits, d, minX, minY, minZ, maxX, maxY, maxZ); got != wantNext || ok != wantNextOK {
				t.Fatalf("Morton3DNextInBox(%d) = %d, %v, expected %d, %v", d, got, ok, wantNext, wantNextOK)
			}
			if got, ok := Morton3DPrevInBox(bits, d, minX, minY, minZ, maxX, maxY, maxZ); got != wantPrev || ok != wantPrevOK {
				t.Fatalf("Morton3DPrevInBox(%d) = %d, %v, expected %d, %v", d, got, ok, wantPrev, wantPrevOK)
			}
		}
	}
}

func TestMorton2DInBoxMaxBits(t *testing.T) {
	const bits = uintBits / 2
	rnd := rand.New(rand.NewPCG(17, 18))
	for range 1000 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		d := uint(rnd.Uint64())
		next, ok := Morton2DNextInBox(bits, d, minX, minY, maxX, maxY)
		if !ok {
			continue
		}
		if x, y := Morton2DDecode(bits, next); next < d || !inAxis(x, minX, maxX) || !inAxis(y, minY, maxY) {
			t.Fatalf("Morton2DNextInBox(%#x) = %#x, which is before d or outside the box", d, next)
		}
		// Nothing in the box lies in [d, next).
		if next > d {
			if prev, ok := Morton2DPrevInBox(bits, next-1, minX, minY, maxX, maxY); ok && prev >= d {
				t.Fatalf("Morton2DNextInBox(%#x) = %#x skips %#x", d, next, prev)
			}
		}
	}
}

// checkRanges verifies ranges are sorted, separated by gaps, and cover
// exactly the indices in the box, or a superset of them when capped.
func checkRanges(t *testing.T, ranges []Range, total uint, maxRanges int, inBox func(d uint) bool) {
	t.Helper()
	if maxRanges > 0 && len(ranges) > maxRanges {
		t.Fatalf("Got %d ranges, cap is %d", len(ranges), maxRanges)
	}
	for i, r := range ranges {
		if r.First > r.Last || r.Last >= total {
			t.Fatalf("Invalid range %v", r)
		}
		if i > 0 && ranges[i-1].Last+1 >= r.First {
			t.Fatalf("Ranges %v and %v overlap or touch", ranges[i-1], r)
		}
	}
	for d := range total {
		_, covered := slices.BinarySearchFunc(ranges, d, func(r Range, d uint) int {
			switch {
			case r.Last < d:
				return -1
			case r.First > d:
				return 1
			}
			return 0
		})
		if in := inBox(d); in && !covered || !in && covered && maxRanges <= 0 {
			t.Fatalf("Index %d: in box %v, covered %v, ranges %v", d, in, covered, ranges)
		}
	}
}

func TestMergeSmallestGaps(t *testing.T) {
	r := []Range{{0, 1}, {5, 6}, {8, 9}, {20, 21}, {23, 24}}
	got := mergeSmallestGaps(r, 3)
	want := []Range{{0, 1}, {5, 9}, {20, 24}}
	if !slices.Equal(got, want) {
		t.Fatalf("mergeSmallestGaps = %v, expected %v", got, want)
	}
}

func TestHilbert2DRanges(t *testing.T) {
	const bits = 4
	rnd := rand.New(rand.NewPCG(19, 20))
	for range 300 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		inBox := func(d uint) bool {
			x, y := Hilbert2DDecode(bits, d)
			return inAxis(x, minX, maxX) && inAxis(y, minY, maxY)
		}
		for _, maxRanges := range []int{0, 1, 2, 3, 5} {
			ranges := Hilbert2DRanges(bits, minX, minY, maxX, maxY, maxRanges)
			checkRanges(t, ranges, 1<<(2*bits), maxRanges, inBox)
			nd := HilbertRanges(bits, []uint{minX, minY}, []uint{maxX, maxY}, maxRanges)
			if !slices.Equal(ranges, nd) {
				t.Fatalf("Hilbert2DRanges = %v, HilbertRanges = %v", ranges, nd)
			}
		}
	}
}

func TestHilbert3DRanges(t *testing.T) {
	const bits = 3
	rnd := rand.New(rand.NewPCG(21, 22))
	for range 200 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		minZ, maxZ := randomAxis(rnd, bits)
		inBox := func(d uint) bool {
			x, y, z := Hilbert3DDecode(bits, d)
			return inAxis(x, minX, maxX) && inAxis(y, minY, maxY) && inAxis(z, minZ, maxZ)
		}
		for _, maxRanges := range []int{0, 1, 4} {
			checkRanges(t, Hilbert3DRanges(bits, minX, minY, minZ, maxX, maxY, maxZ, maxRanges), 1<<(3*bits), maxRanges, inBox)
		}
	}
}

func TestHilbertRangesND(t *testing.T) {
	rnd := rand.New(rand.NewPCG(23, 24))
	for dims := 1; dims <= 4; dims++ {
		nbits := uint(8 / dims)
		lo, hi := make([]uint, dims), make([]uint, dims)
		p := make([]uint, dims)
		for range 50 {
			for i := range dims {
				lo[i], hi[i] = randomAxis(rnd, nbits)
			}
			inBox := func(d uint) bool {
				HilbertDecode(nbits, d, p)
				for i := range p {
					if !inAxis(p[i], lo[i], hi[i]) {
						return false
					}
				}
				return true
			}
			for _, maxRanges := range []int{0, 3} {
				checkRanges(t, HilbertRanges(nbits, lo, hi, maxRanges), 1<<(nbits*uint(dims)), maxRanges, inBox)
			}
		}
	}
}

func TestHilbert2DRangesMaxBits(t *testing.T) {
	const bits = uintBits / 2
	const top = math.MaxUint >> (uintBits - bits)
	if got := Hilbert2DRanges(bits, 0, 0, top, top, 0); !slices.Equal(got, []Range{{0, math.MaxUint}}) {
		t.Fatalf("Full grid ranges = %v", got)
	}
	rnd := rand.New(rand.NewPCG(25, 26))
	for range 20 {
		minX, maxX := randomAxis(rnd, bits)
		minY, maxY := randomAxis(rnd, bits)
		ranges := Hilbert2DRanges(bits, minX, minY, maxX, maxY, 16)
		if len(ranges) > 16 {
			t.Fatalf("Got %d ranges, cap is 16", len(ranges))
		}
		if minX > maxX || minY > maxY {
			if ranges != nil {
				t.Fatalf("Empty box returned %v", ranges)
			}
			continue
		}
		maxX, maxY = min(maxX, top), min(maxY, top)
		points := [][2]uint{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}}
		for range 100 {
			points = append(points, [2]uint{minX + uint(rnd.Uint64N(uint64(maxX-minX)+1)), minY + uint(rnd.Uint64N(uint64(maxY-minY)+1))})
		}
		for _, pt := range points {
			d := Hilbert2DEncode(bits, pt[0], pt[1])
			if !slices.ContainsFunc(ranges, func(r Range) bool { return r.First <= d && d <= r.Last }) {
				t.Fatalf("Point %v (index %#x) not covered by %v", pt, d, ranges)
			}
		}
	}
}

func TestHilbertRangesPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	HilbertRanges(4, []uint{0, 0}, []uint{1}, 0)
}

func BenchmarkMorton2DNextInBox_16(b *testing.B) {
	for b.Loop() {
		Morton2DNextInBox(16, 0x12345678, 0x1000, 0x2000, 0x1fff, 0x2fff)
	}
}

func BenchmarkHilbert2DRanges_16(b *testing.B) {
	for b.Loop() {
		Hilbert2DRanges(16, 1000, 2000, 1999, 2999, 0)
	}
}

func BenchmarkHilbert2DRanges_16_Cap64(b *testing.B) {
	for b.Loop() {
		Hilbert2DRanges(16, 1000, 2000, 1999, 2999, 64)
	}
}
