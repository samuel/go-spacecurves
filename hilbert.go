package spacecurves

import (
	"fmt"
	"math/bits"
)

const uintBits = bits.UintSize

// Hilbert2DEncode encodes a 2D x,y pair into a single value along
// a Hilbert space-filling curve of the specified number of bits per axis
// (at most 32 on 64-bit platforms).
func Hilbert2DEncode(bits, x, y uint) uint {
	d := uint(0)
	for s := uint(1<<bits) / 2; s > 0; s /= 2 {
		var rx, ry uint
		if x&s > 0 {
			rx = 1
		}
		if y&s > 0 {
			ry = 1
		}
		d += s * s * ((3 * rx) ^ ry)
		x, y = hilbertRot(s, x, y, rx, ry)
	}
	return d
}

// Hilbert2DDecode decodes a value along a Hilbert space-filling curve
// of the specified number of bits per axis into a 2D x,y pair.
func Hilbert2DDecode(bits, d uint) (uint, uint) {
	var x, y uint
	n := uint(1 << bits)
	for s := uint(1); s < n; s *= 2 {
		rx := 1 & (d / 2)
		ry := 1 & (d ^ rx)
		x, y = hilbertRot(s, x, y, rx, ry)
		x += s * rx
		y += s * ry
		d /= 4
	}
	return x, y
}

func hilbertRot(n, x, y, rx, ry uint) (uint, uint) {
	if ry == 0 {
		if rx == 1 {
			x = n - 1 - x
			y = n - 1 - y
		}
		x, y = y, x
	}
	return x, y
}

// HilbertEncode encodes an N-dimensional point into a single value along
// a Hilbert space-filling curve of the specified number of bits per axis.
// len(point)*bits must not exceed the width of uint. It panics otherwise.
// The point is not modified.
func HilbertEncode(bits uint, point []uint) uint {
	checkHilbertDims(bits, len(point))
	if bits == 0 {
		return 0
	}
	var buf [uintBits]uint
	x := buf[:len(point)]
	copy(x, point)
	hilbertAxesToTranspose(bits, x)
	return hilbertInterleave(bits, x)
}

// HilbertDecode decodes a value along a Hilbert space-filling curve of the
// specified number of bits per axis into point, whose length sets the number
// of dimensions. len(point)*bits must not exceed the width of uint. It panics
// otherwise.
func HilbertDecode(bits, d uint, point []uint) {
	checkHilbertDims(bits, len(point))
	hilbertDeinterleave(bits, d, point)
	hilbertTransposeToAxes(bits, point)
}

// Hilbert3DEncode encodes a 3D x,y,z point into a single value along
// a Hilbert space-filling curve of the specified number of bits per axis
// (at most 21 on 64-bit platforms).
func Hilbert3DEncode(bits, x, y, z uint) uint {
	p := [3]uint{x, y, z}
	hilbertAxesToTranspose(bits, p[:])
	return hilbertInterleave(bits, p[:])
}

// Hilbert3DDecode decodes a value along a Hilbert space-filling curve
// of the specified number of bits per axis into a 3D x,y,z point.
func Hilbert3DDecode(bits, d uint) (uint, uint, uint) {
	var p [3]uint
	hilbertDeinterleave(bits, d, p[:])
	hilbertTransposeToAxes(bits, p[:])
	return p[0], p[1], p[2]
}

func checkHilbertDims(nbits uint, dims int) {
	if dims > 0 && nbits > uintBits/uint(dims) {
		panic(fmt.Sprintf("spacecurves: %d dimensions of %d bits exceeds %d bits", dims, nbits, uintBits))
	}
}

// hilbertAxesToTranspose converts coordinates in place to the transposed
// Hilbert index, using the algorithm from J. Skilling, "Programming the
// Hilbert curve", AIP Conference Proceedings 707, 2004.
func hilbertAxesToTranspose(nbits uint, x []uint) {
	if nbits == 0 || len(x) == 0 {
		return
	}
	n := len(x)
	m := uint(1) << (nbits - 1)
	for q := m; q > 1; q >>= 1 {
		p := q - 1
		for i := range n {
			if x[i]&q != 0 {
				x[0] ^= p
			} else {
				t := (x[0] ^ x[i]) & p
				x[0] ^= t
				x[i] ^= t
			}
		}
	}
	for i := 1; i < n; i++ {
		x[i] ^= x[i-1]
	}
	var t uint
	for q := m; q > 1; q >>= 1 {
		if x[n-1]&q != 0 {
			t ^= q - 1
		}
	}
	for i := range n {
		x[i] ^= t
	}
}

// hilbertTransposeToAxes is the inverse of hilbertAxesToTranspose.
func hilbertTransposeToAxes(nbits uint, x []uint) {
	if nbits == 0 || len(x) == 0 {
		return
	}
	n := len(x)
	t := x[n-1] >> 1
	for i := n - 1; i > 0; i-- {
		x[i] ^= x[i-1]
	}
	x[0] ^= t
	// end is 0 when nbits is the width of uint, so q stops when it overflows.
	end := uint(2) << (nbits - 1)
	for q := uint(2); q != end; q <<= 1 {
		p := q - 1
		for i := n - 1; i >= 0; i-- {
			if x[i]&q != 0 {
				x[0] ^= p
			} else {
				t := (x[0] ^ x[i]) & p
				x[0] ^= t
				x[i] ^= t
			}
		}
	}
}

// hilbertInterleave packs the transposed index into a single value, taking
// the most significant bit of x[0] first.
func hilbertInterleave(nbits uint, x []uint) uint {
	var d uint
	for j := int(nbits) - 1; j >= 0; j-- {
		for _, v := range x {
			d = d<<1 | (v>>uint(j))&1
		}
	}
	return d
}

// hilbertDeinterleave is the inverse of hilbertInterleave.
func hilbertDeinterleave(nbits, d uint, x []uint) {
	clear(x)
	for j := range nbits {
		for i := len(x) - 1; i >= 0; i-- {
			x[i] |= (d & 1) << j
			d >>= 1
		}
	}
}
