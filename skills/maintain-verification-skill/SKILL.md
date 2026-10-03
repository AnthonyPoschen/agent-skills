---
name: maintain-verification-skill
description: >
  Keep a project's verification skill and feature map honest. Use when the
  user asks to maintain, audit, or update the verify skill, or when that
  map may have drifted from the app. Cover every feature from source and
  with one live drive. Edit only that skill directory.
---

# Maintain a verification skill

A feature map rots when the app changes. Cover every feature file from source and exercise every feature live. The file shape and proof standards are the [feature-map contract](../create-verification-skill/references/feature-map.md).

## Outcomes

End on one, and say which:

- **clean.** Every feature got source and live coverage, and nothing needs shipping.
- **changed.** One set of proven corrections to the verification skill.
- **blocked.** Coverage could not finish, or a proven fix could not be applied safely. Name the blocker.

## Edit scope

Edit only the verification skill directory: its `SKILL.md`, `features/`, and harness scripts it owns. When the app no longer does what the map says, fix the map if the map is wrong. Report a product regression and leave the product code unchanged.

## Pass

0. **Locate.** Find the project verification skill: launch and drive sections, plus a feature map. Several candidates: ask which one. None: stop and load `create-verification-skill`.

1. **Index.** Read the feature index and its sibling feature files. Fix missing, extra, duplicate, or dead entries. Done when the index and the files name the same set.

2. **Source wave.** Launch one read-only subagent per feature file, together. Each returns a feature summary, the source entry points, likely drift with citations or none, and one live recipe. Children do not drive the app and do not edit files. Done when every feature file has a summary.

3. **Reconcile.** Merge overlapping recipes into as few app states as practical. Spot-check cited drift. Add a missing user-facing surface only when you can name its source path.

4. **Live pass.** Drive even when the source looks clean. Follow that skill's Launch section: one long-lived instance driven in series, or a fresh isolated session per drive for a short-lived CLI.

   Hold these for the whole pass:

   - Doctor before the first drive, on each fresh session, and after any failed drive. When doctor cannot see a wedged UI on a healthy process, reset or relaunch.
   - Evidence captured so far is still at its named path after every cleanup. Check the path.
   - Clean residue a drive started when that drive is done. On a shared instance, clean the residue and leave the instance.

   A doctor failure caused by the skill is drift: fix it inside the edit scope and retry once. Restart only what that fix invalidated. A feature you cannot reach is `verified-unreachable` when you name the prerequisite and the route you tried. A prerequisite the map omits is drift. Re-drive a harness fix before it counts.

   Tear down after the last drive, including those re-proofs. Evidence stays.

5. **Triage.** A wrong user-facing description is a map fix. Behavior the harness cannot drive is a harness fix: the script is executable and the skill body shows the command. Broken app behavior is a note for the user, outside this change.

6. **Ship or stop.** For **changed**, re-read every changed file first. The change is the verification skill only. Open a pull request only when the user asked. For **clean** or **blocked**, leave the branch alone. Report the coverage.

Keep run notes in a scratch file. Leave them uncommitted.
