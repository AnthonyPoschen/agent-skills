# Skill Merge Batch 1

## Why

Every installed skill sends its name and description to the agent on every
turn. Fewer, broader skills cost less context, but only if each old use still
triggers. Anthony's main concern was that merging would stop skills from
triggering, so this batch was gated on a trigger test.

## What Merged

| Parent | Absorbed | Reason |
| --- | --- | --- |
| `coding` | `no-comments`, `type-system-discipline` | Both only apply during code work, which `coding` already always triggers on. |
| `git` | `git-commit-workflow`, `git-worktree`, `release` | Commits, worktrees, and releases happen together. |
| `flux` | `flux-ghcr-image-automation`, `flux-kustomize-layout` | Same tool and repo; tasks often span both. |
| `writing` | `technical-writing`, `asd-ste100` | Both govern prose for people; they already deferred to each other. |

Each old skill body is now a reference file under the parent and loads only
when the parent's routing section sends the agent there. Each parent
description keeps the old skill name and its key trigger phrases, so a prompt
that named or described an old skill still matches.

Kept separate on purpose: `grill-with-docs` (revisit later),
`in-your-own-words`, `act-within-scope`, `dockerfile`, `makefile`.

## Trigger Test

Method: the `skill-evaluation` triggering check. 24 realistic prompts (three
per absorbed skill) and four negatives were routed against the full skill
listing on `master` and again on this branch. The router was an agent reading
only the listing, as a harness does. "Expected" is the skill that must load.

| # | Prompt | Before (master) | After |
| --- | --- | --- | --- |
| 1 | Go through my diff and strip out the pointless comments before I open the PR | no-comments + coding | coding (no-comments ref) |
| 2 | There are way too many comments in src/billing, clean them up | no-comments + coding | coding (no-comments ref) |
| 3 | Review this file for comments that just narrate the code | no-comments + coding | coding (no-comments ref) |
| 4 | Design the TypeScript types for an order that can be pending, paid, or refunded | type-system-discipline + coding | coding (types ref) |
| 5 | Review this signature: fn transfer(from: String, to: String, amount: f64) | type-system-discipline + coding | coding (types ref) |
| 6 | We keep swapping userId and orderId; make the compiler catch it | type-system-discipline + coding | coding (types ref) |
| 7 | commit these changes | git-commit-workflow | git |
| 8 | split my working tree into sensible commits | git-commit-workflow | git |
| 9 | write a commit message for the staged diff | git-commit-workflow + writing-type skills | git + writing |
| 10 | do this in a new worktree | git-worktree | git |
| 11 | clean up the worktrees for merged branches | git-worktree | git |
| 12 | clone my repo and set up a feature checkout for the auth rewrite | git-worktree | git |
| 13 | cut a release of agent-skills | release | git |
| 14 | update the changelog for the new version | release | git |
| 15 | the v1.4.0 tag points at the wrong commit, fix the release | release | git |
| 16 | my pods aren't rolling after I push a new image to ghcr | flux-ghcr-image-automation | flux |
| 17 | set up ImagePolicy and ImageUpdateAutomation for the api deployment | flux-ghcr-image-automation | flux |
| 18 | write a GitHub Actions workflow that publishes images Flux can pick up | flux-ghcr-image-automation | flux |
| 19 | scaffold a FluxCD repo layout with dev and prod overlays | flux-kustomize-layout | flux |
| 20 | review my kustomization base/overlay structure | flux-kustomize-layout | flux |
| 21 | how should I structure my GitOps repo for three clusters | flux-kustomize-layout | flux |
| 22 | write a README for this CLI | technical-writing + asd-ste100 | writing |
| 23 | review this RFC draft for clarity | technical-writing + asd-ste100 | writing |
| 24 | write the PR description for this change | technical-writing + asd-ste100 | writing |
| 25 | rewrite this in simplified technical english | asd-ste100 | writing |
| 26 | explain this simply, in STE | asd-ste100 | writing |
| 27 | write the AGENTS.md section for the new lint step | asd-ste100 + technical-writing | writing |
| N1 | add a dark mode toggle to the settings page | frontend-design + coding | frontend-design + coding (unchanged) |
| N2 | write a Dockerfile for the go service | dockerfile + coding | dockerfile + coding (no flux) |
| N3 | add a make target that takes a namespace variable | makefile + coding | makefile + coding (no git/flux) |
| N4 | before you start, restate what I'm asking | in-your-own-words | in-your-own-words (no writing) |

Result: every prompt routes to its parent at least as reliably as before. No
merge was reverted. Listing cost fell from about 2,792 to 2,406 tokens
(characters / 4) and from 29 to 23 skills.

Limit: one routing pass, by an agent rather than each harness. Re-run in a
real harness if a merged skill seems to miss.
