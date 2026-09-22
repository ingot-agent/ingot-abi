# ingot ABI

> The fixed, versioned Go contract between ingot Plugins and generated Runtime
> Images.

[中文](./README.zh.md)

`github.com/ingot-agent/ingot-abi` defines the small set of public types that
the ingot Builder recognizes by exact Go type identity. It separates
runtime-owned host behavior from the replaceable agent contracts in
[`github.com/ingot-agent/sdk`](https://github.com/ingot-agent/sdk).

This is a source-level Go ABI for statically generated Component wiring. It is
not a dynamic-linking ABI, a plugin loader, or a general-purpose framework.

## Why this module exists

Most contracts exchanged by Components are ordinary, replaceable capabilities:
models, tools, sessions, prompts, filesystems, and interaction channels. Those
contracts can live in the Agent SDK, another domain SDK, or a Plugin's own
module.

A few contracts are different. The generated Runtime Image owns process
lifecycle, invocation mode, Plugin state allocation, and the exact Component
constructor shape. Plugins cannot provide equivalent replacements for those
semantics. Keeping them in this dedicated module makes that privilege explicit,
minimal, versioned, and auditable.

## Packages

| Package | Purpose |
|---|---|
| `ingotabi` | Component ABI: `Cleanup`, `Optional[T]`, `Named[T]`, and their helpers. |
| `invocation` | Read-only Runtime arguments and `ModeRun`/`ModeCheck`. |
| `lifecycle` | Graceful process shutdown requests through `Controller`. |
| `state` | The absolute, Plugin-scoped persistent state location. |

The module depends only on the Go standard library.

## Component ABI

Every Component constructor has this exact shape:

```go
package component

import (
	"context"

	ingotabi "github.com/ingot-agent/ingot-abi"
)

type Dependencies struct{}
type Exports struct{}

func New(
	ctx context.Context,
	deps Dependencies,
) (Exports, ingotabi.Cleanup, error) {
	return Exports{}, nil, nil
}
```

The current Core Builder requires the two-argument constructor above. Plugins
load their own configuration through an explicit `state.Scope` dependency;
there is no generated global configuration decoder or `Config` argument.
See the [Core file-format reference](https://github.com/ingot-agent/ingot/blob/main/docs/FILE_FORMATS.md)
and [ADR 0003](https://github.com/ingot-agent/ingot/blob/main/docs/adr/0003-plugin-configuration.md).

The Builder also recognizes `ingotabi.Optional[T]` as an optional dependency
and `ingotabi.Named[T]` as a stable runtime-instance identity. These wrappers
are part of the Component ABI; similarly shaped types from another module are
ordinary contract types.

## Host contracts

Host capabilities are declared explicitly in a Component's `Dependencies`:

```go
import (
	"github.com/ingot-agent/ingot-abi/invocation"
	"github.com/ingot-agent/ingot-abi/lifecycle"
	"github.com/ingot-agent/ingot-abi/state"
)

type Dependencies struct {
	Invocation invocation.Invocation
	Lifecycle  lifecycle.Controller
	State      state.Scope
}
```

The generated Runtime Image injects these exact types as virtual host
providers. They do not add Component-to-Component graph edges, and Plugins may
not export them.

### Invocation

- `Arguments` returns a caller-owned copy and excludes ingot-owned flags.
- `ModeRun` is normal execution.
- `ModeCheck` constructs, validates, and cleans up the complete graph without
  starting an interaction loop or retaining external resources.

### Lifecycle

- `RequestShutdown(nil)` expresses normal completion intent.
- `RequestShutdown(err)` records a process-level failure.
- The first request cancels the Runtime context.
- Requests are concurrent-safe and non-blocking.
- Non-nil shutdown causes and Cleanup errors are aggregated; the generated
  Runtime chooses the final process result.

### State

- `Scope.Dir()` returns an absolute directory assigned from Plugin identity.
- Components in the same Plugin share the same state directory.
- The Plugin owns file access, schema validation, migration, and persistence.

## Builder invariants

- The Builder fixes the module path and exact supported version.
- The selected module version and source identity enter the lock file and Image
  ID.
- Production builds reject user-selected paths, versions, and replacements.
- Only exact ingot ABI type identities receive host injection or wrapper
  semantics.
- Host dependencies remain visible in graph inspection.

## What does not belong here

Do not add HTTP clients, filesystems, loggers, metrics, tracing, clocks,
schedulers, event buses, secrets, business configuration, model/tool/session
contracts, UI protocols, service locators, or extensible registries. Those are
replaceable capabilities or application concerns.

## Compatibility and versioning

The Builder pins this module exactly. Changes to exported shapes or documented
ownership, concurrency, ordering, cancellation, or error semantics require a
coordinated ingot ABI and Builder release. New APIs must satisfy every host ABI
admission rule and keep the module dependency-free.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for contract admission and review rules,
[RELEASE.md](RELEASE.md) for ABI/Core release coordination, and
[SECURITY.md](SECURITY.md) for reporting security concerns.

```sh
go test -race ./...
go vet ./...
```

## License

[Apache License 2.0](./LICENSE)

## Design history

The [v0.1 proposal](docs/design-history/README.md) was migrated from Core for
traceability. It preserves its original MIT notice and is not the current
constructor or configuration specification.
