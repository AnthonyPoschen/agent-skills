# Proposal: quiet tool output (agents prefer small context)

**Status:** evaluation only — no build-out yet. Usage resets Tue ~18:30 Australia/Sydney.
**Owner:** Skills / agent-skills
**Repos:** genos-engine (Rust/Vulkan, default branch `master`), then genos
**Constraint:** no third-party output proxies (RTK rejected). Follow the trust pyramid: repo structure and defaults → CI/static checks → short skill prose → human review.

## Why

Every cloud-agent step re-sends the full history (`fleet-rules.md` § Cloud agent token habits, 15 Oct 2026). Verbose `cargo` / `go test` / `git show` output is paid for again on every later turn. Agents already sometimes redirect by hand (`make test > /tmp/...; tail`); that habit is not structural, so it regresses.

## Light evidence (19 Oct 2026)

Sources: shallow clones of `AnthonyPoschen/genos-engine` (`master`) and `AnthonyPoschen/genos`; cloud-agent transcripts under `/workspace/cloud-agent-transcripts` (schema: `role=tool`, `tool_name=run_terminal_cmd`, `tool_result.value.output`).

### Repo defaults today

| Repo | Quiet defaults? | Agent-facing targets |
|---|---|---|
| genos-engine | No `.cargo/config.toml`. `Makefile` runs bare `cargo run` / `cargo` via `scripts/lint.sh` with no `--quiet`. Semgrep uses `--quiet` already. | `make lint`, `make verify`, `make lighting` print full cargo chatter |
| genos | No quiet `go test` / stylelint defaults in `makefile`. Sparse agent targets (`e2e-dev`, `run`, …). | Agents invent ad-hoc redirects; no shared `make test-quiet` |

### Transcript sample (cargo)

Across ~45 transcripts:

- **900** `cargo …` shell calls, **~473k** chars of tool output total.
- **142** of those show `Compiling` / `Checking` spam (**~139k** chars).
- Rough quietable share of cargo text ≈ **207k chars (~52k tokens)** of *raw* output across the corpus — small per call, large once re-sent.

Cumulative re-send proxy on one heavy engine agent (`bc-5925cb44`, 207 shell calls, 265k raw shell chars): if each shell result is re-included in every later step, proxy cost ≈ **24.8M chars**. Capping each shell result at **800 chars** drops that proxy by **~63%**. That is not a billing meter; it shows why early quieting compounds.

### What actually dumps (ranked)

1. **`read_file` of large sources** — dominates tool chars (often 70%+). Out of scope for quiet-run; already covered by `context-management`.
2. **`git show` / `git diff` / `git grep` of big files** — often 10–35k chars per call. Skill rule: prefer `sed -n`, `rg`, path-scoped diff, or write to a file and read ranges.
3. **`cargo build` / `cargo test` / `cargo clippy`** — usually 1–7k when green; Compiling lines are pure waste. Failures need the error block kept.
4. **`go test ./...`** — seen at ~5k in genos-related runs; less noisy than cargo but still uncapped.

### Existing habits worth encoding

- Fleet rule already: “Scripts print short summaries and write full evidence to files.”
- `context-management` already: targeted reads, don’t paste long transcripts.
- Some genos agents already: `make test > /tmp/make-test.out 2>&1; echo EXIT; tail`.

Gap: nothing in **repo config** or a **shared wrapper** makes the quiet path the default.

---

## 1) Repo structure (trust pyramid: code first)

### 1a. Quiet defaults in each project

**genos-engine (first):**

```toml
# .cargo/config.toml
[term]
quiet = true          # suppress Compiling/Downloading progress; errors still print
# optional later: [build] jobs = N — not required for quieting
```

Makefile (additive targets; keep today’s targets for humans who want noise):

```make
LOGDIR ?= /tmp/genos-engine-logs
$(LOGDIR):
	mkdir -p $(LOGDIR)

# Full log on disk; stdout is verdict + path only.
test-quiet: $(LOGDIR)
	@./scripts/quiet-run --label test --log $(LOGDIR)/test.log -- \
	  cargo test --workspace --quiet -- --test-threads=8

lint-quiet: $(LOGDIR)
	@./scripts/quiet-run --label lint --log $(LOGDIR)/lint.log -- \
	  ./scripts/lint.sh

build-quiet: $(LOGDIR)
	@./scripts/quiet-run --label build --log $(LOGDIR)/build.log -- \
	  cargo build --workspace --quiet
```

