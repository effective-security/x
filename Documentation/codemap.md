# Code Map

Navigation index for this module. Open the owning file instead of grepping
the tree. If a concept is not listed here, search, then add it in the same
change.

This is a **library**, not an application: there is no `cmd/`, no `main`,
and no generated mocks. Each top-level directory is an independent Go
package under `github.com/effective-security/x`.

## How to use this map

1. Look up the concept in [Concept index](#concept-index).
2. Open the listed file; exported entry points are in [Packages](#packages).
3. Grep only when the concept is missing, then update this file.

Module path: `github.com/effective-security/x` (`go.mod`, currently Go 1.27).
Release version lives in `.VERSION`; CI tags from that file.

## Concept index

| I need…                                              | Open                                                                                                     |
| ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| YAML/JSON config load with env interpolation         | [`configloader/configloader.go`](../configloader/configloader.go)                                        |
| `${VAR}`, `env://`, `file://`, `secret://` expansion | [`configloader/expand.go`](../configloader/expand.go), [`configloader/load.go`](../configloader/load.go) |
| Host-specific config override (`.hostmap`)           | [`configloader/configloader.go`](../configloader/configloader.go) `Factory.load`                         |
| Marshal/unmarshal JSON or YAML file                  | [`configloader/load.go`](../configloader/load.go) `Unmarshal`, `Marshal`                                 |
| Kong CLI `--version` flag                            | [`ctl/ctl.go`](../ctl/ctl.go) `VersionFlag`                                                              |
| Kong `*bool` flag mapper                             | [`ctl/bool.go`](../ctl/bool.go) `BoolPtrMapper`                                                          |
| Enum / bit-flag parse, names, masks                  | [`enum/enum.go`](../enum/enum.go)                                                                        |
| Folder/file exists, list names, mkdir                | [`fileutil/folders.go`](../fileutil/folders.go)                                                          |
| Watch a file and reload on mtime                     | [`fileutil/reloader/reloader.go`](../fileutil/reloader/reloader.go)                                      |
| Resolve relative path, `~`, `$VAR`                   | [`fileutil/resolve/resolve.go`](../fileutil/resolve/resolve.go)                                          |
| Distributed unique uint64 ID                         | [`flake/flake.go`](../flake/flake.go)                                                                    |
| Human-readable bool/number/string truncation         | [`format/format.go`](../format/format.go)                                                                |
| CamelCase / snake_case display names                 | [`format/format.go`](../format/format.go) `DisplayName`, `Split`                                         |
| Parse/print times, "ago", elapsed                    | [`format/time.go`](../format/time.go)                                                                    |
| Random UUID string                                   | [`guid/guid.go`](../guid/guid.go) `MustCreate`                                                           |
| Generic map Keys/Filter/Merge/Range                  | [`maps/funcs.go`](../maps/funcs.go)                                                                      |
| Typed `sync.Map`                                     | [`maps/generic_map.go`](../maps/generic_map.go) `SyncMap`                                                |
| Free TCP port                                        | [`netutil/freeport.go`](../netutil/freeport.go)                                                          |
| Local / private IP, wait for network                 | [`netutil/localip.go`](../netutil/localip.go)                                                            |
| Hostname, local IP, node name                        | [`netutil/nodeinfo.go`](../netutil/nodeinfo.go)                                                          |
| Parse/join comma-separated URLs                      | [`netutil/urls.go`](../netutil/urls.go)                                                                  |
| `EADDRINUSE` check                                   | [`netutil/net.go`](../netutil/net.go) `IsAddrInUse`                                                      |
| Print JSON/YAML/table/text                           | [`print/print.go`](../print/print.go)                                                                    |
| Slice equal/contains/unique/hash                     | [`slices/slices.go`](../slices/slices.go)                                                                |
| Sortable `[]uint64`                                  | [`slices/uint64s.go`](../slices/uint64s.go)                                                              |
| Context-aware ticker + callback                      | [`ticker/ticker.go`](../ticker/ticker.go)                                                                |
| Query string + public URL from request               | [`urlutil/urlutil.go`](../urlutil/urlutil.go)                                                            |
| `any` → string/int/bool/time/JSON                    | [`values/value.go`](../values/value.go)                                                                  |
| First non-empty string/number                        | [`values/coalesce.go`](../values/coalesce.go)                                                            |
| Nested `map[string]any` helpers                      | [`values/mapany.go`](../values/mapany.go)                                                                |
| Canonical JSON (sorted keys)                         | [`values/canonicaljson.go`](../values/canonicaljson.go)                                                  |
| Fast non-crypto XXH3 hashes                          | [`values/hash_fast.go`](../values/hash_fast.go)                                                          |
| Drop empty nested values                             | [`values/value.go`](../values/value.go) `Shrink`                                                         |

## Internal dependencies

Packages are independent except:

```
configloader → fileutil/resolve, netutil
values       → enum
print        → slices
```

Do not add new cross-package imports without updating this section.

## Packages

### `configloader`

Load YAML (via `go.uber.org/config`) into a struct, then expand `${VAR}` and
`env://` / `file://` / `secret://` values. JSON is supported only by the
standalone `Unmarshal` / `Marshal` helpers, not by `Factory.Load`.

| File                                                 | Owns                                                                                                                                                                |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`configloader.go`](../configloader/configloader.go) | `Factory`, `NewFactory`, `Load`, `LoadForHostName`, `WithOverride`, `WithEnvironment`, `WithSecretProvider`, `ResolveConfigFile`, `GetAbsFilename`, hostmap overlay |
| [`expand.go`](../configloader/expand.go)             | `Expander`, `ExpandAll`, `${}` substitution                                                                                                                         |
| [`load.go`](../configloader/load.go)                 | `SecretProvider`, `ResolveValue*`, `Unmarshal`, `UnmarshalAndExpand`, `Marshal`, source prefixes                                                                    |

**Entry points:** `NewFactory(nodeInfo, searchDirs, envPrefix)`, then
`Load(configFile, &cfg)`. `nodeInfo` may be nil (uses `netutil.NewNodeInfo`).

**Invariants:**

- Interpolation keys: `HOSTNAME`, `NODENAME`, `LOCALIP`, `USER`,
  `NORMALIZED_USER`, `ENVIRONMENT`, `ENVIRONMENT_UPPERCASE`, plus any OS env
  whose name starts with `envPrefix`. Prefixed copies of the built-in keys
  are also injected (`MYAPP_HOSTNAME`, …).
- `${VAR}` uses `Expander.Variables` first, then `os.Getenv`. A missing
  name becomes empty (`os.Expand`). Leftover `${` after that pass is an
  error. `${secret://…}` fails if the provider is missing or `GetSecret`
  errors (it does not become `""`).
- Maps (including aliases and `map[string]any`) are expanded by copying
  values out and writing them back; string fields in nested structs too.
- `env://NAME` fails if the variable is unset or empty.
- `secret://` needs a `SecretProvider`. `ExpandAll` uses the process-global
  `SecretProviderInstance` unless you set `Expander.SecretProvider`.
- `Marshal` writes files with mode `0600`.
- `USER` / `NORMALIZED_USER` become `"unknown"` if `user.Current` fails
  (no panic). `Environment` is read only when the field is a `string`.
- Host overlay: sibling file `<config>.hostmap` with `override: { HOST: file }`.
  Host is `hostnameOverride`, else `{envPrefix}HOSTNAME`, else `os.Hostname()`.
- `Factory.Load` requires a non-empty config path (`ResolveConfigFile` panics
  on `""`). Relative paths are resolved against `searchDirs`.
- `Environment` is written onto the config struct via reflection when
  `WithEnvironment` is set; missing field is ignored.
- `{envPrefix}CONFIG_DIR` is set to the config base directory when empty.

See [`configloader/README.md`](../configloader/README.md) and
[`configloader/testdata/`](../configloader/testdata/).

### `ctl`

Kong helpers for CLIs.

| File                        | Owns                                                                                                                  |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| [`ctl.go`](../ctl/ctl.go)   | `VersionFlag` — prints version from `kong.Vars["version"]` (else the flag's string value) and `app.Exit(0)`           |
| [`bool.go`](../ctl/bool.go) | `BoolPtrMapper` for `*bool`; accepts `true/1/yes` and `false/0/no` (case-insensitive); flag without value sets `true` |

### `enum`

Generic `int32` enums and bit flags. Types must implement `ValuesMap()` /
`NamesMap()` / `DisplayNamesMap()` as required by the constraint.

**Entry points:** `Parse`, `Convert`, `SupportedNames`, `FlagNames`, `Flags`,
`FlagsInt`, `BitMask`, `SliceNames*`, `SliceDisplayNames*`.

**Invariants:**

- `Parse` splits on `|` or `,` and bitwise-ORs named values. Unknown names
  contribute `0`.
- Flag helpers walk all 32 bits of the `int32` (including bit 31) and skip
  unnamed bits.
- `ProtoEnum` is the protobuf-shaped subset (`String`, `Descriptor`, `Number`).

### `fileutil`

| File                                                       | Owns                                                                                                                                                                                                  |
| ---------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`folders.go`](../fileutil/folders.go)                     | `FolderExists`, `FileExists` (error if missing/wrong kind), `SubfolderNames`, `FileNames` (names only, not paths), `EnsureFolderExists`, `EnsureFolderExistsForFile`                                  |
| [`reloader/reloader.go`](../fileutil/reloader/reloader.go) | Poll mtime; `NewReloader`, `Reload`, `Close`, `LoadedAt`, `LoadedCount`. Callback runs on a new goroutine. `Close` takes a write lock and closes the stop channel; a second `Close` returns an error. |
| [`resolve/resolve.go`](../fileutil/resolve/resolve.go)     | `Directory` (optional create, mode `0744`), `File`, `ExpandPath` (`$VAR` then `~`)                                                                                                                    |

Empty path to `Directory`/`File`/`ExpandPath` is returned unchanged.

### `flake`

Sonyflake-inspired unique `uint64`. Fork notes and bit layout:
[`flake/README.md`](../flake/README.md).

**Entry points:** `NewIDGenerator(Settings)`, `DefaultIDGenerator` (init),
`NextID`, `Decompose`, `IDTime`, `FirstID`, `LastID`, `IsTestRun`.

**Layout (LSB → MSB):** 16-bit machine ID, 6-bit sequence, 41-bit time in
1 ms units from `StartTime` (default `2021-01-01 UTC`). ~69 years, 64 IDs/ms
per machine.

**Invariants:**

- `NewIDGenerator` and `NextID` **panic** on bad start time, machine-ID
  failure, failed uniqueness check, or time overflow. They do not return
  those as errors.
- Default machine ID: when `testing.Testing()` is true, FNV-32 of the
  program path; otherwise lower 16 bits of the first non-loopback interface
  address. Missing interfaces → ID `0`, not an error.
- `NowFunc` is overridable for tests.
- `FirstID`/`LastID`/`IDTime` only inspect `*Flake`; other `IDGenerator`
  implementations return `0` / epoch-based time.

### `format`

Display helpers. No errors; empty/zero inputs become `""`, `"never"`, or
similar display strings.

| File                               | Owns                                                                                                                                                   |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`format.go`](../format/format.go) | `YesNo`, `Enabled`, `Number`, `Float`, `StringMax`, `Strings`, `StringsMax`, `StringsAndMore`, `DisplayName`, `Split`, `TextWithIndent`, `TextOneLine` |
| [`time.go`](../format/time.go)     | `ParseStringTime`, `ParseTime`, `Time`, `LocalTime`, `TimeAgo`, `TimesElapsed`                                                                         |

