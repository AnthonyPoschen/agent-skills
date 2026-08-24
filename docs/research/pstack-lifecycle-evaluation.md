# Pstack lifecycle evaluation

## Scope and conclusion

This note compares Pstack at commit
[`4612556`](https://github.com/cursor/plugins/tree/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack)
with the public skills in this repository. It excludes the coding rules already
adopted here: caller focused boundaries, scoped simplification, direct outcome
proof, selective stable contract tests, and discoverable file layout.

Pstack has useful lifecycle ideas, but its `poteto-mode` is not a good import.
It is a Cursor specific sticky router with twenty two playbooks, model roles,
parallel agents, and unattended delivery assumptions. This repository already
has focused skills and an explicit human controlled ticket workflow. Import the
small practices that fill a real gap. Do not recreate Pstack as a second
orchestration system.

## Already covered

| Pstack idea | Existing coverage | Evaluation |
| --- | --- | --- |
| Caller first architecture and meaningful boundaries | `coding/references/design.md` | Already stronger and better matched to this library. It makes direct cohesive code a valid outcome, not a failure to abstract. Pstack's [architect](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/architect/SKILL.md) reaches the same consumer first design goal, but requires several model passes and sketches for every nontrivial change. |
| Real artifact verification and root cause fixes | `coding/references/verification.md` and `debugging.md` | Already covered. The local rule is more precise for this user: prove the requested outcome through the real path where practical, and test only stable meaningful contracts. Pstack's [prove it works](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-prove-it-works/SKILL.md) and [fix root causes](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-fix-root-causes/SKILL.md) agree with the direction. |
| Domain language and domain shaped structures | `domain-modeling` and `coding/references/design.md` | Covered. Pstack's [model the domain](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-model-the-domain/SKILL.md) is a useful reminder, but its central advice already belongs in the existing domain modeling and coding design workflows. |
| Throwaway prototypes to answer a decision | `prototype` | Covered. Pstack's [prototype playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/prototype.md) agrees that the prototype proves a decision, not production code. |
| Ticket execution, commits, review, releases | `implement-tickets`, `git-commit-workflow`, `changelog`, and `calver-release` | Covered, with better safety for this library. Pstack's central [mode and playbooks](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/README.md) assume a particular Cursor and Graphite workflow. |

## Strong candidates

### Technical writing skill

Add a small standalone `technical-writing` skill, adapted rather than copied.
Use it for public documentation, READMEs, how to guides, architecture notes,
ADRs, PR descriptions, and release text.

Pstack's [technical writing skill](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/technical-writing/SKILL.md)
has a good core: choose the document's purpose, write for a tired engineer,
use ordinary words, state the real symbol or path, and avoid dense noun piles.
For this repository, keep it short and write its examples in plain English.
Do not import every named writing standard or turn it into a rigid sentence
checker. This directly supports the user's preference for simple language and
clear file paths.

### Blinded evaluation for skill changes

This is a skills library. Once a skill has competing wording or a meaningful
change, evaluate behavior rather than deciding from a convincing sounding
review. Give each candidate the same ordinary task in a separate workspace,
hide which variant it received, and judge the results against one shared rubric.
Read the produced work and transcripts rather than trusting a claim that the
skill was followed.

Pstack's [eval playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/eval.md)
is a strong basis for this. Adapt it to the available agent runner. Use it only
when a change has enough impact or uncertainty to justify the experiment; do
not benchmark every wording correction.

### Codebase understanding and decision archaeology

Add a lean skill, or a coding reference pair, for two distinct questions:

- **How does this work now?** Trace entry points, callers, data flow, side
  effects, and ownership before a change that needs system understanding.
- **Why is it this way?** Inspect commits, pull requests, issues, ADRs, and
  documents. Cite evidence, distinguish fact from inference, and state gaps.

Pstack separates these in [how](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/how/SKILL.md)
and [why](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/why/SKILL.md).
The separation is valuable. The implementation should be smaller and tool
agnostic: use only records the current repository and connected systems make
available. Do not require parallel agents or pretend history proves intent.

### Blast radius as a review and change planning check

Add a short optional route to `coding/references/review.md`, rather than a
full new skill. Use it for a change that looks small but crosses a public API,
schema, configuration, lifecycle, cache, concurrency, or shared helper.

The rule is: trace real consumers and affected states, name the specific thing
that could break, then run the direct check that establishes the safety claim.
Pstack's [blast radius skill](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/blast-radius/SKILL.md)
is particularly compatible with the local outcome proof rule because it asks
for running evidence rather than a speculative impact list.

### Measured performance work

Extend `diagnose`, or add a small performance reference, with a distinct
performance loop: capture a baseline on the user relevant path, use traces or
measurements to form hypotheses, make one focused change, capture the same
measurement after, and report the delta plus artifact path. A source reading
is not performance evidence.

Pstack's [performance playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/perf-issue.md)
has the right shape. Keep only the baseline, trace, before and after, and
inconclusive is not a pass rules. Do not copy its long menu of optimization
patterns as a checklist.

### Idempotence for retryable state changes

Add a focused design and review rule for commands, jobs, importers, migrations,
provisioning, and lifecycle operations that can retry after a crash. Ask what
happens on a second run and after partial completion. Require convergence or an
explicit reason the operation is single use.

This is narrow, practical, and often missed. It comes from Pstack's
[idempotence principle](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-make-operations-idempotent/SKILL.md).
It should not become a demand that every in memory helper be idempotent.

### Finish a bounded internal API migration

The current coding references already say to inspect direct consumers, migrate
a bounded set together, and avoid a new workaround or compatibility layer. Add
one explicit finishing rule: when the old API has no external compatibility
commitment, migrate every in scope caller and delete the old API in the same
coherent change. Do not retain a permanent internal adapter merely to postpone
the deletion.

That is Pstack's [migrate callers then delete legacy APIs](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-migrate-callers-then-delete-legacy-apis/SKILL.md).
It is an incremental clarification of an adopted rule, not a new skill. Keep
the existing scope guard: if every caller cannot be migrated safely in the
task, make no half migration and report the wider work.

## Situational candidates

### Project local verification skill and feature map

Pstack's [create verification skill](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/create-verification-skill/SKILL.md)
creates a project local recipe that launches the real application, drives it,
captures evidence, and cleans only what it created. Its companion
[maintenance skill](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/maintain-verification-skill/SKILL.md)
checks a feature map against source and a live run.

This is an excellent fit for a serious runnable application that agents change
often, especially where local databases and real UI or API flows are available.
It directly reinforces the ground truth preference. It is not needed for a
library, a small service with one clear verification command, or a project
without a maintainable local environment. If added, adapt the storage location
to this repository's skill conventions and create it only on request.

### Verifiable units for bulk changes

For data migrations, mechanical edits, package moves, or a long sequence of
related changes, work in small units that each end in a usable checked state.
Pstack calls this [sequence verifiable units](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-sequence-verifiable-units/SKILL.md).

Use the direct proof appropriate to each unit. Do not require a red then green
test commit, since that would conflict with the local selective test rule. This
belongs as a conditional rule in migration and refactoring guidance, not a
default process for a normal small edit.

### Reflection after an expensive or surprising task

Pstack's [reflect skill](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/reflect/SKILL.md)
turns a completed run into a durable skill edit. That fits a skills repository,
but only after a repeated failure, a costly discovery, or an explicit request.
Routine self modification will create noisy rules and make skills worse. A
small post task question is enough: what repeated mistake or missing project
fact would a concise guardrail prevent?

### Decision trail for high trust work

Pstack's [show me your work](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/show-me-your-work/SKILL.md)
records decisions and evidence for review. Use a compact version only for a
multi day migration, high risk incident, audit requirement, or unattended run.
For ordinary changes, the diff, commit message, and direct proof are enough.

### Script the repeatable fragile operation

Pstack's [build the lever](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/principle-build-the-lever/SKILL.md)
rightly values a rerunnable codemod or check. Its default to build a tool for
any nontrivial work is too broad for this library. Apply it only when the same
fragile operation occurs enough times that a small script is simpler and more
trustworthy than manual work. The script must be smaller than the work it
replaces and must not become a framework.

## Not recommended

| Pstack item | Why not import it |
| --- | --- |
| `poteto-mode` and its playbook router | It overlaps many existing skills, relies on Cursor and model configuration, and adds a broad process layer to every substantial task. The repository benefits from explicit, small routing. [Pstack README](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/README.md) |
| Default multiple design candidates and model arena | A useful tool for a genuine one way door, but Pstack's architect requires alternatives for ordinary nontrivial work. That can turn straightforward design into ceremony and conflicts with the smallest coherent change approach. [Architect](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/architect/SKILL.md) |
| TDD as a general lifecycle rule | The user prefers direct real outcome proof and tests only for stable contracts. Pstack itself limits [TDD](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/tdd/SKILL.md) to cheap local seams, but the existing verification rule is the better universal decision rule. |
| Comment removal as a standing rule | Pstack's [no comments](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/no-comments/SKILL.md) is intentionally adversarial and can erase useful contracts, policy, and operational context. Keep the current rule: comments earn their place when they explain non obvious purpose or consequence. |
| Unattended shipping and automatic merge workflows | These require authority, reliable external tools, and organizational trust. The existing `implement-tickets` skill correctly keeps human merge authority. |
| Many tiny principle skills | Their ideas are useful, but copying them as separate skills would recreate Pstack's complexity. Fold only approved rules into the existing design, implementation, debugging, review, and verification references. |

## Suggested order

1. Add the plain language `technical-writing` skill when work resumes on it.
2. Add blinded evaluation before making larger skill changes. This repository
   needs evidence that a skill changed agent behavior, not only a good review.
3. Add codebase understanding and decision archaeology, because it improves
   choices before code changes begin.
4. Add blast radius and idempotence as short conditional checks in coding.
5. Improve performance diagnosis with baseline and after measurement.
6. Consider project local verification only for an application that needs a
   reusable real world drive and evidence recipe.

This order strengthens the lifecycle without replacing direct code, real
verification, or scoped change with a new process for its own sake.