CI (`.github/workflows/lint.yml`): keep full logs as job artifacts; console step can use the quiet wrapper so the Actions UI stays short. Artifacts are free for humans; agents rarely scrape the full CI log.

**genos (second):**

- `make test-quiet` / `make lint-style-quiet` wrapping `go test` and stylelint the same way.
- Prefer `go test -count=1 ./…` with `-json` piped through a tiny filter that prints only fail events + summary (optional phase 2; start with quiet-run + `go test` exit code).

### 1b. Shared wrapper in agent-skills: `scripts/quiet-run`

Auditable, ~30 lines, no network, no deps beyond bash. Interface:

```text
quiet-run --label <name> --log <path> [--keep-tail N] -- <command...>
```

Behaviour:

1. Run `command` with stdout+stderr teed to `--log` (create parent dirs).
2. Capture exit code.
3. Print exactly:
   - `PASS <label>` or `FAIL <label> exit=<code>`
   - up to `--keep-tail` lines (default 40) that match `error|Error|FAILED|panic|FAIL:` from the log (or last 20 lines on failure if no match)
   - `log: <absolute-path>`
4. Exit with the command’s status.
5. Never rewrite arguments; never fetch; never strip secrets specially (log file stays local — same trust as the command itself).

Sketch:

```bash
#!/usr/bin/env bash
set -euo pipefail
# parse --label --log --keep-tail; require -- then cmd
mkdir -p "$(dirname "$log")"
set +e
"$@" > >(tee "$log" >/dev/null) 2>&1
# note: simpler portable form:
set +e
"$@" >"$log" 2>&1
ec=$?
set -e
if [[ $ec -eq 0 ]]; then echo "PASS $label"; else echo "FAIL $label exit=$ec"; fi
# print matched error lines or tail on failure
echo "log: $(cd "$(dirname "$log")" && pwd)/$(basename "$log")"
exit "$ec"
```

Ship under `agent-skills/scripts/quiet-run`, with a one-line test (fake failing/passing commands). Repos vendor or submodule-copy isn’t required: agents can `curl`/`cp` from a checked-out agent-skills, or each repo gets a 5-line `scripts/quiet-run` that is a copy (prefer copy for CI hermeticity; agent-skills remains source of truth).

**Why not RTK / npm proxies:** any opaque interceptor that rewrites tool I/O is a supply-chain and “what did the agent actually see?” risk. A 30-line bash tee is reviewable in one screen.

---

## 2) Coding skill rule (minimal prose)

Add under `skills/coding/references/principles.md` (or a 10-line `references/quiet-output.md` linked from Routing) — **exact proposed text (≤4 lines of rule body):**

```markdown
## Quiet tool output
Prefer quiet flags and repo `*-quiet` / `quiet-run` targets. Filter with `rg`/`jq` before reading; open log files at specific lines. Never paste a full build, test, or `git show` log into the chat — print PASS/FAIL, the error slice, and the log path.
```

Trigger stays on existing `coding` description (already ALWAYS for build/CI/debug). Do **not** lengthen the skill listing description for this.

`context-management` already overlaps; keep one sentence there pointing at the coding rule so the two don’t drift (“build/test logs: coding quiet-output rule”).

---

## 3) Self-correction (agents fix the source)

### Problem

Cloud agents in genos / genos-engine notice huge tool results mid-task, then keep going without changing Makefile / cargo config / CI, so the next agent pays again.

### Mechanism (low listing cost)

**Do not** add a new always-on skill to the global listing if avoidable. Prefer:

**Option A (recommended): reference from `coding`, loaded on demand**

New thin skill `quiet-tool-output` (or a section only) with a **short** description that fires on *what the agent is doing now*:

```yaml
description: >
  Use when a shell/build/test/lint command just returned a large or repetitive
  transcript (Compiling/Downloading spam, hundreds of pass lines, multi-page
  git show), when the same noisy command was run twice this task, or when
  editing Makefile/cargo/CI to reduce agent log noise. Fix the repo default or
  wrapper; do not only re-run with ad-hoc redirects.
```

Listing cost: ~40–60 tokens if global; better as **repo-local** under genos-engine / genos `.agents/skills/` (cloud agents already miss global agent-skills — see skills-visibility probe). Until the Cursor plugin ships skills to cloud agents, put the skill (or AGENTS.md bullet) **in the product repos**.

**Closing step in `coding` / review** (one line, mirrors lint-habit):

