---
name: build-files
description: >
  Create, review, or edit build files. Use for Dockerfiles and container
  builds: containerizing a project, multi-stage builds, scratch vs Alpine
  images, non-root UID/GID, and thin GitHub Actions that run the Dockerfile
  (dockerfile). Use for Makefiles, make targets, make recipes, task runners,
  or docs that tell users to run make, and always when a make target takes
  user-provided variables such as tokens, paths, namespaces, or environments
  (makefile).
---

# Build Files

Load only the reference for the file in front of you.

## Routing

- Dockerfiles, container images, and Docker-based CI: read
  [dockerfile](./references/dockerfile.md). It points to
  [dockerfile-patterns](./references/dockerfile-patterns.md) for concrete
  examples.
- Makefiles, make targets, and docs that run make: read
  [makefile](./references/makefile.md).

A Makefile target that wraps `docker build` needs both.