**Invariants:** `Time` follows package-level `DefaultTimePrintFormat`
(`UTC` / `Local` / `Ago`). Parsed times truncate to
`DefaultTimeTruncate` (millisecond). `NowFunc` is overridable.
`ParseStringTime` tries a layout list (not string length).
`StringMax` truncates on runes. `Split` keeps acronyms (`ID`, `HTML`)
and plural acronyms (`IDs`).

### `guid`

[`guid.go`](../guid/guid.go): `MustCreate()` is **Deprecated**; it returns
`uuid.NewV7().String()` (lowercase). Prefer `uuid.NewV7` / `uuid.NewV4`
from the standard library.

### `maps`

Generic map helpers.

| File                                       | Owns                                                                                                                                                                                                                 |
| ------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`funcs.go`](../maps/funcs.go)             | `Keys`, `OrderedKeys`, `Values`, `Each`, `Map`, `Filter`, `Reduce`, `Find`, `Any`, `All`, `Merge` (later maps win), `Invert`, `GroupBy`, `Partition`, `Count`, `Min`, `Max`, `Take`, `Drop`, `Range`, `OrderedRange` |
| [`generic_map.go`](../maps/generic_map.go) | `SyncMap[K,V]` over `sync.Map`                                                                                                                                                                                       |

`Take`/`Drop`/`Range` follow Go map iteration order. `OrderedKeys` /
`OrderedRange` sort keys. `Invert` groups duplicate values as `[]K`.
`Keys`, `Values`, `OrderedKeys`, and `Merge` are Deprecated in favor of
the standard library (`maps`, `slices.Collect`/`AppendSeq`/`Sorted`).
`SyncMap` remains a typed wrapper because the standard `sync.Map` uses
`any` keys and values.

