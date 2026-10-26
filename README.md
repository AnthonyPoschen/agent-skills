# Agent Skills

Reusable Agent Skills for coding agents that support the open `SKILL.md`
format. This repository is arranged so skills can be installed with the
`skills` CLI:

```bash
npx skills add anthonyposchen/agent-skills --list
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill act-within-scope
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill build-files
npx skills add anthonyposchen/agent-skills --skill domain-modeling
npx skills add anthonyposchen/agent-skills --skill frontend-design
npx skills add anthonyposchen/agent-skills --skill git
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux
npx skills add anthonyposchen/agent-skills --skill grill-with-docs
npx skills add anthonyposchen/agent-skills --skill security-audit
npx skills add anthonyposchen/agent-skills --skill in-your-own-words
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill writing
npx skills add anthonyposchen/agent-skills --skill test-audit
npx skills add anthonyposchen/agent-skills --skill scripted-verification
```

## Available Skills

| Skill | Description |
| --- | --- |
| `act-within-scope` | Make reversible progress within scope without waiting on ordinary execution decisions. |
| `coding` | Coding, refactor, review, debugging, and design-task guardrails for clean, idiomatic, maintainable code. Includes comment audits (no-comments) and type-system discipline. |
| `context-management` | Keep large investigations and multi-phase work focused with targeted reads, compact handoffs, and bounded delegation. |
| `build-files` | Build files: secure multi-stage Dockerfiles (scratch-first, non-root, Docker as CI source of truth) and Makefile targets that validate required variables. |
| `domain-modeling` | Sharpen domain terminology and ADRs; optionally refresh short living alignment docs when the repo declares them. |
| `frontend-design` | Product-fit frontend design and implementation workflow for polished, usable web interfaces. |
| `git` | Commits, worktrees, and releases: commit boundaries and messages, a mirrored worktree layout, and release notes, versions, and tags (CalVer when no convention exists). |
| `encode-lessons` | Turn recurring corrections and preventable mistakes into the smallest durable guard. |
| `experience-first` | Keep product, interface, and API decisions focused on the consumer's core workflow. |
| `flux` | FluxCD GitOps: GitHub Actions to GHCR image automation (tags, policies, rollout diagnostics) and Kustomize base/overlay repo layouts. |
| `grill-with-docs` | Relentless design interview that drives domain-modeling so GLOSSARY.md/ADRs (and declared living shorts) stay current. |
| `security-audit` | **Manual only** (`/security-audit`). Vendored Cloudflare six-phase security audit; agents must not auto-load it. |
| `in-your-own-words` | Restate the user goal, problem, constraints, and done check in plain words before non-trivial work. |
| `lint-rules` | Shared custom lint rules for Go, Rust, Kubernetes and CSS, each with a failing example, plus how to write and contribute new ones. |
| `scripted-verification` | Turn repeated agent judgement into deterministic PASS, FAIL, or UNSURE scripts with exit codes and evidence on disk. |
| `skill-evaluation` | Evidence based comparison of agent skill triggering and task outcomes using realistic prompts and blinded review. |
| `writing` | Prose for people: ASD-STE100 Simplified Technical English for every reply, plus plain English guidance for docs, READMEs, RFCs, PR descriptions, and release notes. |
| `test-behavior-not-implementation` | Call code as its users do and assert the exact result they see; reject loose or call-only assertions. |
| `test-audit` | Decide whether a proposed or existing test should exist, and prune tests that do not protect a real contract. Delete tautologies that only restate the code. |

## Install

List installable skills:

```bash
npx skills add anthonyposchen/agent-skills --list
```

Install a specific skill:

```bash
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill act-within-scope
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill build-files
npx skills add anthonyposchen/agent-skills --skill domain-modeling
npx skills add anthonyposchen/agent-skills --skill frontend-design
npx skills add anthonyposchen/agent-skills --skill git
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux
npx skills add anthonyposchen/agent-skills --skill grill-with-docs
npx skills add anthonyposchen/agent-skills --skill security-audit
npx skills add anthonyposchen/agent-skills --skill in-your-own-words
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill writing
npx skills add anthonyposchen/agent-skills --skill test-audit
```

Install all skills:

```bash
npx skills add anthonyposchen/agent-skills --all
```

Install globally for a specific agent:

```bash
npx skills add anthonyposchen/agent-skills --skill coding --global --agent codex
```

## Cursor plugin

This repository is also a Cursor plugin named `agent-skills`. The plugin's
skills are the existing `skills/` folder, so `npx skills add` is unchanged.

1. In Cursor, open **Customize** in the sidebar.
2. Import this repository with **From GitHub Repository**:
   `https://github.com/AnthonyPoschen/agent-skills`.
3. Find the `agent-skills` plugin, select **Install**, and choose user scope.

Installed plugin skills are available to Cloud Agents. To check, start a
Cloud Agent and ask it to list its skills.

## Shared MCP config

Cross-machine MCP server definitions live in `config/mcp/`. Cua Driver is
`config/mcp/cua-driver.json`: the `cua-driver` command with the `mcp` argument.

After Cua Driver is installed on a machine, register it with the local agents:

```bash
npm run install-cua-mcp
```

The script resolves `cua-driver` on that machine and writes the user-level
Grok, Codex, Claude Code, OpenCode, Cursor, and Gemini entries. Run it again
after moving or upgrading the binary.

