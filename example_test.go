package spacecurves_test

import (
	"fmt"

	"samuelks.com/spacecurves"
)

func ExampleHilbert2DEncode() {
	// A 4x4 grid has 2 bits per axis. Neighboring indices are neighboring cells.
	const bits = 2
	for _, p := range [][2]uint64{{0, 0}, {1, 0}, {1, 1}, {3, 0}} {
		d := spacecurves.Hilbert2DEncode(bits, p[0], p[1])
		x, y := spacecurves.Hilbert2DDecode(bits, d)
		fmt.Printf("(%d, %d) -> %2d -> (%d, %d)\n", p[0], p[1], d, x, y)
	}
	// Output:
	// (0, 0) ->  0 -> (0, 0)
	// (1, 0) ->  1 -> (1, 0)
	// (1, 1) ->  2 -> (1, 1)
	// (3, 0) -> 15 -> (3, 0)
}

func ExampleHilbertEncode() {
	// Points of any number of dimensions, here 4 axes of 16 bits.
	p := []uint64{1, 2, 3, 4}
	d := spacecurves.HilbertEncode(16, p)
	q := make([]uint64, 4)
	spacecurves.HilbertDecode(16, d, q)
	fmt.Println(d, q)
	// Output:
	// 3940 [1 2 3 4]
}

func ExampleMorton2DNextInBox() {
	// Scan points stored in Morton order, skipping the gaps between runs of
	// indices inside the query box x 2..3, y 1..2.
	const bits = 3
	d := uint64(0)
	for {
		next, ok := spacecurves.Morton2DNextInBox(bits, d, 2, 1, 3, 2)
		if !ok {
			break
		}
		x, y := spacecurves.Morton2DDecode(bits, next)
		fmt.Printf("index %2d is (%d, %d)\n", next, x, y)
		d = next + 1
	}
	// Output:
	// index  6 is (2, 1)
	// index  7 is (3, 1)
	// index 12 is (2, 2)
	// index 13 is (3, 2)
}

func ExampleHilbert2DRanges() {
	// The index ranges to scan for the box x 1..2, y 1..2 on an 8x8 grid.
	fmt.Println(spacecurves.Hilbert2DRanges(3, 1, 1, 2, 2, 0))
	// With a cap of 2, the ranges with the smallest gap are merged, adding
	// indices outside the box.
	fmt.Println(spacecurves.Hilbert2DRanges(3, 1, 1, 2, 2, 2))
	// Output:
	// [{2 2} {7 8} {13 13}]
	// [{2 8} {13 13}]
}

func ExampleMoore2DEncode() {
	// The Moore curve is a loop: the last index is next to the first.
	const bits = 2
	x0, y0 := spacecurves.Moore2DDecode(bits, 0)
	x1, y1 := spacecurves.Moore2DDecode(bits, 15)
	fmt.Printf("first (%d, %d), last (%d, %d)\n", x0, y0, x1, y1)
	// Output:
	// first (1, 0), last (2, 0)
}
