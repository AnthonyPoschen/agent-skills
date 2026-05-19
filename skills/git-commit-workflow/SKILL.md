---
name: git-commit-workflow
description: >
  Prepare or create Git commits for the user. Use whenever the user asks to
  commit changes, split work into commits, write commit messages, stage files,
  or clean up a dirty worktree before committing. Enforces Conventional Commits,
  logical commit boundaries, and an Assisted-by trailer for AI-authored or
  AI-assisted changes.
---

# Git Commit Workflow

## Goal

Create small, reviewable commits with accurate Conventional Commit subjects and
an explicit AI assistance trailer. Always separate independent logical units of
work into separate commits. Prefer preserving user changes over making a tidy
history at the cost of mixing unrelated work.

## Safety Rules

- Never discard, reset, restore, clean, or overwrite user changes unless the
  user explicitly asks for that destructive action.
- Never push commits or tags. Do not run `git push`, push aliases, `git
  push --tags`, `git push --follow-tags`, or any command that publishes local
  commits to a remote.
- If the user asks to commit and push, create the local commit only, then stop
  and report that pushing is intentionally left for the user to review and run.
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

Break work down into logical commits by default. If the staged area or worktree
contains multiple features, bug fixes, refactors, docs changes, dependency
updates, generated output, formatting, or infrastructure changes, create
multiple commits instead of one combined commit.

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

## Conventional Commits

Use this subject format:

```text
<type>(<scope>): <summary>
```

Omit the scope when it would be vague.

Preferred types:

- `feat`: user-visible feature or new capability
- `fix`: bug fix or correctness repair
- `refactor`: code restructuring without behavior change
- `test`: test-only changes
- `docs`: documentation-only changes
- `build`: build system, package manager, or dependency changes
- `ci`: CI workflow changes
- `chore`: maintenance that does not fit the above
- `style`: formatting-only changes with no semantic effect
- `perf`: performance improvement
- `revert`: revert a previous commit

Subject rules:

- Use imperative mood: `fix parser panic`, not `fixed parser panic`.
- Keep the subject concise, ideally <= 72 characters.
- Do not end the subject with a period.
- Use `!` for breaking changes, e.g. `feat(api)!: rename token field`.

## Commit Body

Add a body when the "why" is not obvious, the change is risky, or there are
important validation notes. Keep it concise:

```text
<type>(<scope>): <summary>

Explain why this change is needed and call out important tradeoffs.

Validation:
- go test ./...
- npm test

Assisted-by: Codex/GPT-5
```

For breaking changes, include a `BREAKING CHANGE:` footer before
`Assisted-by:`.

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
5. Write a Conventional Commit message with an `Assisted-by` trailer.
6. Run targeted validation when practical before committing, especially for
   code changes.
7. Commit with a multi-line message so the trailer is preserved exactly.
8. Repeat staging, validation, and commit for each remaining logical group.
9. Show the resulting commit hashes and any remaining uncommitted changes.
10. Do not push. Leave all commits local for user review.
