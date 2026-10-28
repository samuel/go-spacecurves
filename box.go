package spacecurves

import (
	"cmp"
	"fmt"
	"slices"
)

// Range is an inclusive range of curve indices.
type Range struct {
	First, Last uint
}

// Morton2DNextInBox returns the smallest index on a 2D Z-order curve of the
// specified number of bits per axis that is at least d and whose point lies
// in the box minX <= x <= maxX, minY <= y <= maxY. It returns false if there
// is none. When d is outside the box this is the BIGMIN of Tropf and Herzog,
// which lets a scan in Morton order skip to the next point in the box.
func Morton2DNextInBox(bits, d, minX, minY, maxX, maxY uint) (uint, bool) {
	minX, maxX, okX := clampAxis(bits, minX, maxX)
	minY, maxY, okY := clampAxis(bits, minY, maxY)
	if !okX || !okY || d > lowMask(2*bits) {
		return 0, false
	}
	if x, y := Morton2DDecode(bits, d); minX <= x && x <= maxX && minY <= y && y <= maxY {
		return d, true
	}
	return mortonBigMin(2, bits, d, Morton2DEncode(bits, minX, minY), Morton2DEncode(bits, maxX, maxY))
}

// Morton2DPrevInBox returns the largest index on a 2D Z-order curve of the
// specified number of bits per axis that is at most d and whose point lies
// in the box minX <= x <= maxX, minY <= y <= maxY. It returns false if there
// is none. When d is outside the box this is the LITMAX of Tropf and Herzog.
func Morton2DPrevInBox(bits, d, minX, minY, maxX, maxY uint) (uint, bool) {
	minX, maxX, okX := clampAxis(bits, minX, maxX)
	minY, maxY, okY := clampAxis(bits, minY, maxY)
	if !okX || !okY {
		return 0, false
	}
	d = min(d, lowMask(2*bits))
	if x, y := Morton2DDecode(bits, d); minX <= x && x <= maxX && minY <= y && y <= maxY {
		return d, true
	}
	return mortonLitMax(2, bits, d, Morton2DEncode(bits, minX, minY), Morton2DEncode(bits, maxX, maxY))
}

// Morton3DNextInBox is the 3D equivalent of Morton2DNextInBox.
func Morton3DNextInBox(bits, d, minX, minY, minZ, maxX, maxY, maxZ uint) (uint, bool) {
	minX, maxX, okX := clampAxis(bits, minX, maxX)
	minY, maxY, okY := clampAxis(bits, minY, maxY)
	minZ, maxZ, okZ := clampAxis(bits, minZ, maxZ)
	if !okX || !okY || !okZ || d > lowMask(3*bits) {
		return 0, false
	}
	if x, y, z := Morton3DDecode(bits, d); minX <= x && x <= maxX && minY <= y && y <= maxY && minZ <= z && z <= maxZ {
		return d, true
	}
	return mortonBigMin(3, bits, d, Morton3DEncode(bits, minX, minY, minZ), Morton3DEncode(bits, maxX, maxY, maxZ))
}

// Morton3DPrevInBox is the 3D equivalent of Morton2DPrevInBox.
func Morton3DPrevInBox(bits, d, minX, minY, minZ, maxX, maxY, maxZ uint) (uint, bool) {
	minX, maxX, okX := clampAxis(bits, minX, maxX)
	minY, maxY, okY := clampAxis(bits, minY, maxY)
	minZ, maxZ, okZ := clampAxis(bits, minZ, maxZ)
	if !okX || !okY || !okZ {
		return 0, false
	}
	d = min(d, lowMask(3*bits))
	if x, y, z := Morton3DDecode(bits, d); minX <= x && x <= maxX && minY <= y && y <= maxY && minZ <= z && z <= maxZ {
		return d, true
	}
	return mortonLitMax(3, bits, d, Morton3DEncode(bits, minX, minY, minZ), Morton3DEncode(bits, maxX, maxY, maxZ))
}

