# Agent Skills Repository Guide

This repository publishes reusable Agent Skills. Keep the repo easy for
installers and agents to scan.

## Structure

- Put every public skill in `skills/<skill-name>/`.
- Use lowercase kebab-case for skill folders.
- Every skill must contain `SKILL.md`.
- Optional supporting resources belong in `scripts/`, `references/`, `assets/`,
  or `agents/`.
- Do not add README files inside individual skill directories; put public docs
  in the repo README and agent instructions in `SKILL.md`.

## Skill Rules

- `SKILL.md` frontmatter must contain only `name` and `description` unless a
  specific consuming tool requires more.
- The `name` value should match the folder name.
- Put all activation and routing language in `description`.
- Prefer concise, imperative instructions in the body.
- Link directly from `SKILL.md` to supporting reference files so agents can load
  them only when relevant.
- Add scripts only when they make repeated or fragile work more reliable.
- When adding a new skill folder, run `npm run link-skills` so the local
  `~/.agents/skills/<skill-name>` entry exists before validation.

## Validation

Run this before publishing changes:

```bash
npm run validate
```

The validator checks skill folder names, frontmatter boundaries, required
fields, duplicate names, references to missing local markdown files, and local
entries in `~/.agents/skills` for every `skills/<skill-name>/` directory.

## Releases

- Always prepare release notes when cutting a release unless the user explicitly
  says to skip release notes.
- If changelog or release-note evidence exists (`CHANGELOG.md`,
  `docs/releases/`, `docs/changelog/`, `releases/`, or `.github/releases/`),
  update it before tagging.
- In-repo changelog or release-note updates must be committed before tagging.
- Tag the commit that includes required release-note updates.