### `netutil`

| File                                    | Owns                                                                                                                     |
| --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| [`freeport.go`](../netutil/freeport.go) | `FindFreePort(host, maxAttempts)` — listen on `:0`, close, return port. Empty host → `localhost`. `maxAttempts < 1` → 1. |
| [`localip.go`](../netutil/localip.go)   | `GetLocalIP` (first non-loopback IPv4), `IsPrivateAddress`, `WaitForNetwork`                                             |
| [`nodeinfo.go`](../netutil/nodeinfo.go) | `NodeInfo`, `NewNodeInfo(extractor)`                                                                                     |
| [`urls.go`](../netutil/urls.go)         | `ParseURLs`, `ParseURLsFromString`, `JoinURLs` (comma-separated)                                                         |
| [`net.go`](../netutil/net.go)           | `IsAddrInUse`; unexported `namedAddress`                                                                                 |

`NewNodeInfo` fails if hostname or local IPv4 cannot be determined.
`IsAddrInUse` uses `errors.Is(err, syscall.EADDRINUSE)`.
`IsPrivateAddress` is Deprecated (`net.IP.IsPrivate` / `IsLoopback` /
`IsLinkLocalUnicast`).

### `print`

[`print.go`](../print/print.go): write values to an `io.Writer`.

**Entry points:** `JSON`, `Yaml`, `Object`, `Print`, `Strings`, `Map`,
`Text`, `TextOneLine`, `RegisterType`, `FindRegistered*`.

