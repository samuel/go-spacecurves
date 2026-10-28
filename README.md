Space-filling Curves for Go
===========================

[![Go Reference](https://pkg.go.dev/badge/samuelks.com/spacecurves.svg)](https://pkg.go.dev/samuelks.com/spacecurves)

Package spacecurves implements [space-filling curves](https://en.wikipedia.org/wiki/Space-filling_curve)
to encode/decode multi-dimensional points into single values along a curve. Different
curves maintain locality of points to various degrees.

Implemented curves:

- [Hilbert](https://en.wikipedia.org/wiki/Hilbert_curve)
- [Moore](https://en.wikipedia.org/wiki/Moore_curve) (a closed-loop variant of Hilbert)
- [Morton (Z-Order)](https://en.wikipedia.org/wiki/Z-order_curve)
- Gray-code (Faloutsos), with locality between Morton and Hilbert at about the
  cost of Morton

The generic functions take the number of bits per axis as their first argument.
On 64-bit platforms this is at most 32 for 2D and 21 for 3D. `HilbertEncode` and
`HilbertDecode` support any number of dimensions, as long as dimensions × bits
fits in a `uint`. Fixed-width Morton
variants (`Morton2DEncode5bit`, `Morton2DEncode16bit`, `Morton3DEncode5bit`,
`Morton3DEncode10bit`, and their decoders) are faster for those widths.
`MortonToHilbert2D/3D` and `HilbertToMorton2D/3D` convert directly between
Morton and Hilbert indices. The 2D conversion's curve is the transpose of
`Hilbert2DEncode`'s. 3D Hilbert curves are not unique: `Hilbert3DEncode` and
`MortonToHilbert3D` produce different curves.

Box queries
-----------

For spatial indexes stored in curve order:

- `Morton2DNextInBox` / `Morton2DPrevInBox` (and 3D versions) return the next or
  previous index inside a query box (BIGMIN/LITMAX), so a scan can skip the gaps.
- `Hilbert2DRanges`, `Hilbert3DRanges`, and `HilbertRanges` (N-D) split a query
  box into sorted, inclusive index ranges. An optional cap on the number of
  ranges merges the smallest gaps, trading extra indices for fewer ranges.

License
-------

MIT. See LICENSE file.
