// Package spacecurves implements space-filling curves to encode/decode
// multi-dimensional points into single values along a curve. Different
// curves maintain locality of points to various degrees.
//
// # Choosing a curve
//
//   - Hilbert (Hilbert2DEncode, Hilbert3DEncode, HilbertEncode for any number
//     of dimensions) has the best locality: consecutive indices are always
//     adjacent cells, and a box maps to few index ranges. It is the usual
//     choice for spatial indexes.
//   - Moore (Moore2DEncode) is a Hilbert variant that forms a closed loop, for
//     cyclic traversals.
//   - Morton, or Z-order (Morton2DEncode, Morton3DEncode, MortonEncode), is
//     the fastest to compute, since it only interleaves bits, but has jumps
//     between distant cells.
//   - Gray-code (Gray2DEncode, Gray3DEncode, GrayEncode) costs about the same
//     as Morton and has locality between Morton and Hilbert.
//
// Curve outputs are stable across releases, so stored indices remain valid.
//
// # Bits and input values
//
// Coordinates and indices are uint64. Each function takes bits, the number
// of bits per axis, so the grid is 2^bits cells along each axis and the index
// holds dims*bits bits. bits is at most 32 in 2D and 21 in 3D, and for the
// N-dimensional functions dims*bits is at most 64. Bits of coordinates and
// indices above that width are ignored. Functions panic if dims*bits exceeds
// 64.
//
// # Box queries
//
// Morton2DNextInBox and Morton2DPrevInBox (and the 3D versions) find the next
// or previous index inside a query box, so a scan of points stored in Morton
// order can skip the gaps between runs. Hilbert2DRanges, Hilbert3DRanges,
// HilbertRanges and the Morton equivalents split a box into sorted index
// ranges to scan; Hilbert order needs fewer ranges than Morton order. An
// optional cap on the number of ranges merges the closest ranges, trading
// extra indices for fewer ranges. Box bounds are inclusive and clamped to the
// grid.
package spacecurves
