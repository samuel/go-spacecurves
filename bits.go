package spacecurves

import "fmt"

// lowMask returns a mask of the low n bits, or all bits when n is at least 64.
func lowMask(n uint) uint64 {
	return uint64(1)<<n - 1
}

// checkBits panics if dims axes of bits each do not fit in a 64-bit index.
func checkBits(bits uint, dims int) {
	if dims > 0 && bits > 64/uint(dims) {
		panicBits(bits, dims)
	}
}

func panicBits(bits uint, dims int) {
	panic(fmt.Sprintf("spacecurves: %d dimensions of %d bits exceeds 64 bits", dims, bits))
}
