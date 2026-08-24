# Verification Strategy

Use this reference when choosing how to verify a code change, especially one
that touches persistence, a framework, a service, a CLI, or a user interface.

## Choose Evidence From The Behavior

Ask first: what is the closest practical proof that this behavior works in the
real system? Choose the check that observes that behavior, not the easiest
substitute to construct.

- Use unit tests for pure, deterministic logic and narrow edge cases.
- Use a real local dependency for behavior that depends on it. For database
  work, prefer an isolated local database using the real engine, real schema or
  migrations, application path, and an assertion on the stored or retrieved
  state.
- Use an integration or end-to-end check for framework configuration, wiring,
  serialization, authentication boundaries, CLI commands, UI flows, and
  multi-component behavior when that path can run locally.
- Use mocks, fakes, or stubs only when the real dependency is unavailable,
  unsafe, prohibitively costly, inherently nondeterministic, or when simulating
  a failure that cannot be produced safely. Keep the mock at an existing
  external boundary.

Do not replace a dependency that is practical to run locally merely because a
stub is easier. A passing mock proves the mock contract; it does not prove the
application works with the real database, framework, or service.

## Real Dependency Checks

When an integration is in scope and the project can run it locally:

1. Use isolated, disposable state: a temporary database, schema, data
   directory, port, account, or namespace owned by the check.
2. Start from the same schema, migrations, configuration, and application entry
   point used by the relevant local environment.
3. Exercise the real behavior through its normal boundary: command, handler,
   service, UI, or public API.
4. Inspect the observable result: response, database row, file, message, UI
   state, or other side effect.
5. Clean up only the state the check created.

Prefer a deterministic scripted check when it is reasonably cheap to maintain.
Do not build a large test harness solely to satisfy a testing rule; use the
closest useful executable check instead.

## Test Scope

Use the narrowest check that still provides meaningful evidence:

- A changed pure transformation usually needs a unit test.
- A changed query, migration, repository, or serialization mapping needs a real
  integration check when available.
- A changed user workflow needs a check that drives the workflow, not internal
  setters or a mock-only equivalent.
- A small, low-risk internal change may need existing checks plus a direct local
  exercise rather than new test scaffolding.

Tests should encode behavior and contracts, not mirror implementation details.
Avoid tests that mostly validate mocks, depend on fragile timing, or require
large unrelated fixture setup for little signal.
