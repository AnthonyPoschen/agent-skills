---
name: create-verification-skill
description: >
  Create a project-local verification skill by reading the codebase. Use when
  the user asks to create a verification skill, a verification process, or a
  drive harness for a repo, or when a runnable app has no scripted way to
  launch, drive, and prove its UI, CLI, or service behavior. Writes or extends
  that project's verify skill and its feature map.
---

# Create a verification skill

Build a project-local skill the next agent can run cold. It launches the real
app, drives a feature the way a user does, and keeps the evidence.

The map contract is [`references/feature-map.md`](references/feature-map.md).
Load it before you seed the map. Load
[`references/feature-map-example/index.md`](references/feature-map-example/index.md)
when you need a filled example.

## 1. Interview the repo

Answer these from the codebase. Ask the user only what the repo cannot show.

- **Surface.** What a person or caller actually touches. Note secondary surfaces.
- **Run.** The repo's own local start command, the ready signal, ports, env, seed data, and auth.
- **Drive.** An existing harness first. Otherwise browser control for a web UI, a PTY or tmux session for a CLI or TUI, HTTP for a service.
- **Observe.** Screenshots, transcripts, response bodies, logs, exit codes, stored rows or files.
- **Isolate.** Whether two instances can run at once. When they cannot, the generated skill refuses to drive a second copy of that shared instance.
- **Home.** The directory where this repo keeps skills. When it has none, use `.agents/skills/`.

When a verification skill already exists, extend it. One skill per app.

When the checkout does not start, fix that or report the exact failure before writing the skill. A skill written against a broken start teaches the wrong steps.

## 2. Write the skill

Write `SKILL.md` in the skill home. Frontmatter is `name` and `description` only. The description names the app, the surface, and when to load it.

Ground every section in the interview. Leave no placeholder.

- **Launch.** The exact start command and the ready signal, plus teardown. A short-lived CLI has no server: build once, then each drive gets its own session.
- **Doctor.** One read-only check that the instance is worth driving: process up, right build, port owned by this run, auth valid when the route needs it.
- **Drive.** Real selectors and commands from this repo. Prefer accessible names, data attributes, prompt strings, and route paths.
- **Evidence.** Where proof goes. The proof standards are the feature-map contract. When a person can see the surface, the look is the `eyes-on` skill. The map still owns the drive.
- **Cleanup.** Stop what this run started. Kill the process you started, not every process with that name. Proof files stay at the named path.
- **Helpers.** Every script the skill ships is executable, and the skill body shows the command.

## 3. Seed the feature map

Create the map index and one file per user-facing feature. Start with the top 3 to 5, taken from routes, commands, menus, or docs. Follow the feature-map contract. A proof that uses one convenient entry point is incomplete when the map lists others.

## 4. Prove it

Run the new skill once: launch, doctor, drive one mapped feature, capture evidence, clean up. After cleanup, the evidence is still at the named path. Fix what fails. Run cleanup after a failed attempt so ports and processes do not stay behind.

A skill you have not executed is a draft. Say that when you could not run this step.

## 5. Hand off

Name the skill path, the feature you drove, the evidence path, and any surface you left unmapped. Later passes use `maintain-verification-skill`.
