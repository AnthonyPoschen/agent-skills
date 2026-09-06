---
name: grill-with-docs
description: >
  Run a relentless interview to sharpen a plan or design while creating and
  updating domain docs (CONTEXT glossary, ADRs) as decisions crystallise. Use
  when starting or refining a feature/design and the project should leave
  durable terminology and decisions behind — not only a chat transcript.
disable-model-invocation: true
---

# Grill With Docs

Run a grilling interview that sharpens the plan or design, and keep the project's
domain docs current as you go.

## Process

1. Load and follow the first-party [`domain-modeling`](../domain-modeling/SKILL.md)
   skill in this repository (not a third-party Matt Pocock install).
2. Interview relentlessly: challenge fuzzy terms, surface contradictions with
   code and existing glossary entries, force concrete scenarios, and refuse to
   leave ambiguous decisions hanging.
3. As terms and decisions crystallise, update docs *inline* via domain-modeling:
   - `CONTEXT.md` (glossary only)
   - `docs/adr/*` when the ADR criteria are met
   - Optional short living alignment files **only if** the repo declares them
     (see domain-modeling "Living short alignment docs")
4. Do not depend on `/setup-matt-pocock-skills` or any mattpocock/skills path.

## Done when

- The plan/design is sharp enough to split or implement without re-litigating
  vocabulary
- Glossary and any warranted ADRs reflect what was decided in this session
- Declared living-short docs (if any) match the new facts — or were correctly
  skipped because the repo does not declare them

## Source

First-party adaptation of the grill-with-docs skill previously used via
mattpocock/skills. Owned here so grill → docs stays iterable without a
third-party runtime dependency.
