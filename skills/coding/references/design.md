# Design Workflow

Use this reference for architecture, API design, module boundaries, dependency
direction, and design-only discussion.

## Design Process

- Start from the existing system shape, not an abstract ideal.
- Identify constraints: public contracts, data model, deployment model,
  framework, performance, security, compatibility, and team conventions.
- Design the public surface around consumer workflows before designing internal
  machinery.
- Prefer designs that reduce coupling at real boundaries: domain/persistence,
  transport/domain, UI/state/effects, external provider/application policy.
- Avoid abstractions that only rename provider APIs or hide concrete technology
  without adding policy.
- Name tradeoffs explicitly when more than one approach is viable.

## Consumer Workflow First

- Start with the shortest realistic usage example. Make the common path read in
  the user's language, not the implementation's language.
- Public names should describe tasks the consumer wants to perform: poll,
  update, attach, configure, save, load, reset, find, query, render, parse.
- If the usage example needs excessive setup, exposes internal concepts, or
  hides the user's intent, redesign the public API before polishing internals.
- Keep common workflows ergonomic even when the internal implementation is more
  generic.

## Public Surface Vs Internal Machinery

- Keep small public convenience methods when they form a clear, stable,
  consumer-facing vocabulary.
- Do not collapse good public names into one generic method just because the
  implementation can be parameterized.
- Move repeated implementation mechanics behind the public API. Public methods
  may delegate to shared helpers when their bodies differ only by small
  parameters, states, fields, or enum cases.
- Let the public surface stay domain-specific while internal helpers carry
  generic traversal, state-transition, validation, or dispatch logic.
- Prefer one obvious public path for each common task. Add lower-level escape
  hatches only when consumers have a clear need.

## Module-Local Public APIs

- Treat each module as having a public surface, even when the module is internal
  to a larger application, engine, service, or tool.
- Design that surface for the module's real consumers: sibling modules,
  application/gameplay code, tools, tests, or integration layers.
- Keep the surface area intentionally limited where the module owns invariants,
  state transitions, platform details, persistence shape, or performance
  constraints.
- Keep the surface simple and ergonomic where other modules need frequent,
  repeated use.
- Hide internal machinery behind the module boundary, but do not make consumers
  fight the API just to keep the module small.
- Prefer a small set of clear, domain-specific operations over exposing broad
  internal state or generic plumbing.
- Review module APIs from the caller's point of view. If a sibling module must
  know too much about internal sequencing, state layout, or edge handling, the
  boundary is too leaky.

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
- Prefer direct composition over broad abstraction when the boundary does not
  need independent policy.

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

- The public API reads in consumer workflow terms.
- The common path is short, predictable, and uses domain vocabulary.
- Public convenience methods are preserved when they make the API clearer.
- Internal module surfaces are intentionally limited but still ergonomic for
  their real callers.
- Repeated internals are centralized without flattening the public surface.
- Related types use consistent names, return shapes, and edge-case behavior.
- Plain data boundaries exist for persistence/configuration/editing when needed.
- Lower-level implementation details do not leak into higher-level workflows.
- A realistic example validates the design ergonomics.
- Risks, migration steps, and testing strategy are clear.

## Design Output

- State the recommended approach first.
- Explain why it fits the existing codebase and consumer workflow.
- Show the expected public usage shape when API ergonomics matter.
- Call out which internals should be shared or hidden behind that surface.
- Call out risks, migration steps, and testing strategy.
- Keep speculative future extensibility out unless there is current evidence it
  matters.
