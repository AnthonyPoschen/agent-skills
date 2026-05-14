# Dockerfile Patterns

Use these as starting points, then adapt commands, versions, ports, and artifact
names to the project.

## Static Binary With Scratch Final Image

Good for Go, Rust, Zig, or C/C++ projects that can produce one static Linux
binary.

```dockerfile
# syntax=docker/dockerfile:1.7

ARG UID=10001
ARG GID=10001

FROM golang:1.24-alpine AS build
ARG UID
ARG GID
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go test ./...
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app ./cmd/app

FROM scratch AS runtime
ARG UID=10001
ARG GID=10001
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt
COPY --chown=${UID}:${GID} --from=build /out/app /app/app
USER ${UID}:${GID}
ENTRYPOINT ["/app/app"]
```

Notes:
- Omit CA certificates only when the binary never performs outbound TLS.
- Do not copy shells, package managers, source trees, or test files into
  `scratch`.
- If CGO or dynamic linking is required, do not use `scratch` unless all needed
  runtime libraries are intentionally copied and verified.

## Python With Alpine Final Image

Good for Python services that need an interpreter and installed dependencies.

```dockerfile
# syntax=docker/dockerfile:1.7

ARG PYTHON_VERSION=3.13
ARG UID=10001
ARG GID=10001

FROM python:${PYTHON_VERSION}-alpine AS deps
ARG UID
ARG GID
WORKDIR /src
RUN apk add --no-cache build-base
COPY requirements.txt ./
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:${PATH}"
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install -r requirements.txt

FROM deps AS test
COPY . .
RUN pytest

FROM python:${PYTHON_VERSION}-alpine AS runtime
ARG UID=10001
ARG GID=10001
WORKDIR /app
RUN addgroup -S -g "${GID}" app \
    && adduser -S -D -H -u "${UID}" -G app app
ENV PATH="/opt/venv/bin:${PATH}" \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1
COPY --chown=${UID}:${GID} --from=test /opt/venv /opt/venv
COPY --chown=${UID}:${GID} --from=test /src/app /app/app
USER ${UID}:${GID}
CMD ["python", "-m", "app"]
```

Notes:
- Replace `pytest` with the project's actual test command.
- Adapt dependency installation for `pyproject.toml`, Poetry, uv, or PDM when
  those tools are the project's source of truth.
- Add runtime Alpine packages explicitly when Python wheels need shared
  libraries at runtime.
- Prefer lock files and deterministic install commands when available.

## Plain Alpine Runtime

Use plain Alpine for dynamic binaries or small runtime bundles that need libc
compatibility, certificates, timezone data, or other OS files.

```dockerfile
FROM alpine:3.22 AS runtime
ARG UID=10001
ARG GID=10001
WORKDIR /app
RUN apk add --no-cache ca-certificates \
    && addgroup -S -g "${GID}" app \
    && adduser -S -D -H -u "${UID}" -G app app
COPY --chown=${UID}:${GID} --from=build /out/app /app/app
USER ${UID}:${GID}
ENTRYPOINT ["/app/app"]
```

## CI Workflow Shape

Keep workflow logic thin and let the Dockerfile own build/test behavior.

```yaml
name: ci

on:
  pull_request:
  push:
    branches: [main]

jobs:
  docker-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/build-push-action@v6
        with:
          context: .
          push: false
          tags: app:ci
```

If tests live behind a dedicated target, build that target in CI:

```yaml
with:
  context: .
  target: test
  push: false
```
