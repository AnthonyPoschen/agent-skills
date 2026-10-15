# Skill Merge Batch 2

## Why

Same goal as [batch 1](skill-merge-batch-1.md): fewer skill listings per turn
without losing triggering. Anthony asked to combine the single-file build
skills and to drop Omarchy work from this repo.

## What Changed

| Change | Reason |
| --- | --- |
| `build-files` = `dockerfile` + `makefile` | Both are narrow rules for one kind of project file and are often edited together. |
| `omarchy-development` deleted | No longer wanted in the shared skill set. |

The parent description keeps "Dockerfile", "Makefile", "make target", and each
old skill name, so prompts that matched the old skills still match.

## Trigger Test

Same method as batch 1: realistic prompts routed against the full listing on
`master` and on this branch, by an agent that sees only the listing.

| # | Prompt | Before (master) | After |
| --- | --- | --- | --- |
| 1 | write a Dockerfile for the go service | dockerfile + coding | build-files + coding |
| 2 | make our image run as non-root and use a scratch final stage | dockerfile | build-files |
| 3 | move the CI build steps into Docker so the workflow stays thin | dockerfile (+ coding) | build-files (+ coding) |
| 4 | add a make target that takes a namespace variable | makefile + coding | build-files + coding |
| 5 | review this Makefile | makefile | build-files |
| 6 | update the README so it tells people to run make deploy ENV=prod | makefile + writing | build-files + writing |
| N1 | my pods aren't rolling after I push a new image to ghcr | flux | flux (no build-files) |
| N2 | commit these changes | git | git (no build-files) |
| N3 | add a dark mode toggle to the settings page | frontend-design + coding | unchanged |
| N4 | fix the Omarchy waybar QML in my checkout | omarchy-development | coding only (skill removed on purpose) |

Result: every Dockerfile and Makefile prompt routes to `build-files` at least
as reliably as before. N4 loses its skill by design.

Limit: one routing pass by an agent, not by each harness; no model API or
agent CLI was available on the box to automate repeated runs. Re-run in a real
harness if `build-files` seems to miss.
