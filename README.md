# Agent Skills

Reusable Agent Skills for coding agents that support the open `SKILL.md`
format. This repository is arranged so skills can be installed with the
`skills` CLI:

```bash
npx skills add anthonyposchen/agent-skills --list
npx skills add anthonyposchen/agent-skills --skill coding
```

## Available Skills

| Skill | Description |
| --- | --- |
| `coding` | Coding, refactor, review, debugging, and design-task guardrails for clean, idiomatic, maintainable code. |

## Install

List installable skills:

```bash
npx skills add anthonyposchen/agent-skills --list
```

Install a specific skill:

```bash
npx skills add anthonyposchen/agent-skills --skill coding
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
5. Run validation:

```bash
npm run validate
```

## License

MIT, unless a skill directory includes its own `LICENSE.txt`.
