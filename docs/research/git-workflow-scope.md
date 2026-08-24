# Git workflow scope

## Conclusion

Rename `git-commit-workflow` to `git-workflow` and broaden it from staging and
committing to the ordinary change path: inspect state, choose a base, create an
isolated branch or worktree when it earns its place, make coherent commits,
push an explicit branch, and prepare a review. Do not turn it into a general
delivery orchestrator.

This gives ordinary requests such as “start a branch”, “open a PR”, or “push
this change” one clear entry point. It does not replace `implement-tickets` or
`release`.

## What Pstack contributes

Pstack's [opening-a-PR playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/opening-a-pr.md)
has four useful ideas:

- Work from an isolated checkout when work could collide with other work.
- Shape ordered commits so each is landable and tells part of the change story.
- Make a review description explain intent, scope, tradeoffs when real, impact,
  and direct verification evidence.
- Keep independent work based on the integration branch. Use a stack only when
  changes genuinely depend on one another.

Its [Babysit playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/babysit.md)
also has sound review rules: only one process should drive a stack, address the
lowest unmerged change first, batch known fixes before restarting checks, and
do not merge without explicit authority. The [shipping playbook](https://github.com/cursor/plugins/blob/46125561306434d8a1d7745d540d8932ab0cd2a2/pstack/skills/poteto-mode/playbooks/shipping.md)
separates green checks from an independent verification decision before merge.

Pstack does **not** prescribe a general branch naming convention. Its worktree
guidance is about isolation, not naming. Its Graphite rules are specifically for
stacked pull requests, not ordinary Git branches.

## Recommended `git-workflow` boundary

The renamed skill should cover one ordinary change through a review-ready
branch:

1. Inspect the worktree, remotes, current branch, and repository instructions.
2. Preserve unrelated work. Identify the target branch and fetch it without
   changing the primary checkout.
3. Create a branch from that target when the user asks for one or a review is
   intended. Follow a repository naming convention; otherwise use a short,
   descriptive task name. Do not invent a global taxonomy.
4. Use a separate worktree when parallel work, a dirty primary checkout, or a
   long-running change makes isolation useful. Do not require one for every
   small local edit.
5. Keep the current selective staging, coherent commit, descriptive subject,
   validation, and assistance-trailer rules.
6. Push only the intended branch with an explicit refspec. Create or update a
   review only with authority, with an explicit base and a concise description
   of purpose, scope, meaningful tradeoffs, impact, and outcome proof.
7. Never force-push, rebase, amend, merge, enable auto-merge, delete branches,
   or delete worktrees without explicit authority. Report the branch, commit,
   validation, and review URL.

The review description wording belongs with `technical-writing`; the Git skill
only specifies which facts a review must carry.

## Keep the other skills separate

| Concern | Existing or proposed owner | Reason |
| --- | --- | --- |
| One change, branch, commits, push, and review preparation | `git-workflow` | This is one coherent everyday Git path. |
| Many dependency-linked tickets, isolated workers, review polling, and human merge | `implement-tickets` | It has a distinct scheduler and authority model. |
| Release records, version selection, and tags | `release` | Release work is not a by-product of commit syntax. |
| Repeated PR monitoring and merge readiness | Later small `pr-follow-through` skill, if needed | Pstack's Babysit ideas are useful but would overload ordinary commits. |

## Do not import

- Do not require Graphite, stacked pull requests, or a fresh worktree for every
  edit. They add process where direct work is simpler.
- Do not adopt Pstack's Conventional Commit requirement. This repository has
  deliberately chosen descriptive subjects instead.
- Do not make CI green alone the definition of a safe merge. Keep direct
  outcome proof and human merge authority.
- Do not combine ticket scheduling, releases, cleanup, and PR monitoring into
  `git-workflow`.

## Local follow-up

The current `implement-tickets` skill and worker assets still ask for
Conventional Commit subjects. If the descriptive subject decision applies to
managed tickets too, update those references and the orchestrator validation in
the same change. Leaving them unchanged creates two conflicting commit rules.
