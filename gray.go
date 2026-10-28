package spacecurves

// Gray2DEncode encodes a 2D x,y pair into a single value along a Gray-code
// curve (Faloutsos) of the specified number of bits per axis (at most 32 on
// 64-bit platforms). The index is the inverse Gray code of the Morton code,
// so consecutive indices map to points that differ in exactly one bit of one
// coordinate.
func Gray2DEncode(bits, x, y uint) uint {
	return grayDecode(Morton2DEncode(bits, x, y))
}

// Gray2DDecode decodes a value along a Gray-code curve of the specified
// number of bits per axis into a 2D x,y pair.
func Gray2DDecode(bits, d uint) (uint, uint) {
	d &= lowMask(2 * bits)
	return Morton2DDecode(bits, d^(d>>1))
}

// Gray3DEncode encodes a 3D x,y,z point into a single value along a
// Gray-code curve of the specified number of bits per axis (at most 21 on
// 64-bit platforms).
func Gray3DEncode(bits, x, y, z uint) uint {
	return grayDecode(Morton3DEncode(bits, x, y, z))
}

// Gray3DDecode decodes a value along a Gray-code curve of the specified
// number of bits per axis into a 3D x,y,z point.
func Gray3DDecode(bits, d uint) (uint, uint, uint) {
	d &= lowMask(3 * bits)
	return Morton3DDecode(bits, d^(d>>1))
}

// grayDecode returns the value whose Gray code is g.
func grayDecode(g uint) uint {
	for s := uint(1); s < uintBits; s <<= 1 {
		g ^= g >> s
	}
	return g
}

// lowMask returns a mask of the low n bits, or all bits when n is at least
// the width of uint.
func lowMask(n uint) uint {
	return uint(1)<<n - 1
}
