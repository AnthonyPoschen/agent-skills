---
name: experience-first
description: >
  Use when product, UX, feature scope, defaults, interfaces, workflows, or
  consumer tradeoffs are involved. Prioritize the real user's experience over
  implementation convenience, and keep every addition focused on the core
  workflow.
---

# Experience First

The product is the experience. Technical choices are good when they make the
real consumer's work clearer, faster, safer, or more satisfying.

The consumer may be an end user, library caller, operator, teammate, or the
next maintainer. Judge the change from their seat, not from the easiest
implementation path.

## Find The Core Workflow

Name who uses the change and what they are trying to accomplish. Follow their
path from starting point to successful outcome, including the feedback that
tells them what happened.

- Keep the common path obvious and short.
- Treat empty, loading, error, permission, and recovery states as part of the
  experience.
- Prefer real use or a realistic usage example over assumptions about what is
  convenient.

## Make Additions Earn Their Place

Every feature, control, option, setting, abstraction, and configuration point
must help the consumer complete the core workflow. Remove, defer, or simplify
anything that does not.

- Prefer a few complete, polished workflows over many partial ones.
- Do not expose internal machinery, generic flexibility, or implementation
  shortcuts when the consumer needs one clear action.
- Choose sensible defaults when they remove a decision without hiding a
  meaningful tradeoff.
- Keep useful lower-level control where consumers have a real, demonstrated
  need for it.

## Decide Through The Consumer's Experience

When implementation convenience conflicts with a clearer or more reliable
experience, favor the experience within the bounds of correctness, safety,
accessibility, and explicit requirements.

For a library or internal API, a caller needing conversions, sequencing
knowledge, or workarounds is experience debt. For a product interface, a user
needing to guess, recover from a silent failure, or navigate unused controls is
experience debt. Fix the owning shape when the problem is systemic.

## Make Uncertain Choices Cheap

When an interaction, workflow, or visual direction is uncertain, explore the
experience before committing to production structure. Use the smallest artifact
that can answer the decision, then keep only the direction that improves the
core workflow.

Do not polish a weak flow into permanence. First remove friction and unused
scope, then refine details such as hierarchy, feedback, spacing, transitions,
and error handling.

## Review From The Outside

Before finishing, answer:

1. Who consumes this, and what are they trying to achieve?
2. Does the common path become simpler, more obvious, or more trustworthy?
3. What can be removed, deferred, or made a default?
4. Can the consumer tell what happened and recover when it does not work?
5. Did implementation convenience introduce experience debt?

## Source

Adapted from Pstack's
[Experience First principle](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-experience-first/SKILL.md).
