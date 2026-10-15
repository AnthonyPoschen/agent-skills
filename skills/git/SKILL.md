---
name: git
description: >
  Git commits, worktrees, and releases. Use when the user asks to commit,
  split work into commits, write a commit message, or stage files, and after
  each completed phase of sequenced work (git-commit-workflow). Use when the
  user asks to work in a worktree, isolate a change on a branch, clone or fork
  a repo, or clean up merged worktrees, or says "make sure its in a worktree"
  or "in a new worktree / branch" (git-worktree). Use when preparing,
  publishing, or correcting a release, changelog, release notes, version, or
  tag (release).
---

# Git

Commits, worktrees, and releases happen together, so they share one skill.
Load only the reference for the current step.

## Routing

- Committing, staging, splitting commits, commit messages, or the
  Assisted-by trailer: read [commit](./references/commit.md).
- Cloning or forking a repo, creating, syncing, or retiring a worktree, or
  isolating a change on a branch: read [worktree](./references/worktree.md).
  Its lifecycle script is `<skill-path>/scripts/worktree-lifecycle.sh`.
  Clone into `~/git/<provider>/<user>/<repo>`; do not
  `gh repo fork --clone` into the current directory.
- Changelogs, release notes, version selection, tags, or publishing: read
  [release](./references/release.md). It loads
  [CalVer](./references/calver.md) when needed.

A release usually ends with a commit and a tag. Follow the release reference
for what to record, then the commit reference for how to commit it.
