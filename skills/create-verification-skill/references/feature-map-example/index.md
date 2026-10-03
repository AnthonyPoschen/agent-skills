# Notes verification map

Specimen of a filled map. Write the project's own commands, handles, and paths.

This directory is the maintained source for verifying the user-facing behavior of Notes. Read the index before driving the app, then use the matching feature file as the recipe.

## Baseline preconditions

- Launch Notes at `http://127.0.0.1:4173` with a disposable data directory.
- Set `NOTES_DATA_DIR=/tmp/notes-verify-$RUN_ID` so concurrent runs do not share state.
- Seed notes titled `Quarterly plan` and `Grocery list`.
- Put `control-notes` and the `notes` CLI on `PATH`.
- Run `control-notes doctor` and require the expected URL, data directory, and build revision.
- Drive only an instance this verification run started.

## Driving conventions

- Start every recipe from the baseline state unless its preconditions say otherwise.
- Prefer ARIA roles and accessible names over CSS selectors or DOM position.
- Treat every command as literal. Keep quoted names and flags unchanged.
- Run browser actions through `control-notes browser`.
- Run terminal actions through `control-notes cli -- <command>`.
- Restore seeded data after a mutation. Proof artifacts stay.

## Proof and skip reporting

Capture the user action and the resulting state. UI proof includes an ARIA snapshot and a screenshot with the app identity visible. CLI proof includes the command, stdout, stderr, and exit code. Mutation proof includes a read-only second view of the stored value. Record the feature id and entry point with every artifact. Report an unreachable path with the command you tried and the unmet precondition. A skipped entry point stays unverified.

## Features

- [Create a note](./create-note.md) covers browser and CLI creation, cancellation, persistence, and cleanup.
