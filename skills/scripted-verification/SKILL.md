---
name: scripted-verification
description: >
  Use when you are about to read logs, screenshots, diffs, timings, or
  benchmark numbers to decide whether something passed; when writing a
  "check", "verify", "retest", or "looks right" brief for another agent; when
  a manual check sits in AGENTS.md, a PR template, or a review comment; when
  about to sleep, wait, or poll for a result; or when re-checking a commit or
  re-running an A/B that was already judged. Turn the check into a
  deterministic script that prints PASS, FAIL, or UNSURE with exit codes and
  keeps the full evidence on disk for review.
---

# Scripted Verification

Agents are expensive judges. Every step re-sends the whole context, so an agent
that re-reads results in prose pays for that context on every step. One audit
found the largest cost was agents re-judging the same kind of result: about 29
runs in a day, some for 25 to 40 minutes at high effort, at 110K to 134K tokens
per step. A script makes that judgement once, for free, every time.

## When To Reach For This

Stop and script the check when you notice yourself:

- about to read logs, screenshots, diffs, or numbers to decide pass or fail,
  especially the same comparison against the same baseline as last time;
- writing a "check X", "verify X", "retest", or "make it look right" brief
  for another agent;
- following a manual check from AGENTS.md, a PR template, or a review comment;
- about to sleep, wait, or poll until a run, build, or resource finishes;
- re-checking a commit, or a commit whose inputs did not change, that already
  has a verdict;
- launching a costly check that fails for a cheap reason: merge conflicts, a
  busy shared resource, or missing arguments;
- arguing whether a borderline number is noise. Measure the noise floor
  instead.

## Script Any Check That Runs More Than Once

The second time an agent judges the same kind of result, write a script.

- Print one verdict line: `PASS`, `FAIL`, or `UNSURE`, with the evidence
  path. Use distinct exit codes: `0` pass, `1` fail, `2` unsure, `3` the
  check itself broke.
- Write the full evidence (screenshots, diffs, raw numbers, logs) to `--out` or
  a debug directory. Do not print it.
- Use `UNSURE` for anything that is not a clean pass or fail: the borderline
  band and states the script cannot decide. Do not force a guess.
- Agents read the verdict line. On `UNSURE`, an agent reviews the saved
  evidence first. Involve a person only when that review still cannot decide.
- Every verdict stays reviewable: the evidence must show why the script decided
  what it did.

Skeleton: [`scripts/verdict-skeleton.sh`](scripts/verdict-skeleton.sh).

## Set Thresholds From Measured Noise

- Measure the noise floor first: run the check several times on an unchanged
  state and record the spread. Set the threshold above it, with an
  `UNSURE` band between noise and a clear failure.
- Do not lock references or goldens until the state is stable. A golden taken
  during churn turns every later run into noise.

## Pre-Check Before Launching An Agent

Run cheap checks first and skip the agent when one fails:

- the branch has merge conflicts;
- a shared resource (GPU, machine, lock) is busy;
- the commit was already judged. Cache verdicts by commit and inputs, and
  reuse the cached verdict.

## Watch Events Instead Of Polling Agents

Wait with CI, webhooks, file watchers, or a cheap script loop. Never keep an
agent idling in a wait loop; wake one only when there is a result to act on.

## Earn Trust

- A new check starts on probation: review its verdicts against the evidence.
- It graduates to trusted after a run of correct verdicts.
- Any miss, false pass, or false fail becomes a new test case for the check.

## Brief Agents Tightly

- Give short, focused tasks. Start a fresh agent per task instead of growing
  one long context.
- Use targeted reads: the failing section or the verdict, not whole files or
  logs.
- Write a precise brief: the goal, the inputs, the done criteria, and the
  exact command to run.

## Related

- `context-management` for targeted reads and handoffs.
- `create-verification-skill` for driving an app and collecting evidence.
- `lint-rules` and `encode-lessons` for turning a repeated judgement into a
  guard.
- `test-audit` for whether a check should exist at all.
