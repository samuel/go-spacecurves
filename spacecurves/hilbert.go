package spacecurves

// Hilbert2DEncode encodes a 2D x,y pair into a single value along
// a Hilbert space-filling curve of the specified number of bits.
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
// of the specified number of bits into a 2D x,y pair.
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
