// Package spacecurves implements space-filling curves to encode/decode
// multi-dimensional points into single values along a curve. Different
// curves maintain locality of points to various degrees.
//
// Coordinates and indices are uint64. Each function takes bits, the number
// of bits per axis, and the index holds dims*bits bits, so bits is at most
// 32 in 2D and 21 in 3D. Bits of coordinates and indices above that width are
// ignored. Functions panic if dims*bits exceeds 64. Box query bounds are
// inclusive and clamped to the grid.
package spacecurves
