# GitLab Work Source

Use this adapter when work items are GitLab issues. Read the project GitLab
instructions before any `glab` command.

## Discovery

Use `glab` to list open work items and read their complete descriptions,
comments, discussions, merge requests, and CI results. Parse dependencies only
from the explicit `Blocked by` section. Treat a blocker as complete after its
merge request is merged into the target branch on every repository listed in
the project contract. A leftover open issue does not block dependents. Close
it after verified integration so the GitLab UI catches up.

## Updates

Use one source branch per issue. Create or update one merge request that closes
the issue. Read the merge request after each write to verify its source branch,
target branch, description, and state.

Use GitLab discussions for threaded feedback. Every agent-written merge request
comment starts with `AI-generated:`. Post that reply on the merge request in
the same poll that discovers the question. A chat-only answer does not
complete the loop. Do not merge merge requests or mark a discussion resolved
on the reviewer's behalf. Close an issue only after its merge request is
merged into the target branch and the merge commit is an ancestor. If GitLab
did not auto-close it, close it then. Never close an unintegrated issue.

## Feedback

Record the authenticated `glab` login as the run's human authority. Read award
emoji on merge-request notes and accept only unambiguous positive awards from
that login as authorization for the exact third-party note. An affirmative
reply authorizes preceding feedback in the same discussion; an individual-note
authorization must identify the source note or `@author`. Accept failed
pipelines and actionable recognized CI/scanner bot notes. Ignore other bot
notes, agent-written `AI-generated:` notes, and unendorsed third-party feedback.
When a human reply grants authority, send both note bodies to the worker.

## Required Information

The supervisor records the authenticated `glab` login, issue IID, explicit blockers, issue notes, merge
request IID and URL, source and target branch, merge state, pipeline state, and
latest feedback timestamp.
