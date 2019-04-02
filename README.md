Space-filling Curves for Go
===========================

[![Documentation](https://godoc.org/github.com/samuel/go-spacecurves/spacecurves?status.svg)](https://godoc.org/github.com/samuel/go-spacecurves/spacecurves)

Package spacecurves implements [space-filling curves](https://en.wikipedia.org/wiki/Space-filling_curve)
to encode/decode multi-dimensional points into single values along a curve. Different
curves maintain locality of points to various degrees.

Implemented curves:

- [Hilbert](https://en.wikipedia.org/wiki/Hilbert_curve)
- [Morton (Z-Order)](https://en.wikipedia.org/wiki/Z-order_curve)

License
-------

3-clause BSD. See LICENSE file.
