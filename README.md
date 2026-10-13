# Agent Skills

Reusable Agent Skills for coding agents that support the open `SKILL.md`
format. This repository is arranged so skills can be installed with the
`skills` CLI:

```bash
npx skills add anthonyposchen/agent-skills --list
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill asd-ste100
npx skills add anthonyposchen/agent-skills --skill act-within-scope
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill dockerfile
npx skills add anthonyposchen/agent-skills --skill domain-modeling
npx skills add anthonyposchen/agent-skills --skill frontend-design
npx skills add anthonyposchen/agent-skills --skill git-worktree
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux-ghcr-image-automation
npx skills add anthonyposchen/agent-skills --skill flux-kustomize-layout
npx skills add anthonyposchen/agent-skills --skill grill-with-docs
npx skills add anthonyposchen/agent-skills --skill makefile
npx skills add anthonyposchen/agent-skills --skill omarchy-development
npx skills add anthonyposchen/agent-skills --skill release
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill technical-writing
npx skills add anthonyposchen/agent-skills --skill test-audit
```

## Available Skills

| Skill | Description |
| --- | --- |
| `asd-ste100` | Write replies and documents in ASD-STE100 Simplified Technical English. |
| `act-within-scope` | Make reversible progress within scope without waiting on ordinary execution decisions. |
| `coding` | Coding, refactor, review, debugging, and design-task guardrails for clean, idiomatic, maintainable code. |
| `context-management` | Keep large investigations and multi-phase work focused with targeted reads, compact handoffs, and bounded delegation. |
| `dockerfile` | Secure multi-stage Dockerfile patterns with scratch-first static binaries, Alpine runtimes, non-root UID/GID ownership, and Docker-based CI source of truth. |
| `domain-modeling` | Sharpen domain terminology and ADRs; optionally refresh short living alignment docs when the repo declares them. |
| `frontend-design` | Product-fit frontend design and implementation workflow for polished, usable web interfaces. |
| `git-worktree` | Create, synchronize, and safely retire isolated Git worktrees in a predictable mirrored layout. |
| `encode-lessons` | Turn recurring corrections and preventable mistakes into the smallest durable guard. |
| `experience-first` | Keep product, interface, and API decisions focused on the consumer's core workflow. |
| `flux-ghcr-image-automation` | GitHub Actions to GHCR to FluxCD image automation loops, including image tags, policies, setter comments, namespace gotchas, and rollout diagnostics. |
| `flux-kustomize-layout` | FluxCD and Kustomize repo layout scaffolds with shared base plus dev/prod overlays. |
| `grill-with-docs` | Relentless design interview that drives domain-modeling so CONTEXT/ADRs (and declared living shorts) stay current. |
| `makefile` | Makefile target patterns with required variable checks and useful input descriptions before commands run. |
| `omarchy-development` | Isolated Omarchy source development and VM testing with disposable, Virtual Machine Manager-managed KVM overlays. |
| `release` | Release records, notes, version selection, tags, and publishing. Uses the repository versioning convention or CalVer when none exists. |
| `skill-evaluation` | Evidence based comparison of agent skill triggering and task outcomes using realistic prompts and blinded review. |
| `technical-writing` | Plain English guidance for documentation, READMEs, RFCs, PR descriptions, release notes, and commit messages. |
| `test-audit` | Decide whether a proposed or existing test should exist, and prune tests that do not protect a real contract. |

## Install

List installable skills:

```bash
npx skills add anthonyposchen/agent-skills --list
```

Install a specific skill:

```bash
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill asd-ste100
npx skills add anthonyposchen/agent-skills --skill act-within-scope
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill dockerfile
npx skills add anthonyposchen/agent-skills --skill domain-modeling
npx skills add anthonyposchen/agent-skills --skill frontend-design
npx skills add anthonyposchen/agent-skills --skill git-worktree
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux-ghcr-image-automation
npx skills add anthonyposchen/agent-skills --skill flux-kustomize-layout
npx skills add anthonyposchen/agent-skills --skill grill-with-docs
npx skills add anthonyposchen/agent-skills --skill makefile
npx skills add anthonyposchen/agent-skills --skill omarchy-development
npx skills add anthonyposchen/agent-skills --skill release
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill technical-writing
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

## Repository Layout

```text
skills/
  asd-ste100/
    SKILL.md
  act-within-scope/
    SKILL.md
  coding/
    SKILL.md
    examples/
  context-management/
    SKILL.md
  create-verification-skill/
    SKILL.md
    references/
  dockerfile/
    SKILL.md
    references/
  domain-modeling/
    SKILL.md
    CONTEXT-FORMAT.md
    ADR-FORMAT.md
  frontend-design/
    SKILL.md
  git-worktree/
    SKILL.md
    scripts/
  encode-lessons/
    SKILL.md
  experience-first/
    SKILL.md
  eyes-on/
    SKILL.md
    scripts/
  flux-ghcr-image-automation/
    SKILL.md
    references/
  flux-kustomize-layout/
    SKILL.md
  grill-with-docs/
    SKILL.md
  makefile/
    SKILL.md
  maintain-verification-skill/
    SKILL.md
  release/
    SKILL.md
    references/
  technical-writing/
    SKILL.md
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

## License

MIT, unless a skill directory includes its own `LICENSE.txt`.
