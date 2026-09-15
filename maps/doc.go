// Package maps provides generic map helpers and a typed wrapper around sync.Map.
//
// Several entry points are deprecated in favor of the standard library:
// Keys/Values/OrderedKeys (maps + slices, Go 1.23) and Merge
// (maps.Clone/Copy). SyncMap remains useful because sync.Map is untyped.
package maps
