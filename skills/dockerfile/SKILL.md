---
name: dockerfile
description: >
  Create, review, or update Dockerfiles and container build workflows. Use for
  containerizing projects, designing multi-stage builds, choosing scratch vs
  Alpine final images, enforcing non-root UID/GID ownership, running build/test
  steps in Docker, and keeping GitHub Actions workflows thin by making the
  Dockerfile the local CI source of truth.
---

# Dockerfile

## Goal

Create Dockerfiles that are secure, reproducible, and useful as the canonical
local CI build. Prefer multi-stage builds unless the user explicitly asks for a
single-stage file.

## Decision Order

1. Identify the final runtime artifact.
2. If the output is a single static binary, aim for a `scratch` final image.
3. If the output needs an interpreter, shared libraries, package files, or
   runtime dependencies, aim for an Alpine final image such as
   `python:<version>-alpine`, `node:<version>-alpine`, or plain `alpine`.
4. If Alpine is incompatible with required native dependencies, explain why and
   choose the smallest appropriate runtime base.
5. Keep build tools, test tools, package managers, and source code out of the
   final image unless they are required at runtime.

## Required Container Policy

- Add `ARG UID=10001` and `ARG GID=10001` in every stage that creates or owns
  files used by the runtime image.
- Run the final container as non-root with `USER ${UID}:${GID}` or a named user
  created from those args.
- Ensure every file used by the final executable is owned by the runtime UID/GID.
  Prefer `COPY --chown=${UID}:${GID}` for final-stage copies.
- For `scratch`, use numeric `USER ${UID}:${GID}` and copy only required files:
  the static binary, CA certificates when outbound TLS is needed, and any
  minimal config/data files the binary actually reads.
- For Alpine runtime images, create the user/group from `UID` and `GID`, then
  copy application files with that owner.
- Set a narrow `WORKDIR`; do not run from filesystem root.
- Use exec-form `ENTRYPOINT` or `CMD`.

## Build And Test Pattern

- Use a dedicated build stage for compilation, dependency installation, and
  generated artifacts.
- Add a test stage when the project has a clear test command. Examples:
  `go test ./...`, `cargo test --locked`, `npm test`, `pytest`.
- Make the default Docker build exercise tests when practical by having the
  production stage depend on the tested artifact path, or document:
  `docker build --target test .`.
- Keep the Dockerfile sufficient for local CI verification so GitHub Actions can
  usually be a simple `docker build` plus optional image push.
- Use BuildKit cache mounts when they improve repeat builds and the project
  already assumes BuildKit.

## Workflow

1. Inspect the project manifest and lock files before writing the Dockerfile.
2. Determine whether the runtime is a static binary, dynamic binary, or dynamic
   language/runtime app.
3. Choose the final image using the decision order above.
4. Write a `.dockerignore` if missing or incomplete.
5. Include test and build stages that match the project tooling.
6. Ensure final-stage ownership, non-root user, and minimal copied files.
7. If asked for CI, make the workflow invoke the Dockerfile rather than
   duplicating build logic.
8. Validate with `docker build` when Docker is available; otherwise report that
   validation was not run.

## Pattern Reference

Read [references/patterns.md](references/patterns.md) when writing concrete
Dockerfile examples or adapting this policy to Go, Rust, Python, Node, or a
plain Alpine runtime.
