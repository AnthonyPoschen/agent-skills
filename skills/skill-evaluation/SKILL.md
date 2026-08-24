---
name: skill-evaluation
description: >
  Use when adding, changing, comparing, or troubleshooting an agent skill and
  you need evidence that it changes agent behaviour. Evaluate triggering and
  task outcomes with realistic prompts, a shared rubric, and blinded review.
---

# Skill Evaluation

Use this skill to decide whether a skill helps on real work. Do not judge a
skill by how convincing its instructions sound or by an agent saying it followed
them.

Use this skill for a new public skill, a meaningful change to an existing skill,
or a skill that seems to trigger at the wrong time or produce poor work. Skip it
for a small wording change unless the change has meaningful uncertainty or risk.

## State The Question

Start with one question the evaluation can answer.

Good questions:

- Does this skill trigger for the writing tasks it claims to handle?
- Does the new version produce clearer release notes than the old version?
- Does this refactoring skill reduce unnecessary abstractions on a real task?

Do not combine several unrelated questions in one evaluation.

## Evaluate Triggering Separately

Test the skill description with a small set of realistic prompts. Include both
prompts that should trigger the skill and prompts that should not.

For each prompt, record:

- whether the skill should trigger;
- whether it did trigger;
- why the result was correct or incorrect.

Use substantive tasks. A trivial prompt may not activate a skill even when its
description is correct, because the agent can solve it without specialised
guidance.

## Evaluate Behaviour On Real Tasks

Compare candidates on the same task. A candidate can be the new skill, the old
skill, a no skill baseline, or another proposed version.

1. Pick a small number of realistic tasks. Prefer prior user requests, real
   repository files, or representative changes over artificial prompts.
2. Write one short shared rubric before running candidates. The rubric must
   describe the outcome that matters, not the wording of the skill.
3. Give every candidate the same task, inputs, constraints, and expected output
   location. Use isolated workspaces when candidates write files.
4. Save the produced artifacts and direct verification evidence.
5. Hide candidate identities before judging the artifacts. Read the work and its
   evidence, not an agent's summary of what it intended to do.
6. Compare quality, correctness, scope, cost, and unnecessary work.

For subjective work such as writing or design, include human judgment in the
rubric. Do not pretend that word count or an automatic score can decide whether
the work is good.

## Write A Useful Rubric

Keep the rubric short and specific to the task. Use observable criteria.

For a technical writing skill, a rubric could check:

- facts, paths, and symbols are correct;
- the document has a clear purpose and structure;
- language is plain and direct;
- the text does not invent claims or add irrelevant detail.

For a coding skill, a rubric could check:

- the requested outcome works through the appropriate path;
- the change stays in scope;
- abstractions and tests follow the repository rules;
- verification directly supports the claimed outcome.

Avoid a rubric that rewards the skill for repeating its own instructions.

## Judge The Result

Compare candidates on more than the final text or diff:

- **Outcome:** did the work solve the real task?
- **Evidence:** did the candidate verify its claim directly?
- **Scope:** did it avoid unrelated edits, scaffolding, and ceremony?
- **Clarity:** is the result easier for a maintainer or user to understand?
- **Cost:** did it take noticeably more time, tokens, files, or steps for no
  meaningful improvement?

Make one decision: keep the change, revise it, remove it, or gather more
evidence. Do not treat a close or inconclusive result as a win.

## Keep Only Durable Evidence

Keep temporary workspaces, raw transcripts, and scratch artifacts out of normal
commits. Commit a prompt and rubric only when they become a useful regression
case for a public skill.

When reporting results, state:

- the question;
- the tasks and rubric;
- the candidates compared;
- the observed result and evidence;
- the decision and next step.

## Source

This workflow is informed by pstack's
[evaluation playbook](https://github.com/cursor/plugins/blob/main/pstack/skills/poteto-mode/playbooks/eval.md),
adapted to remain tool neutral and focused on small, useful comparisons.
