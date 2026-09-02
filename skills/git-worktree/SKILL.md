---
name: git-worktree
description: Create, synchronize, and safely retire isolated Git worktrees for implementation work. Use whenever the user asks to work in a worktree, isolate a change on a branch, create a feature checkout, update a worktree after branch changes, or clean up merged worktrees. New worktrees mirror a source checkout below ~/git at ~/worktree/<provider>/<user>/<repo>/<branch-name>.
---

# Git Worktree

## Goal

Keep implementation changes isolated from the primary checkout while making
every checkout easy to find. A source repository at:

```text
~/git/<provider>/<user>/<repo>
```

has branch worktrees at:

```text
~/worktree/<provider>/<user>/<repo>/<branch-name>
```

Keep slashes in a branch name as path separators. For example,
`feature/add-search` belongs at
`~/worktree/github/acme/widgets/feature/add-search`.

## Workflow

1. Confirm the current directory belongs to a Git working tree with
   `git rev-parse --show-toplevel`. Resolve its canonical path.
2. Find the primary checkout with `git worktree list --porcelain`. Its path
   must be below `~/git/`; otherwise stop and ask for a source checkout or an
   explicitly approved path. Do not guess a different storage layout.
3. Choose the branch named by the user. If none is named, derive a concise,
   descriptive branch name from the requested work and state it before
   creating the worktree. Validate it with `git check-ref-format --branch`.
4. Derive the relative source path below `~/git/`, then construct the exact
   target path below `~/worktree/` by appending the branch name. Do not shorten,
   flatten, hash, or otherwise change either part of the path.
5. If already inside a linked worktree, check that its path equals the derived
   target and that it is on the intended branch. Continue there only when both
   match. If it is an older worktree at another path, leave it untouched and
   ask before changing location.
6. Inspect `git worktree list --porcelain` before creating anything. If the
   target already exists and is registered, resume it. If the branch is already
   checked out elsewhere, do not force it into another worktree; resume that
   checkout when it matches the required target or ask for a different branch.
7. Create the target's parent directories, then create the worktree from the
   source checkout: use `git worktree add -b <branch> <target>` for a new
   branch, or `git worktree add <target> <branch>` for an existing unoccupied
   branch. Base a new branch on the source checkout's current `HEAD` unless the
   user specifies another base.
8. Change into the target and verify `git rev-parse --show-toplevel`,
   `git branch --show-current`, and `git status --short`. Perform all requested
   edits, tests, and commits from this worktree. Never edit the primary checkout
   as part of this workflow.

## Manage the lifecycle

Treat a worktree as the checkout for one branch for its full lifetime.

1. Create or resume the worktree before editing it.
2. When you learn that its branch changed remotely, was rebased, or received
   review changes, synchronize it before doing more work. From the primary
   checkout, run this Bash helper:

   ```text
   <skill-path>/scripts/worktree-lifecycle.sh sync <branch-name>
   ```

   The command fetches the branch's remote, then fast-forwards only a clean
   worktree. If the worktree has local changes or diverged history, leave it
   intact and report the condition instead of stashing, resetting, rebasing, or
   forcing an update.
3. When a branch is known to be merged, run the guarded cleanup immediately
   from the primary checkout:

   ```text
   <skill-path>/scripts/worktree-lifecycle.sh cleanup-merged --apply
   ```

   It fetches the default remote branch, removes only clean mirrored worktrees
   whose branches are ancestors of that branch, and then deletes their merged
   local branches. It never deletes remote branches.
4. When the user abandons an unmerged branch, remove its clean checkout from
   the primary checkout with
   `worktree-lifecycle.sh remove <branch-name>`. Keep the branch unless the
   user explicitly asks to delete it.
5. After a worktree directory was manually removed or a filesystem restore,
   inspect with `git worktree list` and run `git worktree prune` only when the
   stale entries are understood.

Run `worktree-lifecycle.sh cleanup-merged` without `--apply` to preview every
candidate. Report skipped dirty, detached, non-mirrored, and unmerged worktrees
so the user can decide their fate.

## Automate local merge cleanup

When the user asks for automatic cleanup, run
`worktree-lifecycle.sh install-post-merge-hook` from the primary checkout. The
opt-in hook runs the same guarded merged-branch cleanup after a local merge or
pull. It refuses to overwrite another hook. It needs Bash and Git.

Hosted Git providers cannot run a command on this machine when a pull request is
merged remotely. When an agent learns about that merge, it must run the cleanup
step above; the fetch in that command makes remote merges visible locally.

## Safety

- Preserve existing worktrees, branches, and uncommitted changes. Do not use
  `--force`, remove a worktree, prune metadata, reset, or clean as setup.
- Do not create a worktree for the primary branch when it is already checked
  out in the primary checkout. Choose a task branch instead.
- Report the branch and absolute worktree path before beginning the requested
  implementation. Cleanup after a known merge is part of this lifecycle; all
  other cleanup requires an explicit user request.
