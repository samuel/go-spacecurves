package spacecurves

// Moore2DEncode encodes a 2D x,y pair into a single value along a Moore
// curve of the specified number of bits per axis (at most 32 on 64-bit
// platforms). The Moore curve is a closed loop: the first point (2^(bits-1)-1, 0)
// and the last point (2^(bits-1), 0) are adjacent.
func Moore2DEncode(bits, x, y uint) uint {
	if bits == 0 {
		return 0
	}
	mask := uint(1)<<bits - 1
	x &= mask
	y &= mask
	// Each quadrant holds a rotated Hilbert curve of order bits-1.
	h := uint(1) << (bits - 1)
	var q, u, v uint
	switch {
	case x < h && y < h:
		q, u, v = 0, y, h-1-x
	case x < h:
		q, u, v = 1, y-h, h-1-x
	case y >= h:
		q, u, v = 2, 2*h-1-y, x-h
	default:
		q, u, v = 3, h-1-y, x-h
	}
	return q*h*h + Hilbert2DEncode(bits-1, u, v)
}

// Moore2DDecode decodes a value along a Moore curve of the specified number
// of bits per axis into a 2D x,y pair.
func Moore2DDecode(bits, d uint) (uint, uint) {
	if bits == 0 {
		return 0, 0
	}
	h := uint(1) << (bits - 1)
	u, v := Hilbert2DDecode(bits-1, d)
	switch (d >> (2 * (bits - 1))) & 3 {
	case 0:
		return h - 1 - v, u
	case 1:
		return h - 1 - v, h + u
	case 2:
		return h + v, 2*h - 1 - u
	default:
		return h + v, h - 1 - u
	}
}
