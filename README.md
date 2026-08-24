# Agent Skills

Reusable Agent Skills for coding agents that support the open `SKILL.md`
format. This repository is arranged so skills can be installed with the
`skills` CLI:

```bash
npx skills add anthonyposchen/agent-skills --list
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill dockerfile
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux-ghcr-image-automation
npx skills add anthonyposchen/agent-skills --skill flux-kustomize-layout
npx skills add anthonyposchen/agent-skills --skill implement-tickets
npx skills add anthonyposchen/agent-skills --skill makefile
npx skills add anthonyposchen/agent-skills --skill omarchy-development
npx skills add anthonyposchen/agent-skills --skill release
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill technical-writing
```

## Available Skills

| Skill | Description |
| --- | --- |
| `coding` | Coding, refactor, review, debugging, and design-task guardrails for clean, idiomatic, maintainable code. |
| `context-management` | Keep large investigations and multi-phase work focused with targeted reads, compact handoffs, and bounded delegation. |
| `dockerfile` | Secure multi-stage Dockerfile patterns with scratch-first static binaries, Alpine runtimes, non-root UID/GID ownership, and Docker-based CI source of truth. |
| `encode-lessons` | Turn recurring corrections and preventable mistakes into the smallest durable guard. |
| `experience-first` | Keep product, interface, and API decisions focused on the consumer's core workflow. |
| `flux-ghcr-image-automation` | GitHub Actions to GHCR to FluxCD image automation loops, including image tags, policies, setter comments, namespace gotchas, and rollout diagnostics. |
| `flux-kustomize-layout` | FluxCD and Kustomize repo layout scaffolds with shared base plus dev/prod overlays. |
| `implement-tickets` | Dependency-aware selected work-item orchestration with isolated workers, review feedback, and human-controlled merges. |
| `makefile` | Makefile target patterns with required variable checks and useful input descriptions before commands run. |
| `omarchy-development` | Isolated Omarchy source development and VM testing with disposable, Virtual Machine Manager-managed KVM overlays. |
| `release` | Release records, notes, version selection, tags, and publishing. Uses the repository versioning convention or CalVer when none exists. |
| `skill-evaluation` | Evidence based comparison of agent skill triggering and task outcomes using realistic prompts and blinded review. |
| `technical-writing` | Plain English guidance for documentation, READMEs, RFCs, PR descriptions, release notes, and commit messages. |

## Install

List installable skills:

```bash
npx skills add anthonyposchen/agent-skills --list
```

Install a specific skill:

```bash
npx skills add anthonyposchen/agent-skills --skill coding
npx skills add anthonyposchen/agent-skills --skill context-management
npx skills add anthonyposchen/agent-skills --skill dockerfile
npx skills add anthonyposchen/agent-skills --skill encode-lessons
npx skills add anthonyposchen/agent-skills --skill experience-first
npx skills add anthonyposchen/agent-skills --skill flux-ghcr-image-automation
npx skills add anthonyposchen/agent-skills --skill flux-kustomize-layout
npx skills add anthonyposchen/agent-skills --skill implement-tickets
npx skills add anthonyposchen/agent-skills --skill makefile
npx skills add anthonyposchen/agent-skills --skill omarchy-development
npx skills add anthonyposchen/agent-skills --skill release
npx skills add anthonyposchen/agent-skills --skill skill-evaluation
npx skills add anthonyposchen/agent-skills --skill technical-writing
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
  coding/
    SKILL.md
    examples/
  context-management/
    SKILL.md
  dockerfile/
    SKILL.md
    references/
  encode-lessons/
    SKILL.md
  experience-first/
    SKILL.md
  flux-ghcr-image-automation/
    SKILL.md
    references/
  flux-kustomize-layout/
    SKILL.md
  implement-tickets/
    SKILL.md
    references/
  makefile/
    SKILL.md
  release/
    SKILL.md
    references/
  technical-writing/
    SKILL.md
  skill-evaluation/
    SKILL.md
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

## License

MIT, unless a skill directory includes its own `LICENSE.txt`.
