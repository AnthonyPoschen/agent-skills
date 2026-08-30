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

## Outcome Proof Is Required

Before declaring a change complete, state its expected observable outcome and
verify that outcome directly. The right proof depends on the artifact the
request actually concerns:

- When the changed source file is itself the delivered artifact, inspect the
  resulting file or diff.
- For a deterministic transformation, run it with representative input and
  inspect the returned or rendered value.
- For dynamic behavior, run the service, command, application, or UI and inspect
  the result through its normal path.
- For persistence or side effects, inspect the actual row, file, message, or
  other state the operation produced.

A formatter, linter, build, or passing test supports confidence but is not
outcome proof unless it observes the requested behavior. Report the direct
proof separately from supporting checks and any path that could not run.

## Automated Tests Protect Stable Contracts

Add an automated test only when it protects a stable, meaningful contract that
must remain true as the implementation evolves: a product rule, public API,
data or serialization guarantee, security property, or correctness invariant.
The test needs a trustworthy oracle and should survive reasonable refactors.

Do not add a test merely because code changed, a bug occurred, or the behavior
may change again. Do not pin incidental diagnostic text, private call sequences,
temporary data shape, or other implementation details. For example, assert an
exact log message only when it is a deliberately stable operator or machine
interface, not when it is ordinary diagnostic output.

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

- A changed pure transformation benefits from a unit test when its result is a
  stable contract worth protecting; otherwise, directly running it may be
  sufficient outcome proof.
- A changed query, migration, repository, or serialization mapping needs a real
  integration check when available.
- A changed user workflow needs a check that drives the workflow, not internal
  setters or a mock-only equivalent.
- A small, low-risk internal change may need existing checks plus a direct local
  exercise rather than new test scaffolding.

Tests should encode stable behavior and contracts, not mirror implementation
details. Avoid tests that mostly validate mocks, depend on fragile timing, or
require large unrelated fixture setup for little signal.

## Do Not Test Source Presence

Do not write tests that read source files and assert that a particular string,
selector, function name, route literal, CSS declaration, or block of markup
exists. Those checks merely prove that the implementation presently resembles
the implementation just written; they make routine refactors and copy edits
expensive without proving a user-visible outcome.

This includes regex assertions over HTML, CSS, JavaScript, Go, or configuration
source. Delete existing tests of that form unless they validate a generated
artifact or another deliberately stable text-based public contract. Replace
them with an executable test at a public seam, such as a rendered UI flow, an
HTTP request, a command invocation, or a deterministic exported operation.