// clampAxis clamps an inclusive axis range to the grid. It returns false if
// the range is empty or entirely outside the grid.
func clampAxis(bits, lo, hi uint) (uint, uint, bool) {
	top := lowMask(bits)
	if lo > hi || lo > top {
		return 0, 0, false
	}
	return lo, min(hi, top), true
}

// mortonAxisMasks returns, for each axis, the mask of its bits in a Morton
// code of dims axes.
func mortonAxisMasks(dims, bits uint) [3]uint {
	var masks [3]uint
	for i := range dims * bits {
		masks[i%dims] |= 1 << i
	}
	return masks
}

// mortonBigMin returns the smallest Morton code greater than d inside the
// box with corner codes zmin and zmax. d must be outside the box.
func mortonBigMin(dims, bits, d, zmin, zmax uint) (uint, bool) {
	masks := mortonAxisMasks(dims, bits)
	var bigmin uint
	var found bool
	for i := dims*bits - 1; i < dims*bits; i-- {
		bit := uint(1) << i
		// The lower bits of the same axis as bit.
		lower := masks[i%dims] & (bit - 1)
		dBit, minBit, maxBit := d&bit != 0, zmin&bit != 0, zmax&bit != 0
		switch {
		case !dBit && !minBit && maxBit:
			bigmin, found = (zmin|bit)&^lower, true
			zmax = zmax&^bit | lower
		case !dBit && minBit:
			return zmin, true
		case dBit && !maxBit:
			return bigmin, found
		case dBit && !minBit && maxBit:
			zmin = (zmin | bit) &^ lower
		}
	}
	return bigmin, found
}

// mortonLitMax returns the largest Morton code less than d inside the box
// with corner codes zmin and zmax. d must be outside the box.
func mortonLitMax(dims, bits, d, zmin, zmax uint) (uint, bool) {
	masks := mortonAxisMasks(dims, bits)
	var litmax uint
	var found bool
	for i := dims*bits - 1; i < dims*bits; i-- {
		bit := uint(1) << i
		lower := masks[i%dims] & (bit - 1)
		dBit, minBit, maxBit := d&bit != 0, zmin&bit != 0, zmax&bit != 0
		switch {
		case !dBit && !minBit && maxBit:
			zmax = zmax&^bit | lower
		case !dBit && minBit:
			return litmax, found
		case dBit && !maxBit:
			return zmax, true
		case dBit && !minBit && maxBit:
			litmax, found = zmax&^bit|lower, true
			zmin = (zmin | bit) &^ lower
		}
	}
	return litmax, found
}

// Hilbert2DRanges returns the sorted index ranges of the Hilbert2DEncode
// curve that cover the box minX <= x <= maxX, minY <= y <= maxY. Adjacent
// ranges are merged. If maxRanges > 0, at most maxRanges ranges are returned:
// they still cover every point in the box but may include points outside it.
// Without a cap, the number of ranges and the work grow with the perimeter
// of the box.
func Hilbert2DRanges(bits, minX, minY, maxX, maxY uint, maxRanges int) []Range {
	return hilbertRanges(bits, []uint{minX, minY}, []uint{maxX, maxY}, maxRanges, func(d uint, p []uint) {
		p[0], p[1] = Hilbert2DDecode(bits, d)
	})
}

// Hilbert3DRanges is the 3D equivalent of Hilbert2DRanges for the
// Hilbert3DEncode curve.
func Hilbert3DRanges(bits, minX, minY, minZ, maxX, maxY, maxZ uint, maxRanges int) []Range {
	return hilbertRanges(bits, []uint{minX, minY, minZ}, []uint{maxX, maxY, maxZ}, maxRanges, func(d uint, p []uint) {
		p[0], p[1], p[2] = Hilbert3DDecode(bits, d)
	})
}

