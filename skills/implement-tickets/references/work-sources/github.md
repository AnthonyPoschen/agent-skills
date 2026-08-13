# GitHub Work Source

Use this adapter when work items are GitHub Issues and reviews are GitHub pull
requests. Read the project's GitHub conventions before calling `gh`.

## Preflight

Verify `gh auth status`, repository identity, default/target branch, current
login, branch protection expectations, and issue/PR permissions. Infer the
repository from `git remote get-url`; pass `--repo OWNER/REPO` on manager calls
so commands never depend on a worker's current directory.

## Discovery And Dependencies

Use GitHub's native issue dependencies when available. The GraphQL `Issue`
object exposes `blockedBy`; query it with each candidate issue rather than
parsing the visual issue page. Fall back to an explicit `Blocked by` section
only when the repository does not use native relationships.

Query full issue bodies, labels, assignees, comments, state, and blocker states.
The normal ready predicate is:

```text
OPEN + ready label + unassigned + no open/unintegrated blocker
```

A run-owned assignment may resume its matching branch/worktree. Refuse another
assignee unless the user explicitly transfers ownership.

Do not trust `CLOSED` alone as proof of integration. Prefer a linked closing PR
whose merged commit is an ancestor of the fetched target branch. When a project
uses manual closure, require equivalent Git evidence or a recorded user
override.

## Claim, Publish, And Verify

Claim with `gh issue edit <number> --add-assignee @me` before dispatch. Workers
do not call GitHub.

The manager pushes with an explicit refspec:

```sh
git push origin HEAD:refs/heads/issue/<number>-<slug>
```

Create or update a draft PR with an explicit base and head. Include
`Closes #<number>` only when every required repository change will be integrated
by that PR; cross-repository work normally links the canonical issue without
premature closure.

After creation or update, read the PR back with `gh pr view --json` and verify:

- `headRefName`, `headRefOid`, and `baseRefName`;
- `isDraft`, `mergeable`, and `mergeStateStatus`;
- review decision, reviews/comments, and `statusCheckRollup`;
- `mergedAt` and `mergeCommit` after human integration.

`gh pr checks --json` provides machine-readable pass, fail, pending, skipped,
and cancelled buckets. A pending check is waiting, not failure.

## Feedback

Use flat PR comments for general feedback and GraphQL review threads when inline
context or resolution state matters. Read reactions on PR and inline comments
so the authenticated `gh` login can authorize a third party's exact comment.
An affirmative reply authorizes preceding feedback in the same review thread;
an unthreaded authorization must identify the source comment URL or `@author`.
Store stable comment/review/thread IDs and timestamps so polling is idempotent.
Accept provider-native failed checks and actionable output from recognized CI,
coverage, security, or scanner apps. Ignore other bot comments and unendorsed
third-party feedback.

Every agent-written PR comment starts with `AI-generated:`. Keep comments
concise: changed location first, then short action bullets. Do not resolve a
human thread or approve the manager's own PR.

## Human Merge Boundary

Never run `gh pr merge`, enable auto-merge, or mark a draft ready unless the
user explicitly changes the workflow. A human merge is the dependency gate.
After notification of a merge, fetch the target, verify the merge commit is an
ancestor, verify the issue state, safely prune only authorized worktrees, and
recompute readiness.

Treat a human-closed, unmerged PR as terminal `review_closed`; do not relaunch
its branch. It is not integration success. When all selected items are
terminal, run the guarded final cleanup and exit the supervisor.

## Supervisor Fields

Persist repository name, authenticated `gh` login, issue number/node ID, blockers, assignees, labels,
issue update cursor, source/target branches, PR number/URL/head OID, merge state,
check summary, latest human feedback IDs, and merged commit OID.
