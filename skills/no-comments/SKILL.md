---
name: no-comments
description: "Use when auditing a diff or files for comments. Delete comments by default, keep only ones about constraints we cannot change, and encode claimed constraints as types, runtime checks, tests or lint rules."
---

# No Comments

Adapted from poteto's pstack `no-comments` skill (https://github.com/poteto).

Treat every comment as guilty unless it describes something we cannot change.

## Scope

Use the caller's files or diff. Otherwise use the current diff against the base branch, default `main`, including the working tree.

## Steps

1. Review each comment in scope with fresh eyes. Delete narration, restated names, stale notes and contracts the code can express. A keep survives only with proof it is about something we cannot change (an external API quirk, a spec), and it links the source.
2. Audit scoped lint and type suppressions. Correctness or safety suppressions must be fixed, not kept.
3. Fix trivial findings directly by deleting a dead path, dropping a parameter, or using the real API. If a fix needs a new shape, sketch it first, then implement.
4. Implement the smallest root-cause fix in scope. Remove every named workaround. If the root cause is out of scope, land the smallest in-scope fix and report the rest open. Never bolt on symptom guards.
5. Constraint comments say `do not remove`, `do not change wording`, or `talk to X before changing`. Offer the cheapest in-scope type, runtime check, test, or CI lint (see the `encode-lessons` skill). Wait for approval unless the caller pre-approved. If approved, encode then delete. Otherwise delete and report the constraint open.
6. Report the deletion count, kept comments with their proof, fixes, encodings, unenforced constraints, and other open work.
