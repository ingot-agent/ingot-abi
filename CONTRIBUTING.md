# Contributing to the Ingot ABI

The ABI is the fixed source-level Go contract between plugins and generated
Runtime Images. Its module path and exact version are recognized by the Core
Builder. Changes here affect independently developed plugins and generated
programs, so both exported types and documented semantics are compatibility
commitments.

Read the [overview](README.md), [release coordination guide](RELEASE.md) and
[security policy](SECURITY.md) before proposing a change. Ordinary bugs and
documentation fixes can be reported through the repository's issue templates.

## Decide whether the change belongs here

An added host contract must satisfy all of these conditions:

1. The resource or semantic is exclusively owned by the generated runtime.
2. It is needed before graph construction or governs the whole graph lifecycle.
3. A plugin cannot supply an equivalent replaceable capability.
4. Every Runtime Image needs consistent semantics for it.
5. Without it, plugins would need identity-based special cases, global mutable
   state or hidden communication with the Builder/runtime.

Current contracts are Component wrappers, invocation metadata, process shutdown
and plugin state location. Keep the module dependent only on the Go standard
library. HTTP, tools, sessions, models, metrics, secrets, policy and application
protocols belong in ordinary contract modules or plugins. A contract needed by
one official plugin is not automatically an ABI concern.

For new exports or semantic changes, open a proposal describing the behavior,
why a plugin capability cannot implement it, compatibility impact, and the Core
and plugin changes required. Maintainers should agree on this boundary before
implementation expands it.

## Contract review

For every changed exported type or function, document and test:

- accepted inputs and zero/nil behavior;
- caller/callee ownership of slices, maps, pointers and returned values;
- concurrency, ordering and whether methods may block;
- cancellation, deadlines and recognizable errors;
- construction, runtime and cleanup lifetime;
- which behaviors the generated Core runtime must implement.

Preserve ordinary Go error chains and test `errors.Is`/`errors.As` where
applicable. Do not expose plugin settings through `context.Value`, make a
plugin name special, or add a service locator. SDK contracts are optional and
replaceable; the ABI must remain independent of them.

Current constructor validation belongs to
[Core graph.go](https://github.com/ingot-agent/ingot/blob/main/internal/builder/graph.go).
Core currently requires `New(context.Context, Dependencies) (Exports,
ingotabi.Cleanup, error)`. A function that compiles in this module is not by
itself proof that Core accepts it as a component. Review nominal type identity,
generated calls and initialization/cleanup behavior in the consuming Builder.

## Development workflow

Use Go 1.24 or newer, as declared by [go.mod](go.mod). Create a focused branch
from the intended base and preserve existing local work. Implement only the
agreed contract change; format changed Go files with `gofmt`.

From the repository root, in a POSIX shell:

```sh
GOWORK=off go vet ./...
GOWORK=off go test -race ./...
git diff --check
```

In PowerShell:

```powershell
$env:GOWORK = 'off'
go vet ./...
go test -race ./...
git diff --check
```

The race detector requires a supported platform and C toolchain. If unavailable,
report that limitation and run the required race suite in a supported
environment before committing. A plain `go test` result does not replace it.
Workspace-disabled tests check the declared module independently; development
workspace integration is a separate check.

Use external-package tests to model an independent plugin author. They should
protect public shapes and observable semantics. Generated implementations live
in Core, so changes to host injection, wrapper recognition, cancellation or
cleanup also require the matching Core tests and a generated-runtime exercise.
See [RELEASE.md](RELEASE.md) for the coordination sequence.

## Pull requests

A pull request should explain the concrete behavior change, whether old plugin
source still builds, any semantic incompatibility, and the exact tests run.
Link related Core/plugin work when required. Update affected Go documentation,
README examples and migration instructions in the same change. Do not present
unmerged counterpart changes or an unpublished tag as already available.

Keep commits and releases focused. Do not modify or move published version
tags. Root repository version tags use `vX.Y.Z`, without a nested module prefix.
Publishing a tag and changing Core's pin are maintainer release actions, not a
side effect of editing documentation or running tests.

## License

Contributions to the ABI are under [Apache License 2.0](LICENSE). The migrated
documents under [docs/design-history](docs/design-history/README.md) retain
their original MIT license and attribution; preserve those notices.
