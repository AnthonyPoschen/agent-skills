# Jira Work Source (Beta)

Jira selection and contract shape are documented, but the bundled OpenCode
supervisor remains fail-closed until its normalized item identifiers support
Jira keys and it can poll Jira hierarchy and links without falling through to
GitHub discovery. Do not launch Jira-backed workers through the supervisor yet.

Use Atlassian CLI (`acli`) for authentication, JQL selection, work-item reads,
comments, transitions, assignments, hierarchy, and links. During preflight,
discover and record:

- Atlassian site and Jira project key;
- ready and completed workflow states;
- container and executable child issue types;
- dependency link names and their direction;
- development integration linking Jira items to GitHub or GitLab reviews;
- the explicit container or Sub-task keys selecting the run.

For a selected Story or higher container, recursively list child work until the
configured executable child type is reached. A selected Sub-task is the single
item run. Map inward native blocker links to explicit blockers. A completed Jira
status alone does not prove integration:
verify the linked pull/merge request commit is present on the configured Git
target branch for every repository listed in the project contract. A leftover
open Jira item after that Git proof does not block dependents.

Jira supplies work and discussion state; GitHub or GitLab remains the review and
Git integration source. Poll both sides and deduplicate by provider object ID and
revision. Treat Jira comments with the same actionable, question,
informational, and scope-change classification used for review comments.

Before dispatch, verify `acli` authentication discovery, pagination, JQL
escaping, parent traversal, issue-link direction, comment cursors, assignments,
rate limiting, and linked-development reads in the configured project. If any
required convention cannot be discovered from project documentation or APIs,
ask all setup questions together before creating state.
