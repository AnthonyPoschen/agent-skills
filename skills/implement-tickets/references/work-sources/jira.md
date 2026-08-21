# Jira Work Source

The bundled supervisor supports Jira-backed runs with GitLab merge requests.
Use `--source jira --work-item PROJECT-123`; Jira keys are canonical string
work-item IDs, so do not use numeric `--issue` or `--issue-range` selection.

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
item run. The adapter reads native inward `Blocks` links only. If a Story has no
configured child descendants, it fails unless the user has explicitly confirmed
direct Story scope and initialization is rerun with `--confirm-direct-story`. A
completed Jira status alone does not prove integration:
verify the linked pull/merge request commit is present on the configured Git
target branch for every repository listed in the project contract. A leftover
open Jira item after that Git proof does not block dependents.

Jira supplies work state and GitLab supplies review and integration state. Every
GitLab merge request description includes a Jira link. Once the MR is created,
the supervisor posts an ADF Jira comment with a native GitLab MR hyperlink. Only
after GitLab reports the MR merged and its merge commit is an ancestor of the
fetched target branch does the supervisor transition the Jira item to the first
configured `completed_statuses` value.

Before dispatch, verify `acli jira auth status`, pagination, JQL escaping,
parent traversal, issue-link direction, assignments, and transitions in the
configured project. Jira review-comment feedback polling is not implemented;
review feedback continues through the GitLab MR.
