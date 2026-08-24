# Coding Skill Architecture

## Question

Should `coding` become many standalone skills rather than a small entrypoint with conditional reference files?

## Evidence

Agent Skills use progressive disclosure. Hosts keep each skill's name and description available for selection, load the selected `SKILL.md` body, then load linked resources only when needed. The open [Agent Skills specification](https://agentskills.io/specification#progressive-disclosure) puts those stages at roughly metadata, the skill body, and optional resources. It recommends focused, on-demand reference files and says a full `SKILL.md` loads when its skill activates.

[OpenAI's skill authoring guidance](https://github.com/openai/skills/blob/main/skills/.system/skill-creator/SKILL.md) recommends a lean overview and navigation file, with conditional details in one-hop references. [Anthropic's guidance](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills) likewise describes progressive disclosure as loading only the instructions and resources needed for the task. Neither source prohibits routers. Neither gives evidence that a large catalog is cheaper because of prompt caching.

Pstack's `guard-the-context-window` principle says to avoid raw material that will not be used, route large payloads to subagents, and keep frequently used content inline. It does **not** call routers an anti-pattern. In fact, [Pstack's `poteto-mode`](https://github.com/cursor/plugins/blob/main/pstack/skills/poteto-mode/SKILL.md) is a router: it selects from playbooks and invokes other skills as steps need them. The same repository's [`cursor-sdk` skill](https://github.com/cursor/plugins/blob/main/cursor-sdk/skills/cursor-sdk/SKILL.md) explicitly keeps a short main skill and opens focused reference files only for the matching SDK question. The apparent contradiction is therefore a category error: a router is not a context problem; indiscriminate loading is.

## Current State

`skills/coding/SKILL.md` is a 441-word navigation and baseline file. Its references contain roughly 7,500 words, but a normal task reads only the relevant ones. It already has the right broad shape for progressive disclosure: implementation, refactoring, design, review, debugging, verification, file organization, Go layout, and assertions are conditional concerns.

The main issue to watch is not its existence. It is whether the entrypoint stays small and whether the references have overlapping or contradictory rules. The universal rules belong in `coding`; specialized process belongs behind a specific route.

## Recommendation

Do **not** perform a large split of `coding` now. Keep it as the ambient software-work entrypoint and trim it over time to:

- Universal constraints and the user's durable engineering preferences.
- A precise routing table.
- Links that are one level deep.
- No duplicated specialist workflow.

Keep the current reference files as references. Splitting each into a public skill would add always-present metadata, create overlap in activation phrases, and make ordinary work depend on multiple skills firing correctly. It would not by itself reduce runtime context.

Create a standalone skill only when it has all or most of these properties:

- A user can reasonably ask for it by name before implementation starts.
- It has a distinct input, output, tool sequence, or verification standard.
- It should work without applying the full ordinary coding workflow.
- It can own its rules without repeating `coding`.

That supports future standalone skills such as measured performance work, codebase archaeology, or a project-local real-app verification workflow. A separate code-review skill is plausible only if review becomes a frequent, report-only operation with a richer independent contract. It is not justified just because `review.md` exists.

## Next Check

Treat this as an evaluation question rather than an architectural belief. Run representative implementation, refactor, review, debugging, and design prompts against the current router and any proposed split. Measure correct activation, which material was actually read, outcome quality, and missed or conflicting instructions. Change the structure only if the split improves those results.
