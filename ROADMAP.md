# ROADMAP

Bring `github.com/effective-security/x` in line with Go 1.27 and the current
standard library. Concrete bugs and security issues are in
[`FINDINGS.md`](FINDINGS.md). This file is the modernization plan.

The module already sets `go 1.27` in `go.mod`. The root README still says
“Go 1.24 or newer”.

---

## 1. `encoding/json/v2` and `encoding/json/jsontext`

Go 1.27 stabilizes `encoding/json/v2` and `encoding/json/jsontext`. v2
rejects invalid UTF-8 and duplicate object names by default, takes
`Options`, and unmarshals faster. v1 `encoding/json` is implemented on
top of v2 and stays compatible.

Current v1 call sites:

| File | Use |
| --- | --- |
| `configloader/load.go` | `Unmarshal` / `Marshal` for `.json` files |
| `print/print.go` | `JSON` pretty-print |
| `values/value.go` | `JSON`, `JSONIndent`, `IndentJSON` |
| `values/mapany.go` | `FromJSON`, `To`, `Scan` / `Value` |
| `values/canonicaljson.go` | scalars + fallback `json.Marshal`, `Decoder.UseNumber` |

### Suggested work

1. **`configloader.Unmarshal` / `Marshal`**
   Switch JSON files to `jsonv2.Unmarshal` / `jsonv2.Marshal` with
   `jsonv2.Deterministic(true)` on write (stable key order) and
   `jsonv2.RejectUnknownMembers(false)` unless callers want strict
   structs. Use `jsonv2.MarshalWrite` to stream to the file instead of
   buffering then `WriteFile`.

2. **`values.FromJSON` / `MapAny.To` / SQL `Scan`**
   New APIs should return errors (v2’s default is the right contract).
   Keep the swallowing helpers as wrappers, documented, or add
   `FromJSONStrict`. v2 duplicate-name rejection fixes silent last-key-wins
   on attacker-controlled JSON.

3. **`values.JSON` / `JSONIndent` / `print.JSON`**
   `jsonv2.Marshal` + `jsontext.WithIndent("\t")` (or `"  "`). Surface
   marshal errors in new functions (`JSONError`, `print.JSONErr`); leave
   the silent helpers for compatibility.

4. **`canonicaljson`**
   This is the highest-value v2 target. Today it re-encodes with v1
   `json.Marshal` then walks `any`. Replace the walk with
   `jsontext.Encoder` / `jsontext.Decoder`:
   - emit objects with keys sorted (the current contract);
   - use `jsonv2.Deterministic(true)` for the struct fallback;
   - keep the depth cap (1000) as a `jsontext` option or a manual counter.

   v2 `Deterministic` alone is **not** RFC 8785 (JCS): number formatting
   and `\u` escapes still differ. If anyone needs JCS, keep a dedicated
   encoder; otherwise document “sorted keys, v2 scalar encoding”.

5. **Do not blindly change the import**
   v2 defaults differ (UTF-8, duplicates, nil slice/map as empty vs
   `null`). Add tests before switching each call site. `GOEXPERIMENT=nojsonv2`
   is a temporary escape hatch only.

---

## 2. Generic methods (new in Go 1.27)

A method may now declare its own type parameters. Interface methods still
cannot. The useful places in this repo are concrete types that today
expose a family of typed helpers at package scope.

### `values.MapAny` — do this first

Replace the typed getters with generic methods on the map:

```go
func (c MapAny) Get[T any](key string) (T, bool)
func (c MapAny) GetOr[T any](key string, def T) T
func (c MapAny) MustGet[T any](key string) T // or panic → keep out of lib
```

Keep `String` / `Int` / `Bool` as wrappers so existing callers compile.
`Get[T]` should use the same coercion as `values.Int` / `String` via a
generic `Convert[T any](v any) (T, bool)`.

`FromStruct` can grow a generic constructor once tags are honored:

```go
func From[T any](v T) (MapAny, error)
```

### `maps.SyncMap`

If the type stays (see §3), add the methods stdlib already has, as
generic methods only where the type param is *not* `K`/`V`:

```go
func (m *SyncMap[K, V]) Transform[R any](f func(K, V) R) []R
```

