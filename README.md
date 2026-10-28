Space-filling Curves for Go
===========================

[![Go Reference](https://pkg.go.dev/badge/github.com/samuel/go-spacecurves/spacecurves.svg)](https://pkg.go.dev/github.com/samuel/go-spacecurves/spacecurves)

Package spacecurves implements [space-filling curves](https://en.wikipedia.org/wiki/Space-filling_curve)
to encode/decode multi-dimensional points into single values along a curve. Different
curves maintain locality of points to various degrees.

Implemented curves:

- [Hilbert](https://en.wikipedia.org/wiki/Hilbert_curve)
- [Morton (Z-Order)](https://en.wikipedia.org/wiki/Z-order_curve)

The generic functions take the number of bits per axis as their first argument.
On 64-bit platforms this is at most 32 for 2D and 21 for 3D. Fixed-width Morton
variants (`Morton2DEncode5bit`, `Morton2DEncode16bit`, `Morton3DEncode5bit`,
`Morton3DEncode10bit`, and their decoders) are faster for those widths.
`MortonToHilbert2D/3D` and `HilbertToMorton2D/3D` convert directly between
Morton and Hilbert indices. The 2D conversion's curve is the transpose of
`Hilbert2DEncode`'s.

License
-------

3-clause BSD. See LICENSE file.
