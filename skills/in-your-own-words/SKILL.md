---
name: in-your-own-words
description: >
  Use before starting ambiguous or multi-step work, when the user says "in your
  own words" or "restate my goal", or when intent is unclear. Restate the goal,
  the problem being solved, the constraints, and what done looks like in plain
  words; flag assumptions and open questions; wait for a correction if anything
  looks off.
---

# In Your Own Words

Before non-trivial work, play back the request in your own plain words so a
misread goal is caught before the run starts. A short play-back costs one reply;
a wrong goal can waste the whole session.

## When To Play Back

Do this when:

- The task is ambiguous, multi-step, or spans more than one decision.
- The user asks for `/in-your-own-words`, says "in your own words", or "restate
  my goal".
- You are about to spend meaningful effort and the done check is not obvious.

Skip it for a clear, one-step request with an obvious done check.

## Play Back

Read only enough of the message, linked files, and conversation to understand
the ask. Do not investigate causes, design a fix, or start the work yet.

Restate in your own words:

1. **Goal** — what success looks like, and for whom.
2. **Problem** — what is wrong or missing today.
3. **Constraints** — limits, must-nots, deadlines, and authority boundaries.
4. **Done** — the concrete check that would pass or fail.

Then briefly flag:

- **Assumptions** — guesses you filled in; mark each as a guess.
- **Open questions** — at most one question, only when the answer changes the
  work. Say what each answer would change. Never ask the user to restate their
  request.

Keep exact names and values (paths, IDs, amounts, dates, error strings). Do not
parrot the user's sentences. Do not propose a plan or fix. Stay under about 120
words. Leave out any section that has nothing useful in it.

If anything looks off, stop and wait for a correction. If the play-back is
clear and the user only confirms, proceed as described. Do not play it back
again unless they ask.

## Source

Inspired by poteto's `/in-your-own-words` prompt and pstack skill; first-party
wording owned here.
