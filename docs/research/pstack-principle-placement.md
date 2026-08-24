# Pstack principle placement

## Scope

This note compares three Pstack principles with the current skills library. It
recommends the smallest durable placement for each idea. It does not propose
copying Pstack's user-invoked skill structure: all three source skills set
`disable-model-invocation: true`.

## Summary

| Pstack principle | Current coverage | Recommendation |
| --- | --- | --- |
| Redesign from first principles | Partial | Update the coding implementation and refactoring references. |
| Prove it works | Strong | No change. |
| Never block on the human | Partial and scattered | Add one standalone execution-autonomy skill. |

## Redesign from first principles

Pstack says that a new requirement should be treated as though it had been
known on day one, rather than bolted onto the existing design. Its concrete
method is to understand the affected design, ask what a clean-sheet design
would be, propagate the result through types, docs, examples, and rationale,
then deliver the redesign incrementally. Its stated goal is to preserve option
value. [Source](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-redesign-from-first-principles/SKILL.md)

The local coding skill already covers much of the practical boundary case.
`skills/coding/references/implementation.md` says to inspect direct consumers,
diagnose whether boundary friction is systemic, redesign and migrate a bounded
set when that makes those consumers simpler, and avoid special-case escape
hatches. `skills/coding/references/refactoring.md` gives the same direction for
several consumers. `skills/coding/references/design.md` already begins with
the caller's job and rejects awkward consumer conversions.

The gap is a deliberate clean-sheet comparison before accepting the existing
shape as a constraint. Add a short conditional rule to the implementation and
refactoring references:

> When a meaningful requirement exposes friction in a touched design, compare
> the direct change with the smallest shape that would make sense if the
> requirement had existed from the beginning. If the latter is clearly simpler
> for the affected callers and fits the task's natural scope, implement and
> verify that coherent redesign. Update every directly affected type, example,
> document, and rationale. Otherwise keep the change direct and report the
> broader redesign opportunity.

This should not be a standalone skill. It is a design decision within an
implementation or refactor, and a separate skill would create a second routing
path for rules that already belong beside the scope and consumer-migration
rules. The wording must retain the existing scope guard: “from first
principles” does not authorize a repository-wide rewrite or a premature
framework.

## Prove it works

Pstack requires direct observation of the real artifact before declaring work
done. It rejects compilation, cached state, self-reports, and other proxies;
for code it calls for building, exercising the actual feature path, checking
input through output, and using the complete integration path. It also says to
inspect delegated artifacts instead of accepting a delegate's summary, and to
prefer a repeatable script where that is cheap. [Source](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-prove-it-works/SKILL.md)

This library already covers the principle more precisely in
`skills/coding/references/verification.md`. It requires direct outcome proof,
real local dependencies where practical, inspection of the produced state, and
a deterministic script where reasonably cheap. It explicitly says that a
formatter, linter, build, or passing test is supporting evidence unless it
observes the requested behavior. `implementation.md`, `review.md`, and
`debugging.md` repeat the rule in their own contexts. `context-management`
also requires the main agent to inspect important evidence before a
consequential decision.

No change is recommended. A new standalone skill would repeat the exact
ground-truth rule the coding router already loads for verification decisions.
The existing wording is also better aligned with the library's selective-test
policy: it distinguishes direct proof from durable automated contract tests,
rather than turning “script the check” into a blanket testing requirement.

## Never block on the human

Pstack says to proceed with reversible work, present the result, and ask only
when intent is genuinely unknowable. It treats the human as an asynchronous
reviewer and separates product direction from execution. It explicitly keeps
confirmation for irreversible actions such as force-pushes, production data
deletion, and external messages. [Source](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-never-block-on-the-human/SKILL.md)

There is related coverage, but no library-wide rule. The coding skill encourages
scoped action, and `implement-tickets` avoids waiting on routine unblocked work,
but neither tells a normal task agent when an execution choice can be made
without a clarification. The release and ticket skills contain local stop
conditions, but those are not a general policy.

Add a small standalone skill named `act-within-scope` (or another plainer name
chosen during implementation). It should apply when an agent has a requested
outcome and is deciding whether to ask for permission or proceed. Its rule set
should be short:

1. Infer a reasonable execution choice from the request, repository evidence,
   and established conventions; then proceed with reversible, in-scope work.
2. Present the decision, result, and verification so the user can correct the
   course after the fact.
3. Ask only for a genuine product-direction ambiguity, missing authority, or a
   decision that materially changes scope, cost, or external impact.
4. Always stop for destructive, irreversible, security-sensitive, or external
   actions unless the user clearly authorized them.
5. Record a discovered repeatable issue in the smallest appropriate guard or
   follow-up, rather than waiting for a new prompt to notice it again.

Keep it standalone because it applies across coding, research, documentation,
and operations. Do not call it from other skills. Its description can activate
it for autonomy or clarification decisions, while its body preserves the
existing safety and product-direction boundaries. Do not use Pstack's absolute
name: “never block” is rhetorically useful but wrong when authority or intent is
actually missing.

## Sources

- [Pstack: Redesign From First Principles](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-redesign-from-first-principles/SKILL.md)
- [Pstack: Prove It Works](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-prove-it-works/SKILL.md)
- [Pstack: Never Block On The Human](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-never-block-on-the-human/SKILL.md)
- `skills/coding/references/implementation.md`
- `skills/coding/references/refactoring.md`
- `skills/coding/references/design.md`
- `skills/coding/references/verification.md`
- `skills/coding/references/review.md`
- `skills/coding/references/debugging.md`
- `skills/context-management/SKILL.md`
- `skills/implement-tickets/SKILL.md`
