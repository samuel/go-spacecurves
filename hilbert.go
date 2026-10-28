package spacecurves

// Hilbert2DEncode encodes a 2D x,y pair into a single value along
// a Hilbert space-filling curve of the specified number of bits per axis
// (at most 32).
func Hilbert2DEncode(bits uint, x, y uint64) uint64 {
	checkBits(bits, 2)
	var d uint64
	for s := uint64(1) << bits / 2; s > 0; s /= 2 {
		var rx, ry uint64
		if x&s != 0 {
			rx = 1
		}
		if y&s != 0 {
			ry = 1
		}
		d += s * s * ((3 * rx) ^ ry)
		x, y = hilbertRot(s, x, y, rx, ry)
	}
	return d
}

// Hilbert2DDecode decodes a value along a Hilbert space-filling curve
// of the specified number of bits per axis into a 2D x,y pair.
func Hilbert2DDecode(bits uint, d uint64) (uint64, uint64) {
	checkBits(bits, 2)
	var x, y uint64
	n := uint64(1) << bits
	for s := uint64(1); s < n; s *= 2 {
		rx := 1 & (d / 2)
		ry := 1 & (d ^ rx)
		x, y = hilbertRot(s, x, y, rx, ry)
		x += s * rx
		y += s * ry
		d /= 4
	}
	return x, y
}

// hilbertRot rotates and reflects a quadrant. Bits of x and y at or above n
// may be left wrapped, which later levels ignore.
func hilbertRot(n, x, y, rx, ry uint64) (uint64, uint64) {
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
// The point is not modified.
func HilbertEncode(bits uint, point []uint64) uint64 {
	checkBits(bits, len(point))
	if bits == 0 {
		return 0
	}
	var buf [64]uint64
	x := buf[:len(point)]
	m := lowMask(bits)
	for i, v := range point {
		x[i] = v & m
	}
	hilbertAxesToTranspose(bits, x)
	return hilbertInterleave(bits, x)
}

// HilbertDecode decodes a value along a Hilbert space-filling curve of the
// specified number of bits per axis into point, whose length sets the number
// of dimensions.
func HilbertDecode(bits uint, d uint64, point []uint64) {
	checkBits(bits, len(point))
	hilbertDeinterleave(bits, d, point)
	hilbertTransposeToAxes(bits, point)
}

// Hilbert3DEncode encodes a 3D x,y,z point into a single value along
// a Hilbert space-filling curve of the specified number of bits per axis
// (at most 21).
func Hilbert3DEncode(bits uint, x, y, z uint64) uint64 {
	checkBits(bits, 3)
	m := lowMask(bits)
	p := [3]uint64{x & m, y & m, z & m}
	hilbertAxesToTranspose(bits, p[:])
	return hilbertInterleave(bits, p[:])
}

// Hilbert3DDecode decodes a value along a Hilbert space-filling curve
// of the specified number of bits per axis into a 3D x,y,z point.
func Hilbert3DDecode(bits uint, d uint64) (uint64, uint64, uint64) {
	checkBits(bits, 3)
	var p [3]uint64
	hilbertDeinterleave(bits, d, p[:])
	hilbertTransposeToAxes(bits, p[:])
	return p[0], p[1], p[2]
}

// hilbertAxesToTranspose converts coordinates in place to the transposed
// Hilbert index, using the algorithm from J. Skilling, "Programming the
// Hilbert curve", AIP Conference Proceedings 707, 2004.
func hilbertAxesToTranspose(nbits uint, x []uint64) {
	if nbits == 0 || len(x) == 0 {
		return
	}
	n := len(x)
	m := uint64(1) << (nbits - 1)
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
	var t uint64
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
func hilbertTransposeToAxes(nbits uint, x []uint64) {
	if nbits == 0 || len(x) == 0 {
		return
	}
	n := len(x)
	t := x[n-1] >> 1
	for i := n - 1; i > 0; i-- {
		x[i] ^= x[i-1]
	}
	x[0] ^= t
	// end is 0 when nbits is 64, so q stops when it overflows.
	end := uint64(2) << (nbits - 1)
	for q := uint64(2); q != end; q <<= 1 {
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
func hilbertInterleave(nbits uint, x []uint64) uint64 {
	var d uint64
	for j := int(nbits) - 1; j >= 0; j-- {
		for _, v := range x {
			d = d<<1 | (v>>uint(j))&1
		}
	}
	return d
}

// hilbertDeinterleave is the inverse of hilbertInterleave. It reads only the
// low nbits*len(x) bits of d.
func hilbertDeinterleave(nbits uint, d uint64, x []uint64) {
	clear(x)
	for j := range nbits {
		for i := len(x) - 1; i >= 0; i-- {
			x[i] |= (d & 1) << j
			d >>= 1
		}
	}
}
