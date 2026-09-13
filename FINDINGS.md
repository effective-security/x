# FINDINGS

Bugs, security issues, and correctness problems found while reviewing each
package. Modernization ideas (JSON v2, generic methods, stdlib replacements)
are in [`ROADMAP.md`](ROADMAP.md).

Severity: **security** > **bug** > **race** > **correctness** > **docs**.

---

## `configloader`

### [security] `Marshal` writes files as `0777`

```113:113:configloader/load.go
	return os.WriteFile(fn, data, os.ModePerm)
```

`os.ModePerm` is `0777`. After umask this is typically world-readable and
often executable. Config files commonly hold secrets (`secret://`, tokens,
DB URLs). Write with `0600` (or at least `0644` and never executable).

### [security] `file://` reads any path with no bound

`ResolveValueWithSecrets` treats `file:///etc/passwd` (or `file://../../…`)
as a literal path and reads the whole file into memory. There is no
allow-list, `filepath.Clean`, size cap, or restriction to the config
directory. If any interpolated value can be influenced (hostmap overlay,
override file, env), this is arbitrary file read. Cap size and resolve
relative to the config base dir.

### [security] Failed `${secret://…}` lookups become empty strings

```37:45:configloader/expand.go
		s = os.Expand(s, func(env string) string {
			if strings.HasPrefix(env, SecretSource) && f.SecretProvider != nil {
				name := strings.TrimPrefix(env, SecretSource)
				sec, err := f.SecretProvider.GetSecret(name)
				if err != nil {
					logger.KV(xlog.ERROR, "secret", name, "err", err.Error())
				}
				return sec
			}
```

A secret-backend error is logged and replaced with `""`. The later
`${` check does not fire. Callers get a successfully-loaded config with
blank credentials. Return the error instead.

If `SecretProvider` is nil, `${secret://name}` falls through to
`os.Getenv("secret://name")`.

### [security] `Load` mutates process environment

```118:122:configloader/configloader.go
	envName := f.envPrefix + "CONFIG_DIR"
	if variables[envName] == "" {
		variables[envName] = baseDir
		os.Setenv(envName, baseDir)
	}
```

This races with concurrent loads and leaks the config directory into every
child process. Keep the value on `Expander.Variables` only.

### [bug] Interpolation does not apply inside most maps

Only the concrete type `map[string]string` is rewritten in place. The type
check uses `v.Type().String() == "map[string]string"`, so a defined alias
(`type Env map[string]string`) is skipped.

For every other map, `MapRange` values are not addressable. Nested
`map[string]any`, `map[string]struct{…}`, and slices of those structs are
walked but never updated. `ExpandAll` returns nil while `${VAR}` / `env://`
remain.

Fix: copy map entries out, expand, `SetMapIndex` back. Match maps by key
and elem kinds, not `Type().String()`.

### [bug] `Environment` field type-assert can panic

```112:114:configloader/configloader.go
	} else if value, err := reflections.GetField(config, "Environment"); err == nil {
		environment = value.(string)
	}
```

A non-string `Environment` (pointer, enum, `any`) panics. Use a type switch.

### [bug] `userName()` panics on `user.Current` failure

Not listed with the package’s documented panics. A container without
`/etc/passwd` cannot load config. Return `"unknown"` or a wrapped error.

### [correctness] Unresolved-variable comment does not match code

`ExpandAll`’s comment says missing `${VAR}` becomes empty (like
`os.Getenv`). The second `strings.Contains(s, "${")` check then errors —
but only if expansion *reintroduced* `${`. A missing name is silently
empty. Pick one contract and document it.

### [docs] README install path is stale

`configloader/README.md` still says `github.com/effective-security/x/pkg/configloader`.
The module path is `github.com/effective-security/x/configloader`.

---

## `ctl`

No functional bugs. `BoolPtrMapper` and `values.Bool` disagree on accepted
tokens (`1`/`0`/`false`/`no` vs only `true`/`yes`). Align them or document
the split.

---

## `enum`

