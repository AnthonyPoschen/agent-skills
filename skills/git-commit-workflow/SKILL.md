---
name: git-commit-workflow
description: >
  Prepare or create Git commits for the user. Use whenever the user asks to
  commit changes, split work into commits, write commit messages, stage files,
  or clean up a dirty worktree before committing, and after each completed
  phase of a sequenced coding or refactoring workflow. Creates descriptive
  commit messages, logical commit boundaries, and an Assisted-by trailer for
  AI-authored or AI-assisted changes.
---

# Git Commit Workflow

## Goal

Create small, reviewable commits with clear, descriptive subjects and an
explicit AI assistance trailer. Always separate independent logical units of
work into separate commits. Prefer preserving user changes over making a tidy
history at the cost of mixing unrelated work.

## When To Commit

- When the user asks to commit, stage, split commits, or push as part of
  finishing work.
- When a loaded skill or reference completes a phase of sequenced work and
  the tree changed. Treat that as authorization. Commit the phase, then
  continue the next phase without waiting for a separate commit prompt.

## Safety Rules

- Never discard, reset, restore, clean, or overwrite user changes unless the
  user explicitly asks for that destructive action.
- Inspect the worktree before staging: `git status --short`, then targeted
  diffs for changed files.
- Stage only files that belong to the commit being made.
- Treat already-staged changes as provisional, not authoritative. Inspect the
  staged diff and unstage/re-stage by logical unit when staged content contains
  multiple independent concerns.
- If unrelated changes exist, leave them unstaged and mention them.
- If a file contains changes from multiple concerns, use an interactive or
  patch-based staging flow when available. If not practical, ask before mixing
  concerns.
- Do not amend, rebase, squash, or rewrite commits unless the user asks.

## Commit Splitting

Split commits when their changes can be reviewed, validated, and reverted
independently. Keep changes together when splitting would leave either commit
incomplete or broken.

If the staged area or worktree contains multiple independent features, bug
fixes, refactors, documentation changes, dependency updates, generated output,
formatting, or infrastructure changes, create separate commits.

Each commit should represent one coherent unit of work that can be understood,
reviewed, reverted, and tested on its own. Common split points:

- Feature code vs tests or fixtures only when the tests are broad or touch
  shared helpers. Otherwise keep feature and its focused tests together.
- Behavior changes vs formatting, generated files, dependency bumps, or docs.
- Multiple features, bug fixes, or refactors.
- Application code vs infrastructure/config changes.
- Mechanical renames/moves vs semantic edits.

Keep changes together when separating them would leave the repository broken,
unreviewable, or without the tests that explain the behavior.

When multiple logical units are already staged together:

- Do not assume the user's staging is the desired commit boundary.
- Inspect `git diff --cached` and identify logical groups before committing.
- Rebuild the index so each commit contains only one group.
- Ask before mixing groups only when patch-level separation is impractical or
  when the correct boundary is genuinely ambiguous.

## Commit Subject

Write one short, plain-English line that states the concrete behavior or
structural change. A reader should understand what changed without decoding a
type, scope, or local taxonomy.

- Prefer an imperative verb: `Add file organization guidance`, `Keep generated
  wrappers quiet`, or `Split notification parsing from delivery`.
- Name the affected thing when it helps the reader find the change: `Preserve
  failed jobs during retry` is clearer than `Improve retries`.
- Do not use Conventional Commit prefixes or invent broad labels such as
  `chore`, `refactor`, or `cleanup` when they say less than the actual change.
- Keep the subject concise, ideally 72 characters or fewer. Do not end it with
  a period.
- Follow a repository's explicit commit-message convention when it has one. In
  the absence of one, use this format.

## Commit Body

Use a body only when it records a non-obvious reason, risk, migration,
tradeoff, or useful validation evidence. Do not use it to repeat the subject.
Keep it concise:

```text
Add file organization guidance

Explain the reason and any important tradeoffs that the diff does not show.

Validation:
- go test ./...
- npm test

Assisted-by: Codex/GPT-5
```

If the change breaks a public contract, state the affected readers and required
migration plainly before `Assisted-by:`.

## Assisted-by Trailer

For AI-authored or AI-assisted commits, always add:

```text
Assisted-by: <agent>/<model>
```

Use the most specific current agent and model name available. Examples:

- `Assisted-by: Codex/GPT-5`
- `Assisted-by: OpenCode/GPT-5`

If the exact model is unavailable, use the agent name with `AI`, for example:

```text
Assisted-by: Codex/AI
```

Keep this as a Git trailer at the end of the commit message after other
footers.

## Workflow

1. Inspect current status and diffs.
2. Inspect both unstaged and staged diffs; treat staged content as input to
   regroup, not as final commit scope.
3. Group changes into logical commits using the splitting rules above.
4. For each commit, stage only that group.
5. Write a descriptive commit message with an `Assisted-by` trailer.
6. Run targeted validation when practical before committing, especially for
   code changes.
7. Commit with a multi-line message so the trailer is preserved exactly.
8. Repeat staging, validation, and commit for each remaining logical group.
9. Show the resulting commit hashes and any remaining uncommitted changes.
