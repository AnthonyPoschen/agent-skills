# GitLab Work Source

Use this adapter when work items are GitLab issues. Read the project GitLab
instructions before any `glab` command.

## Discovery

Use `glab` to list open work items and read their complete descriptions,
comments, discussions, merge requests, and CI results. Parse dependencies only
from the explicit `Blocked by` section. Treat a blocker as complete only after
its merge request is merged into the target branch.

## Updates

Use one source branch per issue. Create or update one merge request that closes
the issue. Read the merge request after each write to verify its source branch,
target branch, description, and state.

Use GitLab discussions for threaded feedback. Every agent-written merge request
comment starts with `AI-generated:`. Do not close issues manually, merge merge
requests, or mark a discussion resolved without verified work.

## Required Information

The supervisor records the issue IID, explicit blockers, issue notes, merge
request IID and URL, source and target branch, merge state, pipeline state, and
latest feedback timestamp.