### [bug] Flag helpers never see bit 31

```128:128:enum/enum.go
	for i := E(1); i > 0 && i <= val; i <<= 1 {
```

`E` is `~int32`. After `1<<30`, the next shift is negative, so `i > 0`
stops. A flag at `1<<31` is invisible to `FlagNames` / `Flags` /
`FlagsInt`. Walk `uint32` bits, or document that only bits 0–30 are flags.

### [correctness] `Parse` / `Convert` drop unknown names

Unknown tokens OR as `0`. `"Low,Medum"` equals `Low` with no error.
Whitespace is not trimmed (`"Low, Medium"` misses `Medium`). Return an
error (or a `ParseStrict`) and `strings.TrimSpace` each token.

---

## `fileutil`

### [race] `Reloader.Close` writes under `RLock`

```119:128:fileutil/reloader/reloader.go
	k.lock.RLock()
	defer k.lock.RUnlock()

	if k.closed {
		return errors.New("already closed")
	}

	k.closed = true
	k.stopChan <- struct{}{}
```

`closed` is mutated while holding a read lock. Two concurrent `Close`
calls can both pass the check; the second send blocks forever because the
goroutine already exited.

`fileModifiedAt` is written on the ticker goroutine without the lock and
read from `Reload` — `go test -race` should flag it.

### [correctness] `NewReloader` always returns `nil` error

The signature promises an error it never produces. Either validate
`filePath` / `checkInterval` / callback, or drop the error.

### [docs] Tests still use deprecated `io/ioutil`

`fileutil/reloader/reloader_test.go` should use `os.WriteFile`.

### `fileutil/resolve`

- Comment says `NewNotFound`; the helper does not exist.
- Typo: `"crerate dir"`.
- `File` with a relative path and empty `baseDir` stats `""`.
- `ExpandPath` ignores `homedir.Expand` errors (`~` left in place).
- Directory mode is hardcoded `0744` (group-writable). Prefer `0755` or
  a caller-supplied mode.

`github.com/mitchellh/go-homedir` is archived; `os.UserHomeDir` is enough.

---

## `flake`

### [bug] `IsTestRun` matches any path containing `test` or `debug`

```199:201:flake/flake.go
func IsTestRun() bool {
	return strings.Contains(os.Args[0], "test") || strings.Contains(os.Args[0], "debug")
}
```

Binaries under `…/testing/…`, `contest`, `debug-server`, etc. use the FNV
program-path machine ID instead of the interface address. That can collide
across hosts and silently change ID space in production.

Detect tests with `testing.Testing()` (Go 1.21+) or a build tag, not a
substring of `os.Args[0]`.

### [correctness] Failed machine-ID lookup becomes `0`

`defaultMachineID` logs and returns `(0, nil)` when `InterfaceAddrs`
fails or only loopback exists. Every such process shares machine ID `0`.
That breaks the uniqueness contract. Return an error and let
`NewIDGenerator` panic (existing policy) or refuse to start.

### [docs] Package comment and `Settings` comment are wrong

The file header still says 39-bit time / 10 ms / 8-bit sequence. Actual
layout (and the README / codemap) is 41-bit time in 1 ms, 6-bit sequence,
16-bit machine ID. `Settings` still says “lower 8 bits of the private IP”.

---

## `format`

### [bug] `TextOneLine` treats the first **byte** as a rune

```275:275:format/format.go
				if !prevPartDot && unicode.IsUpper(rune(part[0])) {
```

A non-ASCII first character is misclassified. Use `utf8.DecodeRuneInString`.

### [correctness] `StringMax` / `StringsMax` count bytes, not runes

`val[:limit]` can slice mid-code-point and emit invalid UTF-8. Use
`utf8` or `strings.ToValidUTF8`.

### [correctness] `ParseStringTime` keys off string **length**

RFC3339 / RFC3339Nano / zone offsets have several legal lengths. A valid
timestamp of an unexpected length falls through to `RFC3339Nano` and
becomes zero on failure (error discarded). Try a list of layouts instead.

---