> After a command dumps large repetitive output, add or fix a quiet default (`*-quiet`, `.cargo/config.toml`, CI flag) in the same PR when the repo is yours; note `quiet-output: fixed|skipped (reason)` in the end summary.

### What “fix the source” means

| Symptom | Fix (structure, not prose) |
|---|---|
| `Compiling` walls on every `cargo test` | `.cargo/config.toml` `[term] quiet=true` + `--quiet` on agent targets |
| Full `make lint` in chat | `make lint-quiet` via quiet-run; AGENTS.md says agents use `*-quiet` |
| Agent re-runs `git show whole.rs` | coding rule: range/`rg` only; optional repo script `scripts/show-slice` |
| CI log floods Actions UI | quiet console + upload full log artifact |

---

## 4) Measurement (Monday weekly agent-history review)

Extend the existing Monday review (cloud-agent-runs.csv + transcript dump) with a **big-output flag**:

**Threshold (starting point):**

- Any single `run_terminal_cmd` result with `len(output) >= 8_000` chars, **or**
- Any command matching `cargo |go test|go build|make |npm |clippy` with `len(output) >= 4_000`, **or**
- Same command prefix (first 80 chars) appearing **≥3 times** in one transcript with total output ≥ 10_000 chars.

**Report (per noisy agent, top 15/week):**

```text
date | agent_id | repo | cmd_prefix | max_chars | times | quietable? (cargo/go/make) | notes
```

**Action:** Skills bot opens or updates issues/PRs on genos-engine/genos when the same cmd_prefix repeats across ≥2 agents in a week. No auto-commit.

Script location (post-reset build-out): `agent-skills/scripts/flag-big-tool-output.py` reading `/workspace/cloud-agent-transcripts/*.jsonl`.

---

## Phased plan

| Phase | When | Work | Est. savings (evidenced) |
|---|---|---|---|
| **0. This PR** | now | Proposal only | — |
| **1. genos-engine** | after Tue reset | `.cargo/config.toml` quiet; `make {build,test,lint}-quiet` + copied `quiet-run`; AGENTS.md one bullet “agents use `*-quiet`”; coding rule text | ~40–60% of cargo tool chars on green builds; on heavy agents, large cut to cumulative resend (proxy ~63% if shell capped ~800 chars) |
| **2. agent-skills** | same window | Land `scripts/quiet-run` + coding rule + optional thin skill; flag script for Monday review | Listing +~0 if rule is reference-only; +~50 tokens if new skill description |
| **3. genos** | after engine proves out | `make test-quiet` / stylelint-quiet; mirror AGENTS.md bullet | Smaller absolute win than engine; stops ad-hoc `/tmp` redirects from diverging |
| **4. Measurement** | first Monday after land | Flag report in weekly review; tune thresholds | Feedback loop |

Estimated token savings are **per continued turn after a noisy command**, not a one-shot discount. Biggest win is structural quiet early in long engine agents (gi / lighting PRs).

---

## Open questions

1. **`.cargo/config.toml` `[term] quiet=true` for humans too?** Or only agent `*-quiet` targets so interactive `make lighting` stays chatty? (Recommend: repo quiet term + verbose opt-in via `CARGO_TERM_QUIET=false`, so agents inherit quiet without remembering flags.)
2. **Should `quiet-run` live only in agent-skills, or be copied into each repo for hermetic CI?** (Recommend: source of truth in agent-skills, copy into genos-engine/genos `scripts/` with a comment pointing back.)
3. **Cloud agents still don’t load Anthony’s agent-skills** (probe on genos). Until the Cursor plugin is installed on the account, quiet behaviour must be **in-repo** (AGENTS.md + make targets) or it won’t fire in the place that burns the most tokens.
4. **Git dumps vs build dumps:** build quieting is easy; `git show` of multi-kLOC files may need a separate one-liner in `git` skill. In or out of this proposal’s phase 1?
5. **Threshold 8k/4k:** confirm after one Monday report whether false positives (intentional `gh pr view --json body`) need an allowlist.

## Non-goals

- Third-party output proxies / npm interceptors.
- Truncating tool output inside Cursor itself (we don’t control the harness).
- Replacing `context-management` or fleet token-habits prose.
- Changing model/effort tiers.

## Decision asked of Anthony

Approve phase 1 direction (engine quiet defaults + quiet-run + 4-line coding rule + Monday flag) after the usage reset? Any preference on questions 1 and 3?
