package spacecurves

func Morton2DEncode(bits, x, y uint) uint {
	var answer uint
	s := uint(1)
	for i := uint(0); i < bits; i++ {
		answer |= (x & s) << i
		answer |= (y & s) << (i + 1)
		s <<= 1
	}
	return answer
}

func Morton3DEncode(bits, x, y, z uint) uint {
	var answer uint
	s := uint(1)
	for i := uint(0); i < bits; i++ {
		answer |= (x & s) << (2 * i)
		answer |= (y & s) << (2*i + 1)
		answer |= (z & s) << (2*i + 2)
		s <<= 1
	}
	return answer
}

func Morton2DDecode(n, d uint) (uint, uint) {
	var x, y uint
	s := uint(1)
	for i := uint(0); i < n; i++ {
		x |= (d >> i) & s
		y |= (d >> (i + 1)) & s
		s <<= 1
	}
	return x, y
}

func Morton3DDecode(n, d uint) (uint, uint, uint) {
	var x, y, z uint
	s := uint(1)
	for i := uint(0); i < n; i++ {
		x |= (d >> (i * 2)) & s
		y |= (d >> (i*2 + 1)) & s
		z |= (d >> (i*2 + 2)) & s
		s <<= 1
	}
	return x, y, z
}

func MortonToHilbert2D(morton, bits uint) uint {
	hilbert := uint(0)
	remap := uint(0xb4)
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

func HilbertToMorton2D(hilbert, bits uint) uint {
	morton := uint(0)
	remap := uint(0xb4)
	block := bits << 1
	for block != 0 {
		block -= 2
		hcode := (hilbert >> block) & 3
		mcode := (remap >> (hcode << 1)) & 3
		remap ^= 0x330000cc >> (hcode << 3)
		morton = (morton << 2) + mcode
	}
	return morton
}

func MortonToHilbert3D(morton, bits uint) uint {
	hilbert := morton
	if bits > 1 {
		block := (bits * 3) - 3
		hcode := (hilbert >> block) & 7
		shift := uint(0)
		signs := uint(0)
		for block != 0 {
			block -= 3
			hcode <<= 2
			mcode := (uint(0x20212021) >> hcode) & 3
			shift = (0x48 >> (7 - shift - mcode)) & 3
			signs = (signs | (signs << 3)) >> mcode
			signs = (signs ^ (0x53560300 >> hcode)) & 7
			mcode = (hilbert >> block) & 7
			hcode = mcode
			hcode = ((hcode | (hcode << 3)) >> shift) & 7
			hcode ^= signs
			hilbert ^= (mcode ^ hcode) << block
		}
	}
	hilbert ^= (hilbert >> 1) & 0x92492492
	hilbert ^= (hilbert & 0x92492492) >> 1
	return hilbert
}

func HilbertToMorton3D(hilbert, bits uint) uint {
	morton := hilbert
	morton ^= (morton & 0x92492492) >> 1
	morton ^= (morton >> 1) & 0x92492492
	if bits > 1 {
		block := ((bits * 3) - 3)
		hcode := ((morton >> block) & 7)
		shift := uint(0)
		signs := uint(0)
		for block != 0 {
			block -= 3
			hcode <<= 2
			mcode := (uint(0x20212021) >> hcode) & 3
			shift = (0x48 >> (4 - shift + mcode)) & 3
			signs = (signs | (signs << 3)) >> mcode
			signs = (signs ^ (0x53560300 >> hcode)) & 7
			hcode = (morton >> block) & 7
			mcode = hcode
			mcode ^= signs
			mcode = ((mcode | (mcode << 3)) >> shift) & 7
			morton ^= (hcode ^ mcode) << block
		}
	}
	return morton
}

func Morton2DEncode5bit(x, y uint) uint {
	x &= 0x0000001f
	y &= 0x0000001f
	x *= 0x01041041
	y *= 0x01041041
	x &= 0x10204081
	y &= 0x10204081
	x *= 0x00108421
	y *= 0x00108421
	x &= 0x15500000
	y &= 0x15500000
	return (x >> 20) | (y >> 19)
}

func Morton2DDecode5bit(morton uint) (uint, uint) {
	value1 := morton
	value2 := value1 >> 1
	value1 &= 0x00000155
	value2 &= 0x00000155
	value1 |= value1 >> 1
	value2 |= value2 >> 1
	value1 &= 0x00000133
	value2 &= 0x00000133
	value1 |= value1 >> 2
	value2 |= value2 >> 2
	value1 &= 0x0000010f
	value2 &= 0x0000010f
	value1 |= value1 >> 4
	value2 |= value2 >> 4
	value1 &= 0x0000001f
	value2 &= 0x0000001f
	return value1, value2
}

func Morton2DEncode16bit(x, y uint) uint {
	x &= 0x0000ffff
	y &= 0x0000ffff
	x |= x << 8
	y |= y << 8
	x &= 0x00ff00ff
	y &= 0x00ff00ff
	x |= x << 4
	y |= y << 4
	x &= 0x0f0f0f0f
	y &= 0x0f0f0f0f
	x |= x << 2
	y |= y << 2
	x &= 0x33333333
	y &= 0x33333333
	x |= x << 1
	y |= y << 1
	x &= 0x55555555
	y &= 0x55555555
	return x | (y << 1)
}

func Morton2DDecode16bit(morton uint) (uint, uint) {
	value1 := morton
	value2 := value1 >> 1
	value1 &= 0x55555555
	value2 &= 0x55555555
	value1 |= value1 >> 1
	value2 |= value2 >> 1
	value1 &= 0x33333333
	value2 &= 0x33333333
	value1 |= value1 >> 2
	value2 |= value2 >> 2
	value1 &= 0x0f0f0f0f
	value2 &= 0x0f0f0f0f
	value1 |= value1 >> 4
	value2 |= value2 >> 4
	value1 &= 0x00ff00ff
	value2 &= 0x00ff00ff
	value1 |= value1 >> 8
	value2 |= value2 >> 8
	value1 &= 0x0000ffff
	value2 &= 0x0000ffff
	return value1, value2
}

func Morton3DEncode5bit(x, y, z uint) uint {
	x &= 0x0000001f
	y &= 0x0000001f
	z &= 0x0000001f
	x *= 0x01041041
	y *= 0x01041041
	z *= 0x01041041
	x &= 0x10204081
	y &= 0x10204081
	z &= 0x10204081
	x *= 0x00011111
	y *= 0x00011111
	z *= 0x00011111
	x &= 0x12490000
	y &= 0x12490000
	z &= 0x12490000
	return (x >> 16) | (y >> 15) | (z >> 14)
}

func Morton3DDecode5bit(morton uint) (uint, uint, uint) {
	value1 := morton
	value2 := value1 >> 1
	value3 := value1 >> 2
	value1 &= 0x00001249
	value2 &= 0x00001249
	value3 &= 0x00001249
	value1 |= value1 >> 2
	value2 |= value2 >> 2
	value3 |= value3 >> 2
	value1 &= 0x000010c3
	value2 &= 0x000010c3
	value3 &= 0x000010c3
	value1 |= value1 >> 4
	value2 |= value2 >> 4
	value3 |= value3 >> 4
	value1 &= 0x0000100f
	value2 &= 0x0000100f
	value3 &= 0x0000100f
	value1 |= value1 >> 8
	value2 |= value2 >> 8
	value3 |= value3 >> 8
	value1 &= 0x0000001f
	value2 &= 0x0000001f
	value3 &= 0x0000001f
	return value1, value2, value3
}

func Morton3DEncode10bit(x, y, z uint) uint {
	x &= 0x000003ff
	y &= 0x000003ff
	z &= 0x000003ff
	x |= x << 16
	y |= y << 16
	z |= z << 16
	x &= 0x030000ff
	y &= 0x030000ff
	z &= 0x030000ff
	x |= x << 8
	y |= y << 8
	z |= z << 8
	x &= 0x0300f00f
	y &= 0x0300f00f
	z &= 0x0300f00f
	x |= x << 4
	y |= y << 4
	z |= z << 4
	x &= 0x030c30c3
	y &= 0x030c30c3
	z &= 0x030c30c3
	x |= x << 2
	y |= y << 2
	z |= z << 2
	x &= 0x09249249
	y &= 0x09249249
	z &= 0x09249249
	return x | (y << 1) | (z << 2)
}

func Morton3DDecode10bit(morton uint) (uint, uint, uint) {
	value1 := morton
	value2 := value1 >> 1
	value3 := value1 >> 2
	value1 &= 0x09249249
	value2 &= 0x09249249
	value3 &= 0x09249249
	value1 |= value1 >> 2
	value2 |= value2 >> 2
	value3 |= value3 >> 2
	value1 &= 0x030c30c3
	value2 &= 0x030c30c3
	value3 &= 0x030c30c3
	value1 |= value1 >> 4
	value2 |= value2 >> 4
	value3 |= value3 >> 4
	value1 &= 0x0300f00f
	value2 &= 0x0300f00f
	value3 &= 0x0300f00f
	value1 |= value1 >> 8
	value2 |= value2 >> 8
	value3 |= value3 >> 8
	value1 &= 0x030000ff
	value2 &= 0x030000ff
	value3 &= 0x030000ff
	value1 |= value1 >> 16
	value2 |= value2 >> 16
	value3 |= value3 >> 16
	value1 &= 0x000003ff
	value2 &= 0x000003ff
	value3 &= 0x000003ff
	return value1, value2, value3
}