A better outcome is to delete `SyncMap` and tell callers to use
`sync.Map[K, V]`.

### `print`

```go
func Register[T any](fn func(io.Writer, T))
```

Removes `reflect.TypeOf` at the call site and the `any` cast inside
`CustomFn`. Keep `RegisterType` as a deprecated alias.

### `enum`

Package-level `Parse[E]`, `Convert[E]`, `FlagNames[E]` are already
generic functions; they cannot move onto the enum type unless each
concrete enum defines them. A small helper type works:

```go
type Of[E Values] struct{}
func (Of[E]) Parse(s string) (E, error)
```

Use this when adding `Parse` errors (FINDINGS).

### `format` / `values` coalescing

`Select`, `NumbersCoalesce`, `format.Number` are already generic
functions. No method form is needed. `cmp.Or` (Go 1.22) can replace
`StringsCoalesce` / `NvlString` for comparable zeros.

---

## 3. Standard library that already replaced this code

Adopt these without waiting for new APIs. Deprecate the local names;
do not break callers in one cut.

| Local API | Stdlib (since) | Action |
| --- | --- | --- |
| `maps.Keys` / `Values` | `maps.Keys` / `maps.Values` iterators (1.23) | Return `iter.Seq` or collect with `slices.Collect` |
| `maps.Merge` | `maps.Copy` / `maps.Clone` (1.21) | `Clone` then `Copy` |
| `maps.OrderedKeys` | `slices.Sorted(maps.Keys(m))` (1.23) | One-liner |
| `maps.SyncMap` | `sync.Map[K,V]` (1.24) | Deprecate the wrapper |
| `slices.Contains` | `slices.Contains` (1.21) | Alias |
| `slices.*SlicesEqual` | `slices.Equal` (1.21) | Alias |
| `slices.CloneStrings` | `slices.Clone` (1.21) | Alias (nil-preserving) |
| `slices.Deduplicate` / `UniqueStrings` | `slices.Compact` after sort, or a set | Keep order-preserving helper |
| `slices.Uint64s` | `slices.Sort` (1.21) | Delete `sort.Interface` |
| `slices.Truncate` | `s[:min(n,len(s))]` | Keep if the `uint` API matters |
| `netutil.IsPrivateAddress` | `net.IP.IsPrivate` (1.17) + `IsLoopback` / `IsLinkLocalUnicast` | Thin wrapper |
| `netutil.IsAddrInUse` | `errors.Is` / `errors.As` | Rewrite |
| `guid.MustCreate` | `uuid.New` / `uuid.NewV4` / `uuid.NewV7` (1.27) | See §4 |
| `fileutil/resolve.ExpandPath` | `os.UserHomeDir` + `os.ExpandEnv` | Drop `mitchellh/go-homedir` |
| `values.StringsCoalesce` | `cmp.Or` (1.22) | Alias |
| `format.NowFunc` / `flake.NowFunc` | keep (tests) | — |

`github.com/deckarep/golang-set` v1 appears only in `flake` tests. Use
`map[uint64]struct{}` or `golang-set/v2`.

---

## 4. New Go 1.27 packages this repo should use

### `uuid`

`guid` is a 16-byte hex formatter, not a UUID. Plan:

1. Implement `MustCreate` with `uuid.New()` (v4) so the string is a
   valid RFC 9562 UUID (version + variant bits). **This is a format
   change** for the 13th hex digit; call it out in the changelog.
2. Add `Create() (string, error)` that does not panic.
3. Add `V7()` for time-sortable IDs when callers would otherwise reach
   for `flake` just to sort.
4. Keep uppercase vs lowercase as an explicit option; stdlib `UUID.String`
   is lowercase.

Do not replace `flake` with UUID v7. Flake is a 64-bit sortable ID with
an explicit machine field; v7 is 128-bit and leaks time. They solve
different problems.

### `hash/maphash.Hasher`

`values` already has a hand-rolled `hasher` interface for XXH3.
`maphash.Hasher` / `ComparableHasher` is the stdlib contract for future
generic maps. Implement `Hasher` for `MapAny` (via canonical JSON or
`XXH3HashArgs*`) so callers can put maps in hash tables without
re-encoding ad hoc.

