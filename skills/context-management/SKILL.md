---
name: context-management
description: >
  Use when a task involves large logs, long files, broad search results,
  screenshots, repeated reads, or fan-out planning. Keep the working context
  focused through targeted extraction, compact handoffs, scoped phases, and
  bounded delegation.
---

# Context Management

Protect the working context for the decisions that remain. Do not spend it on
raw material that will not affect the next action.

## Read With A Question

Before opening a large file, log, diff, document, image set, or search result,
state the question it must answer. Then read only the relevant section, fields,
or examples.

- Prefer a targeted search, range, filter, or screenshot page over reading the
  whole source.
- Do not reread material already understood. Keep a short note with the fact,
  evidence location, uncertainty, and next action.
- Read more only when the current evidence cannot support the decision. Do not
  collect background merely because it is available.

## Keep A Compact Working Record

When a task has several phases, preserve only what the next phase needs:

- the goal and constraints;
- decisions made and why;
- direct evidence and its location;
- open questions, risks, and the next action.

Summarise verbose tool output in those terms. Do not paste a long transcript
into later work unless its exact wording is necessary.

## Delegate Only When It Helps

Delegate a large investigation when it has a clear question, can be explored
independently, and the main task can proceed from a concise evidence-backed
summary.

- Give the delegate a bounded scope and the facts it must return.
- Ask for source locations and uncertainty, not a retelling of every file read.
- Inspect important evidence yourself before making a consequential decision.
- Keep small, tightly coupled reads in the main thread. Do not create agents
  merely to avoid reading a few relevant files.

## Load Instructions Progressively

Keep rules used on nearly every invocation in the active skill. Keep variants,
language details, and rare procedures in directly linked references. Load one
reference when the task needs it; do not load a catalog speculatively.

## Reassess Before Context Becomes A Problem

Pause before a broad new phase, a long fan-out, or a large raw payload. Restate
the current decision and choose the smallest next read or delegated task that
can resolve it. If the context is already crowded, make a compact handoff
record before continuing.

## Sources

- [Pstack guard the context window](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-guard-the-context-window/SKILL.md)
- [OpenAI skill authoring guidance](https://github.com/openai/skills/blob/main/skills/.system/skill-creator/SKILL.md)
