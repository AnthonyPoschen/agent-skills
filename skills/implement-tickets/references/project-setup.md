# Project Setup And Discovery

Do not assume that a repository uses a particular issue pack, label vocabulary,
default branch, test command, provider, or implementation skill. Resolve setup
in this order:

1. explicit supervisor arguments;
2. the repository contract at `.github/implement-tickets.json`;
3. linked project documentation and read-only tool discovery;
4. direct questions to the user.

Never create labels, rewrite issues, invent dependencies, or select a plausible
test command merely to make preflight pass.

## Discoverable Documentation Chain

Projects that use this workflow repeatedly should make it discoverable from the
nearest `AGENTS.md`:

```text
AGENTS.md
  -> docs/agents/issue-execution.md
       -> .github/implement-tickets.json
       -> tracker-specific conventions or issue templates
```

`AGENTS.md` should contain a short routing statement, not a duplicate copy of
the workflow. For example:

```md
For dependency-gated issue execution, see
[`docs/agents/issue-execution.md`](docs/agents/issue-execution.md).
```

The linked document should explain the meaning of the ready label, where native
dependencies are recorded, who may merge, the target branch, verification
expectations, cross-repository responsibilities, and any implementation skill
invocation. It should link the machine-readable contract rather than allowing
the two sources to drift silently.

## Project Contract

Copy [`assets/implement-tickets.json`](../assets/implement-tickets.json) to
`.github/implement-tickets.json` and adapt it. Supported fields are:

- `version`: contract version; currently `1`;
- `source`: `github` or `gitlab`; `jira` is recognized as a fail-closed beta;
- `target` and `remote`: integration ref;
- `ready_label`: the label that explicitly authorizes scheduling;
- `concurrency`: bounded active-worker limit;
- `verification`: independent commands run by the supervisor after a worker;
- `implementation_invocation`: optional project implementation skill prefix;
- `harness`, `worker_agent`, `worker_model`, and `manager_model`: optional
  execution settings when the project deliberately owns them.

Prefer leaving provider/model fields out of a shared project contract unless
the project truly requires them. User or machine-level selection normally owns
those values.

Unknown fields and unsupported contract versions fail closed. Command-line
arguments override contract fields for a specific run.

## Required Preflight Questions

When documentation, the contract, and read-only inspection do not resolve a
setting, stop before assignment, worktree creation, or worker launch. Ask all
known questions together, including:

- Which branch is the integration target?
- Which existing label makes an issue ready, and which issues should carry it?
- Where are explicit blockers recorded?
- Which worker harness and model should run?
- Which independent commands verify a completed worker checkpoint?
- Does the project use an implementation skill invocation?
- Are pushes and draft reviews authorized, and who owns merges?

Record the answers in project documentation and the contract when they are
stable project facts. Keep one-run choices in supervisor state instead.

## Run Selection

Select work for one run without changing the project contract:

- repeat `--issue <number>` for a fixed set;
- repeat `--issue-range <start>-<end>` for inclusive numeric ranges;
- use `--issue-query <GitHub search>` for a dynamic category;
- omit all three to select every open issue carrying the ready label.

An issue query cannot be combined with fixed numbers or ranges. Initialization
resolves every selection mode to a fixed list of open, ready-labeled issue
numbers and persists it in run state. Selection does not bypass assignment or
dependency gates, and a restart does not absorb unrelated new work.
