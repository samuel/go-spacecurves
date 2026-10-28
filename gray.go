package spacecurves

// Gray2DEncode encodes a 2D x,y pair into a single value along a Gray-code
// curve (Faloutsos) of the specified number of bits per axis (at most 32).
// The index is the inverse Gray code of the Morton code, so consecutive
// indices map to points that differ in exactly one bit of one coordinate.
func Gray2DEncode(bits uint, x, y uint64) uint64 {
	return grayDecode(Morton2DEncode(bits, x, y))
}

// Gray2DDecode decodes a value along a Gray-code curve of the specified
// number of bits per axis into a 2D x,y pair.
func Gray2DDecode(bits uint, d uint64) (uint64, uint64) {
	checkBits(bits, 2)
	d &= lowMask(2 * bits)
	return Morton2DDecode(bits, d^(d>>1))
}

// Gray3DEncode encodes a 3D x,y,z point into a single value along a
// Gray-code curve of the specified number of bits per axis (at most 21).
func Gray3DEncode(bits uint, x, y, z uint64) uint64 {
	return grayDecode(Morton3DEncode(bits, x, y, z))
}

// Gray3DDecode decodes a value along a Gray-code curve of the specified
// number of bits per axis into a 3D x,y,z point.
func Gray3DDecode(bits uint, d uint64) (uint64, uint64, uint64) {
	checkBits(bits, 3)
	d &= lowMask(3 * bits)
	return Morton3DDecode(bits, d^(d>>1))
}

// grayDecode returns the value whose Gray code is g.
func grayDecode(g uint64) uint64 {
	for s := uint(1); s < 64; s <<= 1 {
		g ^= g >> s
	}
	return g
}

// GrayEncode encodes an N-dimensional point into a single value along a
// Gray-code curve of the specified number of bits per axis. In 2D and 3D it
// matches Gray2DEncode and Gray3DEncode. The point is not modified.
func GrayEncode(bits uint, point []uint64) uint64 {
	return grayDecode(MortonEncode(bits, point))
}

// GrayDecode decodes a value along a Gray-code curve of the specified number
// of bits per axis into point, whose length sets the number of dimensions.
func GrayDecode(bits uint, d uint64, point []uint64) {
	checkBits(bits, len(point))
	d &= lowMask(uint(len(point)) * bits)
	MortonDecode(bits, d^(d>>1), point)
}
