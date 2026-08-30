# Repository Instructions

## Project role

`github.com/ingot-agent/ingot-abi` is the fixed, non-replaceable host ABI
between ingot Plugins and generated Runtime Images. Unlike the Agent SDK it
is not optional: the Builder recognizes exactly these module identities,
pins the exact version, and generates the host implementations behind every
contract in this module.

## Minimal surface is mandatory

Only protocols that satisfy every Runtime ABI rule belong here:

1. the resource or semantic is exclusively owned by the generated runtime;
2. the protocol must exist before the Component Graph is constructed, or it
   governs the whole graph lifecycle;
3. a Plugin cannot provide a semantically equivalent replacement;
4. every Runtime Image needs uniform, predictable behavior;
5. without this module the Builder would need plugin-name special cases,
   globals, or hidden channels.

Never add HTTP clients, filesystem access, loggers, metrics, tracing,
clocks, schedulers, event buses, secrets, business configuration,
model/tool/session contracts, UI protocols, service locators, or any
`Get(string)` registry. Do not depend on the Agent SDK, ingot core, or any
concrete plugin; prefer the Go standard library.

## Contract discipline

- Treat exported types, interfaces, sentinel errors, and documented
  semantics as compatibility commitments pinned by the Builder.
- Document ownership, concurrency, ordering, cancellation and error
  semantics for every exported identifier.
- Aggregates are immutable by contract; callers own returned aggregate
  outputs unless documented otherwise.
- Preserve ordinary Go error chains for `errors.Is` and `errors.As`.
- Keep the module dependency-free; adding a dependency is a breaking change
  to the Runtime ABI.

## Change discipline

- The Builder pins exact versions. Bumping the ingot ABI requires a
  coordinated Builder + ingot ABI release.
- Format only changed Go files with `gofmt` and avoid unrelated rewrites.
- Preserve existing work and keep changes focused on the requested task.

## Git workflow

- Never edit, commit, or push directly on `main`. Use a focused working
  branch and submit changes through a pull request.
- Do not create commits, push branches, force-push, rewrite shared history,
  or open pull requests unless the user explicitly requests that action.

## Required validation

Before every commit, run the complete race-enabled test suite:

```sh
go test -race ./...
git diff --check
```