**Invariants:**

- `Print` / `Object` order: registered `CustomFn`, then `Printer` interface,
  then `[]string` / `map[string]string`, else JSON. `Print` does not hold
  the registry lock across lookup or I/O.
- `Object(format)`: `"yaml"` / `"json"` take precedence over custom printers.
- `RegisterType` is process-global (mutex-protected).
- `JSON`/`Yaml` ignore marshal errors.
- Table `Map` truncates values to 80 bytes via `slices.StringUpto`.

### `slices`

| File                                 | Owns                                                                                                                                                                                                                                                                                       |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`slices.go`](../slices/slices.go)   | Equal/contains for common types, `CloneStrings` (preserves nil vs empty), `NvlString`, `Prefixed`/`Suffixed`/`Quoted`, `UniqueStrings`, generic `Deduplicate`/`Contains`/`Truncate`/`Replace`, `StringArrayToMap` (`key=value`), `HashStrings` (SHA-1, base64 raw URL), `StringsSafeSplit` |
| [`uint64s.go`](../slices/uint64s.go) | `Uint64s` for `sort.Sort`                                                                                                                                                                                                                                                                  |

Contains helpers are linear; comments say use a map for large sets.
`Replace` mutates in place. `Truncate` returns a prefix of the same slice
(not a copy) when shortened.
`Contains`, `*SlicesEqual`, `CloneStrings`, `NvlString`, `Uint64s`, and
`HashStrings` are Deprecated (stdlib `slices` / `cmp.Or`, or
`values.XXH3HashArgs128Hex` / SHA-256). `UniqueStrings` / `Deduplicate`
are kept (order-preserving; `slices.Compact` needs a sorted input).

### `ticker`

[`ticker.go`](../ticker/ticker.go): `New(ctx, every, value, run)` starts a
goroutine immediately. `Stop` cancels the derived context. `GetValue` /
`SetValue` / `Count` are mutex-safe. Nil receiver methods are no-ops.
The callback is **not** invoked on start, only on ticks. `count` increments
per tick, starting at 1 for the first callback.

### `urlutil`

[`urlutil.go`](../urlutil/urlutil.go): `GetQueryString`, `GetValue` (first
value or `""`), `GetPublicEndpointURL`.

Public URL: `X-Forwarded-Proto` overrides scheme; empty scheme → `https`;
host from `r.URL.Host` else `r.Host`. Path is the given relative endpoint
(query/fragment of the request are not copied).

### `values`

Coerce `any` and work with nested maps. Conversion helpers **do not return
errors**: unsupported types log at DEBUG and yield the zero value.

