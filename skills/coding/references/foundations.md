# Foundation First

Use this reference when a change introduces a consequential data model,
persistent representation, shared scaffold, or concurrent state. Do not apply
it to a small local helper that can stay direct.

## Establish The Shape Before The Logic

Decide the data shape and its dominant creation, read, update, and ownership
paths before distributing logic through callers. A clear representation makes
the ordinary code obvious. A poor one spreads conversions, conditions, and
coordination across the system.

- Start from real caller workflows and access patterns, not a generic schema.
- Make related types and data models converge where they represent one concept.
- Prefer explicit data and operations over a clever abstraction that hides the
  actual model.
- Do not turn three similar lines into a general mechanism unless the mechanism
  removes a real invariant, repeated decision, or invalid state.

## Build Foundations Only When They Pay Now

Put a shared type, test support, CI check, or other scaffold first when it
materially simplifies several later parts of the same bounded change. Do not
build a foundation for a hypothetical future.

Ask whether each later phase becomes simpler, safer, or more direct because
the foundation exists. If not, keep the current change direct and add support
only when evidence requires it.

Before adding a foundation, remove dead paths, pass-throughs, misleading state,
and redundant supporting machinery in the affected area. Then ask whether the
smaller design still needs a shared type, scaffold, validator, or configuration
point. Subtraction makes the real foundation easier to see and prevents a
temporary workaround from becoming permanent structure.

## Give Concurrent State An Owner

Before two actors share mutable state, ask what happens if either changes it
while the other is working. Prefer independent ownership when they do not need
one canonical state. When state must be shared, make the owner and coordination
rule explicit rather than relying on timing or caller convention.

## Grow One Coherent Model

Let each increment establish or deepen one coherent model. Do not implement a
new capability by scattering special-case coordination across callers when a
bounded model or owning operation would make the whole path simpler.

Keep the model proportional to the task. A direct local flow remains the right
foundation when no shared data shape, invariant, lifecycle, or concurrent
ownership needs to exist.
