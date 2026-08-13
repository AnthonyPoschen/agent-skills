# Harness Capability Contract

Treat the Go supervisor as a composable CLI, not as a mandatory replacement for
the host harness. Detect what the harness already guarantees, then retain the
smallest supervisor surface that supplies missing durable behavior.

| Capability | Capable interactive harness (for example Codex) | Headless harness |
| --- | --- | --- |
| Worker fan-out and sandboxing | Prefer native agents and sandbox controls | Use supervisor launch and policy adapters |
| Reasoning and feedback classification | Prefer manager judgment after deterministic normalization | Use conservative supervisor classification |
| User-visible progress | Keep an active relay turn; use commentary headed by `Events` | Stream `events --follow` to the host UI |
| Tracker polling and deduplication | Keep supervisor polling or implement the same durable cursor contract | Keep supervisor polling |
| Dependencies and integration gates | Keep deterministic supervisor state | Keep deterministic supervisor state |
| Recovery after session/process loss | Keep supervisor state, logs, and ownership tokens | Keep the full supervisor loop |
| Commit, push, and draft review publication | Keep centralized manager/supervisor authority | Keep centralized supervisor authority |

A harness may call `init`, `sync`, `status`, `events`, `inbox`, `dispatch`, or
`publish` independently. It does not need to run `supervise` when its own event
loop already provides scheduling and process ownership.

The provider-neutral boundary is normalized items, blockers, reviews, feedback,
checks, and events. Harness-native agent IDs and tool events may remain native;
store only the correlation IDs needed for recovery.

Important supervisor events are checkpoint commits, branch pushes, review
publication, new feedback, failed checks, requested input, integration, worker
loss/failure, and a queue that becomes blocked or idle. Interactive managers
surface these immediately. Routine unchanged polls stay silent. Emit a compact
heartbeat only after five minutes without an important event.

Background execution and chat delivery are separate capabilities. A detached
supervisor can keep implementing, but it cannot inject output into a Codex turn
that has already ended. Codex therefore keeps a bounded wait active until the
run reaches a terminal state and replays durable sequence-numbered events after
any reconnect. Only events confirmed as shown advance the chat relay cursor.
