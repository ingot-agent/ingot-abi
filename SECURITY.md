# Security policy: Ingot ABI

## Reporting a suspected vulnerability

Do not post exploit details, credentials, private prompts, state files or
sensitive reproduction data in a public issue or pull request.

**Reporting setup checked on 2026-09-22:** GitHub's API reported private
vulnerability reporting disabled for this repository. No dedicated security
email or guaranteed response schedule is published in this checkout.

Open the repository's [Security page](https://github.com/ingot-agent/ingot-abi/security).
If maintainers have since enabled **Report a vulnerability**, use that private
form. Otherwise, open a [contact request](https://github.com/ingot-agent/ingot-abi/issues/new?title=Private%20security%20reporting%20contact%20request)
containing only: “Please provide a private channel for a security report.”
Do not include affected code paths, reproduction steps or attachments in that
public coordination request. Wait for maintainers to establish a private
channel before sending technical details.

Once a private channel is available, include:

- exact ABI version, Core/Builder version and pin, affected contract, Go version, OS/architecture, and a minimal generated-component reproducer;
- the impact and the access/conditions needed to reproduce it;
- the smallest reproduction with synthetic data and credentials;
- expected versus observed behavior, and any suggested mitigation.

If the report spans repositories, name all affected modules in one private
report so maintainers can coordinate it. Ordinary non-security defects belong
in the public bug-report form.

## Scope and trust boundaries

The ABI defines Component wrappers and runtime-owned invocation, lifecycle
and state-location contracts. Report ambiguities or contract defects affecting
host injection, cancellation, cleanup, state separation or ownership here.
Generated host implementations live in [Core](https://github.com/ingot-agent/ingot);
include its version when reporting runtime behavior.

The ABI is a source-level Go contract, not a security boundary around native
plugin code. A state scope assigns a location; it does not restrict filesystem
access, encrypt contents or migrate data. A plugin's code runs with the runtime
account's permissions. The exact ABI version pinned by a Builder must be
considered together with that Builder's generated implementation.

## Version information and disclosure

Please report the exact affected versions, including older releases and source
checkouts. This repository does not currently publish an LTS or guaranteed
security-backport schedule. Fix availability and migration requirements must
be stated in the corresponding release notes; an unreleased branch fix should
not be described as available in an existing tag.

Use the established private channel to coordinate investigation and disclosure.
Do not assume this document guarantees a response deadline or authorizes
testing systems, accounts or data that you do not control.

## Maintainer release requirement

Before public release, enable **Private vulnerability reporting** in the GitHub
repository's security settings, verify the reporter-facing form with an
appropriate account, and update the dated setup statement above. Monitor the
chosen channel and document any support/response policy only after it has been
agreed. Adding this file or an issue-template link does not enable reporting
in GitHub settings.

