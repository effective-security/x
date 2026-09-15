# FINDINGS

Active, actionable bugs, security issues, and correctness problems.
Completed work and deprecation decisions are omitted. Modernization plans
live in [`ROADMAP.md`](ROADMAP.md).

Use the **ID** when commenting or assigning work. Update **Status** in the
same change as the code or decision, and remove the item when it is complete.

## Status

| Status | Meaning |
| --- | --- |
| Open | Not started |
| In Progress | Being fixed |
| Needs Approval | Behavior or compatibility change that needs a decision |

Severity: **security** > **bug** > **race** > **correctness** > **docs**.

## Index

| ID | Package | Title | Severity | Status |
| --- | --- | --- | --- | --- |
| [F-CFG-02](#f-cfg-02) | configloader | `file://` reads any path with no bound | security | Needs Approval |
| [F-CFG-04](#f-cfg-04) | configloader | `Load` mutates process environment | security | Needs Approval |
| [F-CTL-01](#f-ctl-01) | ctl | Bool tokens disagree with `values.Bool` | correctness | Needs Approval |
| [F-ENUM-02](#f-enum-02) | enum | `Parse` / `Convert` drop unknown names | correctness | Needs Approval |
| [F-REL-02](#f-rel-02) | fileutil/reloader | `NewReloader` always returns `nil` error | correctness | Needs Approval |
| [F-RES-02](#f-res-02) | fileutil/resolve | Directory mode `0744`; archived homedir | correctness | Needs Approval |
| [F-FLK-02](#f-flk-02) | flake | Failed machine-ID lookup becomes `0` | correctness | Needs Approval |
| [F-NET-04](#f-net-04) | netutil | `WaitForNetwork` has no context | correctness | Needs Approval |
| [F-NET-05](#f-net-05) | netutil | `namedAddress` is dead production code | correctness | Needs Approval |
| [F-NET-06](#f-net-06) | netutil | `FindFreePort` returns an unreserved port | correctness | Needs Approval |
| [F-PRT-02](#f-prt-02) | print | `JSON` / `Yaml` ignore marshal errors | correctness | Needs Approval |
| [F-TCK-01](#f-tck-01) | ticker | Zero interval panics; `Stop` does not wait | correctness | Needs Approval |
| [F-URL-01](#f-url-01) | urlutil | Public URL trusts forwarded proto/host | security | Needs Approval |
| [F-VAL-04](#f-val-04) | values | `Bool` inconsistent with `ctl` | correctness | Needs Approval |
| [F-VAL-07](#f-val-07) | values | JSON helpers swallow errors | correctness | Needs Approval |
| [F-VAL-09](#f-val-09) | values | `Shrink` drops `false` and numeric zero | correctness | Needs Approval |
| [F-X-02](#f-x-02) | repo | Process-global mutable hooks | correctness | Open |
| [F-X-03](#f-x-03) | repo | `golang-set` v1 only used in flake tests | docs | Open |
| [F-X-04](#f-x-04) | repo | Missing `doc.go` usage examples | docs | Open |

---

## `configloader`

### F-CFG-02

- **Status:** Needs Approval
- **Severity:** security
- **Package:** `configloader`

`file://` reads any path with no allow-list, size cap, or config-dir bind.
Bounding this can break callers who point at `/run/secrets` or absolute
certificates. The ROADMAP proposal adds an optional base directory and size
cap, with the restriction disabled until a major version.

### F-CFG-04

- **Status:** Needs Approval
- **Severity:** security
- **Package:** `configloader`

`Load` calls `os.Setenv(prefix+CONFIG_DIR, baseDir)`. Removing it can break
child processes that read that variable. Decide whether to keep and document
the process mutation or inject the value only into `Expander.Variables`.

---

## `ctl`

### F-CTL-01

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `ctl`

`BoolPtrMapper` accepts `true/1/yes` and `false/0/no` in any case.
`values.Bool` only treats `"true"` and `"yes"` as true. Aligning them would
change coercion of `"false"` and `"1"`; decide the shared token contract.

---

## `enum`

### F-ENUM-02

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `enum`

Unknown tokens OR as `0`, and whitespace is not trimmed. Add a strict parser
that returns `(E, error)` rather than changing the permissive `Parse` contract
in place.

---

## `fileutil`

### F-REL-02

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `fileutil/reloader`

`NewReloader` returns an `error` that is always nil. Decide whether to validate
the path, interval, and callback or retain the signature solely for source
compatibility.

### F-RES-02

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `fileutil/resolve`

Created directories use `0744`: group and other users can list names but
cannot traverse the directory. Changing this unusual mode to `0700` or `0755`
may affect callers. `homedir.Expand` errors are also ignored; replace the
archived dependency with `os.UserHomeDir` as described in the ROADMAP.

---

## `flake`

### F-FLK-02

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `flake`

Missing interfaces still yield machine ID `0`. Returning an error would make
`NewIDGenerator` panic under its existing policy and may affect Mac or CI
environments that currently rely on the fallback. Decide the failure contract.

---

## `netutil`

### F-NET-04

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `netutil`

`WaitForNetwork` polls once per second and cannot be cancelled independently.
Add a context-aware companion API or approve a signature change.

### F-NET-05

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `netutil`

`namedAddress` is unexported and used only by tests. Export it if it is part of
the intended raft integration, or delete it and its tests.

### F-NET-06

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `netutil`

`FindFreePort` closes its temporary listener before returning the port, so
another process can claim the port before the caller binds it. Add a
`ReservePort` API that returns the open listener; changing `FindFreePort`'s
return type would break callers.

---

## `print`

### F-PRT-02

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `print`

`JSON` and `Yaml` intentionally ignore marshal and writer errors. Add
error-returning companions such as `JSONErr` rather than silently changing the
existing API contract.

---

## `ticker`

### F-TCK-01

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `ticker`

`time.NewTicker(0)` panics, and `Stop` does not wait for an in-flight callback.
Decide whether to reject non-positive intervals with an error and whether to
add `Wait` or a blocking stop API.

---

## `urlutil`

### F-URL-01

- **Status:** Needs Approval
- **Severity:** security
- **Package:** `urlutil`

`GetPublicEndpointURL` trusts `X-Forwarded-Proto` and `Host`. Restricting these
values can break reverse-proxy deployments. The ROADMAP proposes an
`http`/`https` allow-list and an optional trusted-proxy policy.

---

## `values`

### F-VAL-04

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `values`

`Bool` does not use the same token set or case handling as `ctl.BoolPtrMapper`.
Resolve this with F-CTL-01.

### F-VAL-07

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `values`

`JSON`, `JSONIndent`, `FromJSON`, `FromYAML`, and `MapAny.YAML` swallow
serialization errors. Add error-returning JSON v2 APIs while keeping the
documented coercing helpers for compatibility.

### F-VAL-09

- **Status:** Needs Approval
- **Severity:** correctness
- **Package:** `values`

`MapAny.Shrink` uses `IsEmpty`, which treats `false` and numeric zeros as
empty. Removing those values can change a serialized document's meaning.
Add an option or a separate helper that preserves meaningful zero values.

---

## Cross-cutting

### F-X-02

- **Status:** Open
- **Severity:** correctness
- **Package:** repo

`SecretProviderInstance`, the `print` registry, and
`flake.DefaultIDGenerator` are process-global mutable state. Prefer instance
APIs where they already exist and design scoped replacements for the rest.

### F-X-03

- **Status:** Open
- **Severity:** docs
- **Package:** repo

`github.com/deckarep/golang-set` v1 is used only in flake tests. Replace it
with `map[uint64]struct{}` or `golang-set/v2`, then remove the dependency.

### F-X-04

- **Status:** Open
- **Severity:** docs
- **Package:** repo

`configloader`, `maps`, `values`, and `guid` package comments lack the short
usage examples required by the repository documentation rules.