| File                                             | Owns                                                                                                                                                                                             |
| ------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`value.go`](../values/value.go)                 | `String`, `StringSlice`, `IntSlice`, `Bool`, `Time`, `Int`, `UInt64`, `Int64`, `Float32`, `Float64`, `JSON`, `JSONIndent`, `IndentJSON`, `IsEmpty`, `IsSlice`, `IsMap`, `IsCollection`, `Shrink` |
| [`mapany.go`](../values/mapany.go)               | `MapAny`, `FromStruct`, `FromJSON`, `FromYAML`, getters, `Extract`, `Merge`, `Shrink`, `RenameKeys`, `TraverseSubMaps`, SQL `Scan`/`Value`, `CanonicalJSON`, `CastMapAny`                        |
| [`coalesce.go`](../values/coalesce.go)           | `StringsCoalesce`, `NumbersCoalesce`, `Coalesce`, `Select`, `SelectFunc`                                                                                                                         |
| [`canonicaljson.go`](../values/canonicaljson.go) | `MarshalCanonicalJSON`, `CanonicalizeJSON`, `WriteCanonicalJSON` — object keys sorted; max depth 1000                                                                                            |
| [`hash_fast.go`](../values/hash_fast.go)         | `XXH3Hash*` — **not cryptographic**                                                                                                                                                              |

**Invariants:**

- `String` prefers `HasDisplayName`, then `HasName`, then `fmt.Stringer`.
- Slices join with `,`. `Bool` uses `"true"`/`"false"`.
- `FromJSON` / `FromYAML` swallow unmarshal errors (DEBUG log) and may
  return nil/partial maps. Empty / `{}` / `[]` / `null` (JSON) yield nil.
- `FromStruct` panics if the value is not a struct; skips unexported fields;
  keys are **Go field names**, not JSON tags.
- `MapAny.Add` skips empty keys, `nil`, `""`, and empty collections; it
  stores `false` and numeric zeros. May allocate if the receiver is nil.
- `MapAny.Slice` converts slices to `[]any`; missing or non-slice keys
  return nil (no panic).
- `Coalesce` on an empty vararg returns the zero value.
- Canonical JSON re-encodes unknown types via `encoding/json` then
  canonicalizes; structs keep `encoding/json` field order only after that
  round-trip (keys are then sorted).
- `XXH3HashArgs*`: strings are length-prefixed and followed by `0xff`.
  Unknown types go through `String(v)` with no length prefix.
- `StringsCoalesce` is Deprecated (`cmp.Or`).

## Test layout

| Location                                              | Role                                                     |
| ----------------------------------------------------- | -------------------------------------------------------- |
| `<pkg>/<file>_test.go`                                | Primary tests next to the code                           |
| `<pkg>/<file>_extra_test.go`                          | Additional coverage (often black-box `package foo_test`) |
| `<pkg>/stdlib_compare_test.go`, `*_benchmark_test.go` | Stdlib compatibility and performance comparisons         |
| `values/shrink_example_test.go`                       | Godoc examples                                           |
| `configloader/testdata/`                              | YAML configs, hostmap, overrides                         |
| `fileutil/testdata/`                                  | Sample files for resolve/reloader                        |

Mix of white-box (`package foo`) and black-box (`package foo_test`). Tests
use `github.com/stretchr/testify/assert` and `require`. There are no gomock
interfaces in this module.

CI (`.github/workflows/unittest.yml`) runs `make tools generate` then
`make build covtest`, and fails PRs under **80%** coverage. Markdown and
`Documentation/**` changes do not trigger that job.

Overridable time/ticker hooks for tests: `flake.NowFunc`, `format.NowFunc`,
and `reloader` `makeTicker`.

## Docs and samples that can go stale

Keep these compiling/accurate when behavior changes:

- Package `doc.go` comments (every package).
- [`README.md`](../README.md) package table.
- [`FINDINGS.md`](../FINDINGS.md) / [`ROADMAP.md`](../ROADMAP.md) — review
  IDs and modernization plan.
- [`configloader/README.md`](../configloader/README.md) — env var list and hostmap format.
- [`flake/README.md`](../flake/README.md) — bit layout and panic policy.
- Godoc examples under `values/`.

## When to update this file

Update in the **same change** when you:

- Add, remove, or rename a package or file that owns a concept.
- Add, move, or rename an exported entry-point type, func, or interface.
- Change an invariant (panic vs error, global, conversion zero-value,
  bit layout, interpolation syntax, callback/goroutine rules).
- Change internal package imports.
- Change test layout or testdata conventions.

Add a row to the concept index even for a small helper if someone would
otherwise grep for it.
