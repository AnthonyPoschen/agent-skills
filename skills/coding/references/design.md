# Design Workflow

Use this reference for architecture, API design, module boundaries, dependency
direction, and design-only discussion.

## Design Process

- Start from the caller's job and the existing system shape, not an abstract
  ideal.
- Identify constraints: public contracts, data model, deployment model,
  framework, performance, security, compatibility, and team conventions.
- Design every meaningful boundary around consumer workflows before designing
  its internal machinery. Consumers include external users and internal callers.
- Prefer designs that reduce coupling at real boundaries: domain/persistence,
  transport/domain, UI/state/effects, external provider/application policy.
- Keep cohesive work direct when the caller already owns the decisions and can
  understand the flow in place. Do not create boundaries just to divide code.
- Name tradeoffs explicitly when more than one approach is viable.

## Caller Workflow First

- Start with the shortest realistic usage example. Make the common path read in
  the caller's language, not the implementation's language.
- Names should describe tasks the caller wants to perform: poll,
  update, attach, configure, save, load, reset, find, query, render, parse.
- If the usage example needs excessive setup, exposes internal concepts, or
  forces conversions, adapters, or sequencing hacks, redesign the boundary
  before polishing internals.
- Let a meaningful boundary's real workflows drive its internal data structures
  when necessary. Do not preserve an existing representation at the cost of an
  awkward caller experience.
- Keep common workflows ergonomic even when the internal implementation is more
  generic, while retaining direct or lower-level access where callers have a
  real need for it.

## Boundaries That Earn Their Cost

- A boundary earns its cost when it removes knowledge callers should not carry:
  an invariant, lifecycle, policy, representation conversion, or integration
  detail.
- It must make the real caller's job shorter, clearer, or safer. A single caller
  is enough when the boundary owns meaningful complexity; multiple callers are
  not enough when it merely moves the same reasoning elsewhere.
- Prefer one obvious, domain-specific path for each common task. Keep
  lower-level escape hatches only when callers have a clear need.
- Do not collapse clear operations into a generic method merely because the
  implementation can be parameterized.
- Do not add a wrapper, interface, or module that only renames provider APIs,
  passes through the same types and arguments, or exists for a hypothetical
  future implementation.

## Internal Module APIs

- Treat a module as having an API when it owns a coherent responsibility that
  callers should be able to use without understanding its internals. Do not
  manufacture an API for every small private operation.
- Design that API for its real consumers: sibling modules,
  application/gameplay code, tools, tests, or integration layers.
- Keep the surface intentionally limited where the module owns invariants, state
  transitions, platform details, persistence shape, or performance constraints.
- Hide internal machinery, but do not make callers fight the API merely to keep
  the module small.
- Prefer a small set of clear, domain-specific operations over broad internal
  state or generic plumbing. If a sibling module must know internal sequencing,
  state layout, or edge handling, the boundary is too leaky.

## Layering And Vocabulary

- Separate layers by user intent and responsibility. For example: raw device
  state vs gameplay actions, transport DTOs vs domain models, persistence
  records vs runtime services.
- Keep vocabulary consistent across related types. If several objects support
  the same concept, prefer the same verb names and return shapes.
- Do not leak lower-layer concepts into higher-layer workflows unless they are
  part of the user's mental model.
- Keep boundaries explicit: higher layers compose lower layers; lower layers do
  not know higher-level policy.
- Prefer direct composition when the caller and implementation share the same
  concepts. Create a boundary only when it reduces total reader work across the
  system, rather than adding layers for their own sake.

## Plain Data Boundaries

- Use plain data structs for persistence, serialization, configuration,
  user-editable state, snapshots, and import/export surfaces.
- Keep data-exchange types boring: stable field names, simple values, optional
  fields where omission is meaningful, and validation at the boundary that
  applies them.
- Provide helpers to convert between runtime state and plain data when that
  makes save/load/edit workflows straightforward.
- Avoid forcing consumers through a complex object graph when they need to store,
  inspect, edit, diff, or transmit data.

## Symmetry Across Related Types

- Design related objects as a family. Shared concepts should use shared names,
  semantics, and error behavior.
- Symmetry is valuable at the public API even when implementations differ.
- Review similar methods together. If one has different edge-case behavior,
  decide whether that difference is intentional and document it through naming,
  types, or tests.
- Keep implementation deduplication behind the symmetric public API, not at the
  expense of consumer clarity.

## Examples As API Tests

- Write or inspect a realistic example before considering the design finished.
- The example should show setup, common operations, error handling, and the
  resulting data flow.
- Treat awkward examples as design feedback. If the example needs comments to
  explain basic usage, the API names or layering may be wrong.
- Keep examples focused on real workflows rather than demonstrating every low
  level capability at once.

## Design Review Checklist

- Each meaningful boundary reads in its real callers' workflow terms.
- The common path is short, predictable, and uses domain vocabulary.
- Cohesive local flows remain direct rather than being fragmented into layers.
- Each abstraction removes more caller burden than it adds in concepts,
  configuration, or indirection.
- Internal module APIs are intentionally limited but ergonomic for their real
  callers.
- Related types use consistent names, return shapes, and edge-case behavior.
- Plain data boundaries exist for persistence/configuration/editing when needed.
- Lower-level implementation details do not leak into higher-level workflows.
- A realistic example validates the design ergonomics.
- Risks, migration steps, and testing strategy are clear.

## Design Output

- State the recommended approach first.
- Explain why it fits the existing codebase and consumer workflow.
- Show the expected caller usage shape when API ergonomics matter.
- Call out which internals should be shared or hidden behind that surface.
- Call out risks, migration steps, and testing strategy.
- Keep speculative future extensibility out unless there is current evidence it
  matters.