### `strings.CutLast` / `bytes.CutLast`

Use in `format.TextOneLine` and any “split from the right” helpers
instead of `LastIndex` + slice.

### `net/url.URL.Clone` / `Values.Clone`

`urlutil.GetPublicEndpointURL` should start from `r.URL.Clone()` (or a
configured public base) instead of assembling a bare `url.URL`.

### Experimental `simd`

Not a fit. This library is not compute-heavy. Revisit only if
`values.XXH3*` or `flake` ever needs a vectorized path.

---

## 5. Features worth adding (by package)

### `configloader`

- Expand maps correctly (FINDINGS) and support `map[string]any`.
- Bound `file://` to the config directory; cap read size.
- Fail closed on secret errors.
- JSON configs through `Factory.Load` (today only YAML via uber/config).
- `context.Context` on `SecretProvider.GetSecret`.
- Replace `oleiade/reflections` with `reflect` or a small helper; the
  extra dependency is used for one field.

### `enum`

- `Parse` that returns `(E, error)` and trims tokens.
- `TextMarshaler` / `TextUnmarshaler` / `json.Marshaler` so enums work
  with json v2 without per-type boilerplate.
- Generic method form in §2 if you add a helper type.

### `fileutil/reloader`

- `fsnotify` (or `os.Root` + poll) with a context; fix the races in
  FINDINGS.
- Initial callback optional (`WithLoadOnStart`).
- `testing/synctest` + `httptest.NewTestServer` are 1.27-friendly ways
  to drop `time.Sleep` in tests (reloader and ticker still sleep).

### `flake`

- Fix `IsTestRun` and machine-ID `0` (FINDINGS).
- Expose `Decompose` as a struct, not `map[string]uint64`.
- Optional `NextID() (uint64, error)` beside the panicking API.

### `guid`

- Stdlib `uuid` (§4). Optional parse/validate.

### `maps`

- Deprecate `SyncMap`.
- Iterator-returning `Keys` / `Values` / `Range` (`iter.Seq2[K,V]`).
- `Collect` helpers that sit on `iter`.

### `netutil`

- `FindFreePort` returns `net.Listener` (or a `ReservePort` API).
- `WaitForNetwork(ctx)`.
- Export or delete `namedAddress`.

### `print`

- Generic `Register[T]`.
- `JSON` that uses json v2 and returns `error`.
- Drop the process-global registry or scope it to a `Printers` value
  (generic methods on that type).

### `slices`

- Mark typed equals/contains as deprecated aliases.
- Fix `HashStrings` (FINDINGS) or move callers to `values.XXH3HashArgs*`.

### `ticker`

- Reject `every <= 0`.
- `Stop` that waits; or `Wait()`.
- Invoke-on-start option (common request; today the first callback is
  tick 1).

### `urlutil`

- Trusted-proxy option for `X-Forwarded-Proto` / `X-Forwarded-Host`.
- Allow-list schemes `http`/`https` only.
- Build from a configured public base URL (`url.URL.Clone`).

### `values`

- Generic `Get` / `Convert` (§2).
- Honor `json` tags in `FromStruct` (breaking if done in place — add
  `FromStructJSON`).
- json v2 for `FromJSON` / `Scan` with an error-returning API.
- `IsEmpty` / `Add` / `Shrink` that do not treat `false` and `0` as empty.
- Align `Time` parsing with `format.ParseStringTime`.

---

## 6. Suggested order

1. Security fixes in FINDINGS (`configloader` file mode / `file://` /
   secrets, `urlutil` forwarded headers, `slices.HashStrings`).
2. Concurrency fixes (`print.Print` lock, `reloader.Close`).
3. json v2 on `values` + `configloader` + `print` (tests first).
4. `guid` → stdlib `uuid` (changelog).
5. Generic methods on `MapAny` and `print.Register`.
6. Deprecate wrappers that duplicate `slices`, `maps`, and `sync.Map`.
7. Features in §5 as follow-ups.

Keep `Documentation/codemap.md` updated in the same change as each item
(invariants, entry points, panic vs error).