## `guid`

No bugs in the current contract (raw 16 bytes, uppercase hex, panic on
`crypto/rand` failure). It is **not** RFC 9562 / 4122 — version and variant
bits are unset. Go 1.27 adds stdlib `uuid`; see ROADMAP.

---

## `maps`

### [correctness] `SyncMap` is a partial, panic-prone wrap of `sync.Map`

`Load` / `Range` type-assert and panic if the map was corrupted or used
with mixed types. Stdlib `sync.Map` has been generic (`sync.Map[K,V]`)
since Go 1.24 and includes `Swap`, `CompareAndSwap`, `CompareAndDelete`,
and `Clear`. This wrapper is behind the stdlib.

`Take` / `Drop` / `Range` use random map iteration (documented).
`OrderedKeys` still uses `sort.Slice` instead of `slices.Sort`.

---

## `netutil`

### [bug] `FindFreePort` closes a failed listener

```29:36:netutil/freeport.go
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			l.Close()
			logger.KV(xlog.ERROR,
				"reason", "unable to listen",
				"addr", addr,
				"err", err.Error())
			continue
		}
```

On listen failure `l` is nil. `Close` happens to be nil-safe today, but
the call is wrong and hides the real error. Do not close; return or
continue with `err`.

The returned port is already closed, so the caller still has a TOCTOU
race. Document that, or return the open `net.Listener`.

### [bug] `IsAddrInUse` does not unwrap

```47:54:netutil/net.go
func IsAddrInUse(err error) bool {
	if err, ok := err.(*net.OpError); ok {
		if err, ok := err.Err.(*os.SyscallError); ok {
			return err.Err == syscall.EADDRINUSE
		}
	}
	return false
}
```

`errors.Wrap` / `%w` / `net.OpError` behind another wrapper all return
false. Use `errors.As` and `errors.Is(err, syscall.EADDRINUSE)`.

### [correctness] `IsPrivateAddress` reinvents `net.IP` helpers

Custom CIDR table includes loopback and link-local, so the name is
misleading. `127.0.0.1/8` should be `127.0.0.0/8`. Prefer
`ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast()`, and decide
whether CGNAT `100.64.0.0/10` counts.

### [correctness] `WaitForNetwork` has no context

One-second poll, no cancel, no jitter. After timeout it returns the last
error (possibly nil IP + err). Accept `context.Context`.

### [correctness] `namedAddress` is dead production code

Unexported, only referenced from tests. Export it (the comment mentions
raft identity) or delete it.

---

## `print`

### [race] `Print` deadlocks with `RegisterType`

```84:88:print/print.go
func Print(w io.Writer, value any) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	if printFunc, found := FindRegistered(value); found {
```

`FindRegistered` takes the same `RLock`. `sync.RWMutex` is not reentrant:
if `RegisterType` is blocked in `Lock`, the nested `RLock` waits forever.

`Print` also holds the lock across JSON/table I/O. Lookup under the lock,
copy the func, unlock, then print.

### [correctness] `JSON` / `Yaml` ignore marshal errors

Documented, but a failed marshal writes nothing (or a partial buffer) with
no signal. At least write a fallback or log at ERROR, matching `values`.

---

## `slices`

### [security] `HashStrings` is not a safe digest

```281:288:slices/slices.go
func HashStrings(values ...string) string {
	h := crypto.SHA1.New()
	for _, v := range values {
		_, _ = h.Write([]byte(v))
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
```

SHA-1 is retired for integrity. Adjacent strings are concatenated with no
length prefix, so `HashStrings("ab", "c") == HashStrings("a", "bc")`.
Length-prefix (as `values.XXH3HashArgs*` does) and use SHA-256 or XXH3,
and say it is non-cryptographic if that is the intent.

### [correctness] Overlap with stdlib `slices`

`Contains`, `Equal`, `Clone`, `Compact` already exist in `slices`. The
typed `*SlicesEqual` helpers and `Uint64s` (`sort.Interface`) can be
thin wrappers or deprecated.

