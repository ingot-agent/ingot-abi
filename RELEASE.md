# ABI releases and Core coordination

ABI releases are Go module tags for `github.com/ingot-agent/ingot-abi`.
Core pins an exact version and source identity; installing a newer ABI module
does not update the ABI used by an existing Builder or Runtime Image. This
repository currently has no checked-in release workflow.

The checked-in Core pins **`v0.1.0`**, which is also the ABI tag present at the
time of this documentation update. Check the target Core release's
[resolve.go](https://github.com/ingot-agent/ingot/blob/main/internal/builder/resolve.go)
and its generated lock instead of assuming a future release uses this pin.

## Identify the compatibility boundary

| Change | Required review |
|---|---|
| Documentation correction with unchanged semantics | Correct examples and references; no invented API release or version bump |
| New exported shape or helper | API/semantic review, module tests and Core acceptance of the exact intended ABI version |
| Changed ownership, concurrency, errors, cancellation or lifecycle | Treat as a compatibility change even if callers still compile |
| Wrapper or host contract identity change | Coordinated Core type recognition, code generation and consumer migration |
| Constructor shape change in Core | Review Builder compatibility and generated wiring; an unchanged ABI tag alone does not identify constructor acceptance |

Do not silently reinterpret a published tag. The module remains pre-1.0;
publish explicit migration notes for incompatible changes and identify the
required Core and plugin versions. No multi-version host-ABI compatibility or
security-backport schedule is implied by this document.

## Prepare a release candidate

1. Review [CONTRIBUTING.md](CONTRIBUTING.md) and the host-contract admission
   rules. Keep `go.mod` free of SDK, Core and concrete plugin dependencies.
2. Record every public shape or semantic change, the previous and new behavior,
   affected consumers, and whether state written by affected plugins remains
   readable. The ABI supplies locations and lifecycle, not plugin migrations.
3. Run `GOWORK=off go vet ./...`, `GOWORK=off go test -race ./...` and
   `git diff --check` from this repository. Review test results, not just command
   availability.
4. In a separate Core checkout, prepare the corresponding exact ABI pin and
   changes to [graph validation](https://github.com/ingot-agent/ingot/blob/main/internal/builder/graph.go)
   and [generated runtime support](https://github.com/ingot-agent/ingot/blob/main/internal/builder/generate.go)
   when affected. Update pin-dependent fixtures and tests found by searching
   for `IngotABIVersion` and the old version, rather than changing only one
   constant.
5. Verify all exposed host types and wrappers with generated components:
   successful construction, missing/invalid dependencies, forbidden host
   exports, exact constructor shape, normal shutdown, failure cleanup and
   reverse cleanup order. Include invocation/check mode and state-location
   behavior when changed. A handwritten function in ABI tests is insufficient
   evidence for these Builder behaviors.
6. Test the coordinated candidate in an explicit, temporary development
   workspace. Record this as local-source integration; it does not establish
   that released modules can be resolved independently.
7. Prepare release notes and related Core/plugin pull requests. Confirm the
   repository's [security reporting setup](SECURITY.md) before public release.

## Publish and verify released inputs

Publication requires the maintainer's normal review and release authorization.
After the ABI release change is merged, create its `vX.Y.Z` tag on the reviewed
commit and publish release notes. Core and consumers must refer to that exact
version, with matching semantic import paths if a future major version requires
them. Do not move an existing tag to correct a bad release.

From an empty temporary directory, verify availability without a local
workspace replacement. Replace the placeholder with the actual release:

```sh
GOWORK=off go mod download -json github.com/ingot-agent/ingot-abi@vX.Y.Z
```

Check the returned version and checksums through the intended Go module proxy.
Then rerun affected consumer checks with `GOWORK=off`, and use the coordinated
Core candidate to resolve and build its official profiles in a fresh managed
Home. The resulting lock must select the intended ABI path/version/sum;
attempted MVS upgrades or production replacements must still be rejected.

The publication sequence follows actual dependencies: publish the ABI version
first, then consumers that require it, then the Core/profile changes that refer
to those available versions. An SDK release is needed only if SDK contracts
also change; the SDK is not an ABI dependency.

Core's existing immutable images retain their compiled code. Users must build
new images to adopt the new ABI. Changing a Runtime binding does not migrate
plugin state or restart a running process automatically. Link the relevant
[Core upgrade guide](https://github.com/ingot-agent/ingot/blob/main/docs/UPGRADING.md)
and plugin migration instructions in release notes.

## Failure handling

If the module tag is unavailable from the intended proxy, or independent
consumers fail, stop dependent releases and investigate the exact recorded
version/commit. Workspace replacements must not be used to conceal unresolved
published dependencies.

For an already-published defective version, prepare a new reviewed release and
explicit consumer updates. If retreating to a previous Core/Image, check plugin
state compatibility and restore a validated backup when necessary; ABI rollback
does not roll back data. Preserve failing test evidence and the exact dependency
graph in a redacted issue or private security report as appropriate.
