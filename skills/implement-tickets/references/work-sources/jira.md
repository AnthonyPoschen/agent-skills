# Jira Work Source (Beta)

Jira is a documented beta adapter target. Do not launch Jira-backed workers yet:
the bundled supervisor fails closed until the Atlassian CLI contract has been
tested against a real project.

Plan to use Atlassian CLI (`acli`) for authentication, JQL selection, work-item
reads, comments, transitions, and assignments. During preflight, discover and
record:

- Atlassian site and Jira project key;
- ready and completed workflow states;
- issue types eligible for implementation;
- dependency link names and their direction;
- development integration linking Jira items to GitHub or GitLab reviews;
- the JQL query or explicit keys selecting the run.

Map Jira keys to normalized supervisor items. Map inward blocker links to
explicit blockers. A completed Jira status alone does not prove integration:
verify the linked pull/merge request commit is present on the configured Git
target branch for every repository listed in the project contract. A leftover
open Jira item after that Git proof does not block dependents.

Jira supplies work and discussion state; GitHub or GitLab remains the review and
Git integration source. Poll both sides and deduplicate by provider object ID and
revision. Treat Jira comments with the same actionable, question,
informational, and scope-change classification used for review comments.

Before promoting this adapter from beta, test `acli` authentication discovery,
pagination, JQL escaping, issue-link direction, comment cursors, assignments,
rate limiting, and linked-development reads in a disposable project. If any
required convention cannot be discovered from project documentation or APIs,
ask all setup questions together before creating state.
