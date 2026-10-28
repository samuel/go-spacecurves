package spacecurves

// Morton2DEncode encodes a 2D x,y pair into a single value along
// a Z-order curve of the specified number of bits per axis (at most 32).
func Morton2DEncode(bits uint, x, y uint64) uint64 {
	checkBits(bits, 2)
	m := lowMask(bits)
	return spread2(x&m) | spread2(y&m)<<1
}

// Morton2DDecode decodes a value along a Z-order curve
// of the specified number of bits per axis into a 2D x,y pair.
func Morton2DDecode(bits uint, d uint64) (uint64, uint64) {
	checkBits(bits, 2)
	d &= lowMask(2 * bits)
	return compact2(d), compact2(d >> 1)
}

// Morton3DEncode encodes a 3D x,y,z point into a single value along
// a Z-order curve of the specified number of bits per axis (at most 21).
func Morton3DEncode(bits uint, x, y, z uint64) uint64 {
	checkBits(bits, 3)
	m := lowMask(bits)
	return spread3(x&m) | spread3(y&m)<<1 | spread3(z&m)<<2
}

// Morton3DDecode decodes a value along a Z-order curve
// of the specified number of bits per axis into a 3D x,y,z point.
func Morton3DDecode(bits uint, d uint64) (uint64, uint64, uint64) {
	checkBits(bits, 3)
	d &= lowMask(3 * bits)
	return compact3(d), compact3(d >> 1), compact3(d >> 2)
}

// spread2 moves bit i of the low 32 bits of x to bit 2i.
func spread2(x uint64) uint64 {
	x &= 0x00000000ffffffff
	x = (x | x<<16) & 0x0000ffff0000ffff
	x = (x | x<<8) & 0x00ff00ff00ff00ff
	x = (x | x<<4) & 0x0f0f0f0f0f0f0f0f
	x = (x | x<<2) & 0x3333333333333333
	x = (x | x<<1) & 0x5555555555555555
	return x
}

// compact2 is the inverse of spread2, ignoring odd bits.
func compact2(x uint64) uint64 {
	x &= 0x5555555555555555
	x = (x | x>>1) & 0x3333333333333333
	x = (x | x>>2) & 0x0f0f0f0f0f0f0f0f
	x = (x | x>>4) & 0x00ff00ff00ff00ff
	x = (x | x>>8) & 0x0000ffff0000ffff
	x = (x | x>>16) & 0x00000000ffffffff
	return x
}

// spread3 moves bit i of the low 21 bits of x to bit 3i.
func spread3(x uint64) uint64 {
	x &= 0x00000000001fffff
	x = (x | x<<32) & 0x001f00000000ffff
	x = (x | x<<16) & 0x001f0000ff0000ff
	x = (x | x<<8) & 0x100f00f00f00f00f
	x = (x | x<<4) & 0x10c30c30c30c30c3
	x = (x | x<<2) & 0x1249249249249249
	return x
}

// compact3 is the inverse of spread3, ignoring bits not at multiples of 3.
func compact3(x uint64) uint64 {
	x &= 0x1249249249249249
	x = (x | x>>2) & 0x10c30c30c30c30c3
	x = (x | x>>4) & 0x100f00f00f00f00f
	x = (x | x>>8) & 0x001f0000ff0000ff
	x = (x | x>>16) & 0x001f00000000ffff
	x = (x | x>>32) & 0x00000000001fffff
	return x
}

// swapMorton2D swaps the x and y bits of a 2D Morton code.
func swapMorton2D(m uint64) uint64 {
	return (m&0x5555555555555555)<<1 | (m>>1)&0x5555555555555555
}

// MortonToHilbert2D converts a 2D Morton code of the specified number of bits
// per axis (at most 32) to the index of the same point on the Hilbert2DEncode
// curve.
func MortonToHilbert2D(bits uint, morton uint64) uint64 {
	checkBits(bits, 2)
	// The state table produces the transpose of Hilbert2DEncode, so swap axes.
	morton = swapMorton2D(morton & lowMask(2*bits))
	var hilbert uint64
	remap := uint64(0xb4)
	block := bits << 1
	for block != 0 {
		block -= 2
		mcode := (morton >> block) & 3
		hcode := (remap >> (mcode << 1)) & 3
		remap ^= 0x82000028 >> (hcode << 3)
		hilbert = (hilbert << 2) + hcode
	}
	return hilbert
}

// HilbertToMorton2D is the inverse of MortonToHilbert2D.
func HilbertToMorton2D(bits uint, hilbert uint64) uint64 {
	checkBits(bits, 2)
	hilbert &= lowMask(2 * bits)
	var morton uint64
	remap := uint64(0xb4)
	block := bits << 1
	for block != 0 {
		block -= 2
		hcode := (hilbert >> block) & 3
		mcode := (remap >> (hcode << 1)) & 3
		remap ^= 0x330000cc >> (hcode << 3)
		morton = (morton << 2) + mcode
	}
	return swapMorton2D(morton)
}
