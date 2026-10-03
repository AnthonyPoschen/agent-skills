# Feature map contract

The map is the maintained source for verifying user-facing behavior. Keep implementation detail out of it. Name user paths, stable handles, required state, commands, and observable proof.

## Index

The index is `features/README.md` inside the verification skill. It has:

- **Baseline preconditions.** How this run starts, including disposable state and the doctor check. Drive only an instance this run started, unless the skill's Isolate rule says the app has one shared instance.
- **Driving conventions.** Where each recipe starts, which handles to prefer, and how mutation state is restored.
- **Proof and skip reporting.** The standards below, plus where artifacts go.
- **Features.** One link per feature file and the behavior that file covers.

## Feature file

Each feature file starts with an H1 and one paragraph of user-visible behavior. Then these H2 sections, in order:

1. `Sub-features` lists a short id and one line for each behavior.
2. `How to get to it (user POV)` lists every user entry point.
3. `Driving it with <harness>` starts with `Preconditions:` and then labeled steps. Each step pairs the user action, the exact command, and the observable result.
4. `Gotchas` lists traps that waste or invalidate a run.

## Proof standards

- Exercise the real user path. An internal setter or a test-only endpoint is not that proof.
- Capture the action and the resulting state. The final screen alone is not enough.
- Verify the side effect from a second view: the reopened record, the stored row, the written file, or the sent message.
- Mock only at a boundary the production app already isolates.
- When the safe path is a dry-run or test mode, observe what it skipped (files, network, git refs). The name of the mode is not the proof.
- Record the feature id and the entry point with every artifact.
- An unreachable path reports the command you tried and the unmet precondition.
- A skipped entry point stays unverified. Another path does not cover it.

## Cleanup

Cleanup removes instances and scratch state this run created. On a shared instance, clean the residue and leave the instance. Proof artifacts stay at the path the skill names.