`StringArrayToMap` and `UniqueStrings` use `errors.New` from stdlib, not
`github.com/cockroachdb/errors`, against the repo rule.

---

## `ticker`

### [correctness] Zero interval panics; `Stop` does not wait

`time.NewTicker(0)` panics. Validate `every > 0`.

`Stop` cancels the context but does not wait for an in-flight callback.
A callback that uses `GetValue` / `SetValue` after `Stop` is fine, but
callers cannot know the last tick has finished. Offer `Stop` + wait, or
document it.

A panicking `run` kills the goroutine with no recovery or log.

---

## `urlutil`

### [security] Public URL trusts `X-Forwarded-Proto` and `Host`

```33:51:urlutil/urlutil.go
	if specifiedProto := r.Header.Get(XForwardedProtoHeader); specifiedProto != "" {
		proto = specifiedProto
	}
	host := r.URL.Host
	if host == "" {
		host = r.Host
	}
```

Any client can set `X-Forwarded-Proto: javascript` (or `http`) and a
hostile `Host`. The result is used as a “public” absolute URL (password
resets, redirects, links). Allow only `http`/`https`, and take host from
a configured public base unless a trusted proxy already sanitized
`X-Forwarded-*`.

---

## `values`

### [bug] `MapAny.Add` drops `false`, `0`, and other zeros

`Add` skips `IsEmpty(value)`. `IsEmpty` uses `reflect.Value.IsZero()`, so
`Add("enabled", false)` and `Add("count", 0)` are no-ops. `Shrink` also
strips zeros, which can drop meaningful JSON fields.

Treat empty as nil / `""` / empty collection; do not treat `false` or `0`
as empty unless a separate option says so.

### [bug] `MapAny.Slice` panics on the wrong type

```223:228:values/mapany.go
func (c MapAny) Slice(k string) []any {
	if c == nil {
		return nil
	}
	return c[k].([]any)
}
```

A JSON array unmarshaled as `[]string` (or a missing key) panics. Type-
assert and return nil, or convert via `[]any`.

### [correctness] `Time` string parsing is too narrow

Strings longer than 20 characters must match
`2006-01-02T15:04:05.000-0700` (no colon in the zone). RFC3339
(`…+00:00`) fails. Shorter strings are parsed only as unix integers, so
`2006-01-02` and `DateTime` become nil. Share `format.ParseStringTime`
or a layout list.

### [correctness] `Bool` is inconsistent with `ctl`

Only `"true"` and `"yes"` (case-sensitive). `"True"`, `"1"`, `"false"`
are false / unsupported. Document or align.

### [correctness] `Coalesce` panics on no args

`return args[0]` on an empty vararg. Documented; still surprising. Return
the zero value.

### [docs] `XXH3HashArgs*` comments disagree with the code

Comments (and the old codemap text) say string fields end with a `0x00`
delimiter. The implementation uses `hashDelimiter = []byte{0xff}` and
length-prefixes strings. Callers who reimplement the scheme from the
comment will not match.

The `default` branch hashes `String(v)` with **no** length prefix, so
different types that stringify the same way collide.

### [correctness] JSON helpers swallow errors

`JSON`, `JSONIndent`, `FromJSON`, `FromYAML`, `MapAny.YAML` ignore
marshal/unmarshal errors (documented). `FromJSON` can return a **partial**
map after a mid-document failure. Prefer `encoding/json/v2` with
`RejectUnknownMembers` / duplicate-name rejection for new APIs (ROADMAP).

---

## Cross-cutting

| Issue | Where |
| --- | --- |
| Stale “Go 1.24 or newer” in root README; `go.mod` is 1.27 | `README.md` |
| Process-global mutable hooks (`SecretProviderInstance`, `print` registry, `flake.DefaultIDGenerator`) | several packages |
| Unused / outdated dependency: `github.com/deckarep/golang-set` v1 only in `flake` tests | `go.mod` |
| Several packages have no usage example in `doc.go` despite non-obvious entry points | `configloader`, `maps`, `values`, `guid` |
