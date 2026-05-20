# Assertion Strategy (TigerBeetle-Inspired)

Use this reference when deciding where and how to assert while writing or
reviewing code.

## Goal

- Prefer fail-fast invariant checks over silent drift.
- Treat assertions as executable design contracts for internal assumptions.
- Keep user-input validation and operational error handling separate from
  assertion checks.

Use this mental model:

- **Errors are for what the world does to the program.**
- **Assertions are for what the program should never do to itself.**

This matches the TigerBeetle-style split: convert correctness bugs into
fail-stop behavior instead of silent corruption.

## Boundaries vs Core

At boundaries, return errors. Inside validated core logic, assert invariants.

- Boundary examples: filesystem, network, database, OS calls, CLI/API input,
  config files, untrusted serialized data, plugins.
- Core examples: internal state-machine transitions, prevalidated indices,
  ownership/alias assumptions, conservation properties, impossible enum states.

## What To Assert

Assert conditions that must always be true if the code is correct:

- Internal invariants (state machine expectations, index bounds after prior
  checks, non-null internal references, ownership assumptions).
- Preconditions for private/internal helpers that are guaranteed by callers.
- Postconditions that protect critical transformations.
- Cross-field consistency and impossible enum/state combinations.

Do **not** use assertions for expected runtime failures:

- Invalid user input.
- Network, filesystem, database, or external-service failures.
- Authorization, rate limits, timeouts, contention, retries.

Handle expected runtime failures with normal error paths.

## Default Policy

When project-specific rules are absent, use this default split:

1. **Production assertions by default for correctness-critical invariants**
   where continuing would risk data
   corruption, security boundary violations, irreversible side effects, or
   silent correctness bugs that are worse than fail-stop behavior.
2. **Debug assertions** for expensive, high-frequency, or mostly-diagnostic
   checks where crash-on-production would create unacceptable availability risk.

If uncertain, choose one of these explicitly and document why. Do not hide the
decision behind vague “maybe debug-only” defaults.

## Three Assertion Levels

When language/runtime allows, use three levels:

1. **Error return** (`error`, `Result`, etc.) for expected operational failure.
2. **Always-on invariant assert/panic** for correctness-critical internal
   contracts.
3. **Debug-only assert** for expensive diagnostics.

This avoids the false binary of “everything is debug asserts” vs “everything
must crash production.”

## Promotion Criteria (Debug -> Production)

Promote an invariant to production assertions when most are true:

- Violation implies corrupted state, unsafe output, or incorrect irreversible
  writes.
- Recovery is ambiguous or likely to hide a deeper defect.
- The check is cheap enough for hot paths.
- The assertion message can drive direct diagnosis.

Keep as debug assertion when most are true:

- The check is expensive and mainly diagnostic.
- A violation is serious but the process has safer recovery/containment.
- The system has strong external guardrails and crash cost dominates risk.

## Writing Assertions

- Assert closest to where the invariant becomes guaranteed.
- Use one assertion per logical claim unless a grouped check improves
  diagnosis.
- Include concrete context in the message (ids, sizes, state names), without
  leaking secrets.
- Avoid side effects inside assertion expressions.
- Do not duplicate equivalent assertions at every layer; assert once at the
  strongest boundary.

## Language Guidance

- **Go:**
  - Keep errors-as-values at boundaries.
  - Use an always-on `Assert` helper (panic on violation) for critical
    invariants.
  - Use build-tagged `DebugAssert` only for expensive checks.
- **Zig:**
  - Keep `error` returns for expected operational failures.
  - Use explicit always-on invariant checks (`if (!ok) @panic(...)`) for
    critical contracts.
  - Use `std.debug.assert` for debug/release-safe diagnostics, understanding
    optimize-mode behavior in release-fast/release-small.

## Testing Strategy (TigerStyle-Aligned)

Treat assertions as part of the behavior contract. Tests should intentionally
surface invariant violations early, before production runtime paths depend on
undefined or corrupt state.

### What Changes In Test Design

- Keep boundary tests that expect error returns for invalid input and
  environment failures.
- Add core-invariant tests that exercise internal state transitions and verify
  invariant-preserving outcomes.
- Add targeted assertion-surfacing tests around mutation-heavy paths,
  edge-of-capacity scenarios, and state-machine transitions.
- Prefer tests that prove correctness properties (conservation,
  monotonicity/order, boundedness, idempotence) over only happy-path examples.

### Practical Patterns

1. **Boundary vs core pair tests**
   - Boundary test: invalid external input returns expected error.
   - Core test: equivalent validated flow preserves invariants after mutation.
2. **Property-oriented checks**
   - Verify invariants across many generated or table-driven cases.
   - Keep assertions enabled in these runs so invariant violations fail fast.
3. **Capacity and limit tests**
   - Hit fixed limits exactly, just below, and just above.
   - Confirm above-limit behavior is error-return at boundary, not silent drift.
4. **State-transition tests**
   - Cover legal transitions and at least one illegal/impossible transition
     attempt (through boundary API) that is rejected safely.
5. **Regression tests for assertion failures**
   - When an assertion bug is found, add a test that reproduces the path and
     passes after the fix.

### Build/Mode Coverage

- Run normal test suites in the default local mode.
- Add at least one CI path with assertion-rich mode enabled (for example,
  debug/release-safe profiles) so debug assertions are exercised regularly.
- For languages with stripped debug asserts in fast-release modes, keep
  always-on invariant checks for correctness-critical contracts and test those
  paths explicitly.

### Anti-Patterns

- Do not remove or weaken assertions to make tests pass.
- Do not convert invariant violations into generic error returns unless the
  condition is truly expected at a boundary.
- Do not rely only on end-to-end tests for invariant discovery; include focused
  unit/module tests near mutation logic.

## Review Checklist

- Every assertion checks an internal contract, not expected runtime failure.
- Debug vs production level is justified by risk and cost.
- Assertion messages are actionable and specific.
- Critical flows fail fast before unsafe mutation or side effects.
- Tests cover behavior that demonstrates the asserted contract.
- Boundary APIs return errors/null for expected bad input or environment
  failures; internal helpers assert impossible states.
- Tests are structured to exercise both boundary errors and core invariants,
  with assertion-enabled runs in regular validation.