## Repository Layout

```text
skills/
  act-within-scope/
    SKILL.md
  coding/
    SKILL.md
    examples/
    references/
  context-management/
    SKILL.md
  create-verification-skill/
    SKILL.md
    references/
  build-files/
    SKILL.md
    references/
  domain-modeling/
    SKILL.md
    GLOSSARY-FORMAT.md
    ADR-FORMAT.md
  frontend-design/
    SKILL.md
  git/
    SKILL.md
    references/
    scripts/
  encode-lessons/
    SKILL.md
  experience-first/
    SKILL.md
  eyes-on/
    SKILL.md
    scripts/
  flux/
    SKILL.md
    references/
  grill-with-docs/
    SKILL.md
  security-audit/
    SKILL.md
    LICENSE-UPSTREAM
    *.md reference files, validators, report-schema.json
  in-your-own-words/
    SKILL.md
  maintain-verification-skill/
    SKILL.md
  writing/
    SKILL.md
    references/
  scripted-verification/
    SKILL.md
    scripts/
  skill-evaluation/
    SKILL.md
  test-audit/
    SKILL.md
    references/
```

Each skill lives in `skills/<skill-name>/` and must include a `SKILL.md` file
with YAML frontmatter containing `name` and `description`.

Optional skill resources:

- `scripts/` for deterministic helper scripts.
- `references/` for supporting docs loaded only when needed.
- `assets/` for templates, images, fonts, or other reusable output resources.
- `agents/openai.yaml` for Codex-facing UI metadata.

## Add A Skill

1. Create `skills/<new-skill-name>/SKILL.md`.
2. Use lowercase kebab-case for the folder and `name` field.
3. Keep triggering context in the frontmatter `description`.
4. Keep `SKILL.md` concise; move long details into `references/`.
5. Link local agent skills:

```bash
npm run link-skills
```

6. Run validation:

```bash
npm run validate
```

## Local Development

Install Git hooks that relink skills after `git pull` updates the working tree:

```bash
npm run install-git-hooks
```

This installs `post-merge` and `post-rewrite` hooks so both merge-based pulls
and rebase-based pulls refresh local links.

## Two-Way Skills Sync With The Bot Box

**Intent.** The skills in `skills/` and the shared bot skills on the box
(`/home/box/agent-data/workflows`) should stay one set. An edit on either
side, whether from a bot, a person or a cloud agent, should reach the other. This
repo's `master` is canonical, so box edits arrive here as PRs for review and are
never pushed straight to `master`.

**How it decides.** `scripts/sync-skills.sh` compares each skill on both sides
against the last synced version, recorded in a state file on the box (outside
this repo). If only the repo changed, the change is copied to the box. If only
the box changed, the sync opens or updates a `skills-sync/<skill>` PR. Deletions
follow the same rule. Skills listed in `sync-exclude.txt` are bot-only and are
never touched. A PR does not count as synced until it is merged, so the sync
never triggers itself.

**Run it.** `scripts/sync-skills.sh --help` lists the options and defaults. Add
`--dry-run` to print the plan without changing anything. Tests run offline with
`bash tests/sync-skills-test.sh`.

**Conflicts.** If both sides changed the same skill, or one side deleted it
while the other edited it, the sync prints a `FLAG conflict` line, opens a PR
carrying the box version with the conflict noted in its body, and leaves both
sides as they are. Resolve it in that PR, or close the PR and fix one side. The
skill settles once both sides match.

**Routine (add after merge).** Run the script every 15 minutes on the box as a
Grok Bot routine. Surface `FLAG` lines to Anthony.

## Manual-only skills

Some skills set `disable-model-invocation: true` so agents do not auto-load them.
Invoke them by name (for example `/security-audit`).

| Harness | Honours manual-only? | How |
| --- | --- | --- |
| Cursor (Agent Skills / plugin) | Yes | Top-level `disable-model-invocation: true`. Still available via `/security-audit`. The Cursor plugin manifest points at all of `skills/`; there is no separate exclude list. |
| Claude Code | Yes | Same frontmatter field. Description is kept out of the auto-listing; `/security-audit` still works. |
| OpenCode | Yes (via metadata) | OpenCode ignores `disable-model-invocation`. This skill sets `metadata.opencode/autoinvoke: "false"` and `metadata.opencode/slash: "true"` so it stays out of `<available_skills>` but remains a `/security-audit` slash command. |
| Grok Build / this repo's link and sync scripts | No equivalent flag | `scripts/link-agent-skills.sh`, `scripts/sync-skills.sh`, and `scripts/validate-skills.py` still link, sync, and validate every skill folder, including manual-only ones. Grok has no documented `disable-model-invocation` or `opencode/autoinvoke` equivalent, so the skill may still appear in a Grok skill listing after `npm run link-skills`. Prefer explicit `/security-audit` (or an equivalent ask) rather than relying on auto-routing. |

`security-audit` is vendored from [cloudflare/security-audit-skill](https://github.com/cloudflare/security-audit-skill) at commit `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8` (MIT). See `skills/security-audit/LICENSE-UPSTREAM`.

## License

MIT, unless a skill directory includes its own `LICENSE.txt` or `LICENSE-UPSTREAM`.
