# Changelog

## 1.0.0 (unreleased)

First stable release. Curve outputs are now fixed, so stored indices stay valid
across releases.

### Added

- N-dimensional Hilbert (`HilbertEncode`/`HilbertDecode`) and 3D Hilbert
  (`Hilbert3DEncode`/`Hilbert3DDecode`) curves.
- Moore curve (`Moore2DEncode`/`Moore2DDecode`), a closed-loop variant of Hilbert.
- Gray-code curve (Faloutsos) in 2D, 3D, and N dimensions.
- N-dimensional Morton (`MortonEncode`/`MortonDecode`).
- Box queries:
  - `Morton2DNextInBox`/`Morton2DPrevInBox` and 3D versions (BIGMIN/LITMAX).
  - `Hilbert2DRanges`, `Hilbert3DRanges`, `HilbertRanges`, `Morton2DRanges`,
    `Morton3DRanges` and `MortonRanges` split a box into index ranges, with an
    optional cap on the number of ranges.

### Changed

- Coordinates and indices are `uint64` on every platform. `bits` is always the
  first argument.
- `MortonToHilbert2D`/`HilbertToMorton2D` produce `Hilbert2DEncode` indices.
  Previously they produced its transpose.
- Generic Morton encoding uses bit tricks at full width and is fast for any
  number of bits.
- Bits of coordinates and indices above the curve's width are ignored, and
  functions panic if dimensions × bits exceeds 64.

### Removed

- `MortonToHilbert3D`/`HilbertToMorton3D`, which produced a different 3D curve
  from `Hilbert3DEncode`. Use `Hilbert3DEncode` with decoded coordinates.
- Fixed-width Morton functions (`Morton2DEncode5bit`, `Morton2DEncode16bit`,
  `Morton3DEncode5bit`, `Morton3DEncode10bit` and their decoders). The generic
  functions are as fast.

### Fixed

- `MortonToHilbert3D` and `HilbertToMorton3D` were wrong above 10 bits per axis
  on 64-bit platforms. (These functions are removed in this release.)