// HilbertRanges is the N-dimensional equivalent of Hilbert2DRanges for the
// HilbertEncode curve. The box is lo[i] <= point[i] <= hi[i]. Each refinement
// step visits 2^len(lo) children per partially covered cell, so it is meant
// for low dimensions. It panics if lo and hi differ in length or if
// len(lo)*bits exceeds the width of uint.
func HilbertRanges(bits uint, lo, hi []uint, maxRanges int) []Range {
	if len(lo) != len(hi) {
		panic(fmt.Sprintf("spacecurves: box corners have %d and %d dimensions", len(lo), len(hi)))
	}
	checkHilbertDims(bits, len(lo))
	return hilbertRanges(bits, lo, hi, maxRanges, func(d uint, p []uint) {
		HilbertDecode(bits, d, p)
	})
}

// hilbertRanges finds the ranges covering a box by refining the curve one
// level at a time. Every aligned block of 2^(dims*k) indices fills an aligned
// cube of side 2^k, so a block is kept whole if its cube is inside the box,
// dropped if disjoint, and split into its 2^dims children otherwise.
func hilbertRanges(bits uint, lo, hi []uint, maxRanges int, decode func(d uint, p []uint)) []Range {
	dims := uint(len(lo))
	if dims == 0 {
		return []Range{{0, 0}}
	}
	if dims >= uintBits {
		panic(fmt.Sprintf("spacecurves: Hilbert ranges support fewer than %d dimensions", uintBits))
	}
	lo, hi = slices.Clone(lo), slices.Clone(hi)
	for i := range lo {
		var ok bool
		if lo[i], hi[i], ok = clampAxis(bits, lo[i], hi[i]); !ok {
			return nil
		}
	}

	type block struct {
		first, last uint
		full        bool
	}
	p := make([]uint, dims)
	// classify reports whether the cube holding the block at base with side
	// 2^k intersects the box and whether it lies entirely inside it.
	classify := func(base, k uint) (bool, bool) {
		decode(base, p)
		side := lowMask(k)
		full := true
		for i, v := range p {
			start := v &^ side
			end := start + side
			if end < lo[i] || start > hi[i] {
				return false, false
			}
			if start < lo[i] || end > hi[i] {
				full = false
			}
		}
		return true, full
	}
	toRanges := func(blocks []block) []Range {
		var r []Range
		for _, b := range blocks {
			if n := len(r); n > 0 && r[n-1].Last+1 == b.first {
				r[n-1].Last = b.last
			} else {
				r = append(r, Range{b.first, b.last})
			}
		}
		return r
	}

	_, full := classify(0, bits)
	blocks := []block{{0, lowMask(dims * bits), full}}
	for k := bits; k > 0; k-- {
		if !slices.ContainsFunc(blocks, func(b block) bool { return !b.full }) {
			break
		}
		childBits := dims * (k - 1)
		var next []block
		for _, b := range blocks {
			if b.full {
				next = append(next, b)
				continue
			}
			for i := range uint(1) << dims {
				base := b.first + i<<childBits
				if in, full := classify(base, k-1); in {
					next = append(next, block{base, base + lowMask(childBits), full})
				}
			}
		}
		if maxRanges > 0 {
			if r := toRanges(next); len(r) > maxRanges {
				return mergeSmallestGaps(r, maxRanges)
			}
		}
		blocks = next
	}
	return toRanges(blocks)
}

// mergeSmallestGaps reduces sorted ranges to n by closing the smallest gaps
// between them, which adds the fewest extra indices.
func mergeSmallestGaps(r []Range, n int) []Range {
	gaps := make([]int, len(r)-1)
	for i := range gaps {
		gaps[i] = i
	}
	gapSize := func(i int) uint { return r[i+1].First - r[i].Last }
	slices.SortStableFunc(gaps, func(a, b int) int { return cmp.Compare(gapSize(a), gapSize(b)) })
	closed := make([]bool, len(r)-1)
	for _, g := range gaps[:len(r)-n] {
		closed[g] = true
	}
	out := make([]Range, 0, n)
	out = append(out, r[0])
	for i := 1; i < len(r); i++ {
		if closed[i-1] {
			out[len(out)-1].Last = r[i].Last
		} else {
			out = append(out, r[i])
		}
	}
	return out
}
