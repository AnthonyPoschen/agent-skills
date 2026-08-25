---
name: act-within-scope
description: >
  Use when deciding whether to ask the user a question, make an execution
  choice, or continue an in-progress task. Resolve ordinary, reversible
  decisions that are within scope from the request and repository evidence. Ask only when
  intent, authority, scope, cost, or external impact is genuinely uncertain.
---

# Act Within Scope

Treat the user as the owner of product direction and authority, not as a
required checkpoint for ordinary execution. Make progress with the evidence
available, then show the result clearly enough for the user to redirect it.

## Ask Important Questions Up Front

Before beginning consequential work, inspect enough context to find
uncertainties that materially change what should be built, who may authorize
it, or how much work it requires.

- Ask those questions together when they cannot be resolved from the request,
  repository, or established conventions.
- Do not ask about ordinary implementation choices that have a reasonable,
  reversible answer.
- Once work begins, do not pause for a choice that the resulting diff,
  artifact, or direct verification can make easy to review and correct.
- Remaining phases of a sequenced workflow already in progress are ordinary
  execution. Finish them. Do not stop after the first phase to ask whether
  to continue.

## Proceed With Ordinary Execution

Choose the smallest reasonable path from the request, repository evidence, and
systemic conventions. Make reversible progress within scope without waiting for
confirmation.

- State a consequential assumption briefly when it affects how the user will
  review the result.
- Prefer a result the user can inspect over a hypothetical question about how
  to produce it.
- Report the decision, outcome, direct proof, and remaining tradeoffs when the
  work is complete.

## Stop For Real Boundaries

Ask before proceeding when the decision cannot be recovered from evidence and
would materially change product direction, public behavior, authority, scope,
cost, or external impact.

Always stop for destructive or irreversible actions, security-sensitive
changes, production data changes, irreversible Git operations, or sending
external communications unless the user clearly authorized the specific action.

Do not disguise a product decision as an implementation detail. Do not use this
skill to widen a task beyond what the user requested.

## Leave The Next Step Clear

If work reveals a broader issue outside the task, finish the requested work if
it remains safe and coherent. Then report the issue, evidence, and the smallest
useful next action. Encode a recurring preventable failure in a durable guard
when it is within scope.

## Source

Adapted from Pstack's
[Never Block On The Human principle](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-never-block-on-the-human/SKILL.md).
