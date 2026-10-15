# Tagging Policy

## Recommended Branch Tag

Use a tag that is both commit-traceable and orderable:

```text
<branch>-<short-git-sha>-<unix-timestamp>
```

Example:

```text
master-5838cfa-1779289200
```

This gives Flux a sortable numeric field while preserving the Git commit in the
selected image tag.

## GitHub Actions Example

Add a preparation step before `docker/metadata-action`:

```yaml
- name: Prepare image tag
  id: image-tag
  shell: bash
  run: |
    short_sha="${GITHUB_SHA::7}"
    timestamp="$(date -u +%s)"
    echo "commit_tag=${GITHUB_REF_NAME}-${short_sha}-${timestamp}" >> "$GITHUB_OUTPUT"
```

Then include the raw tag in Docker metadata:

```yaml
- name: Extract Docker metadata
  id: meta
  uses: docker/metadata-action@v6
  with:
    images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
    tags: |
      type=ref,event=branch
      type=ref,event=tag
      type=sha,prefix=sha-
      type=raw,value=${{ steps.image-tag.outputs.commit_tag }}
      type=raw,value=latest,enable={{is_default_branch}}
```

Keep `type=sha` and `latest` if humans or other tools use them, but do not make
Flux rely on plain SHA tags for newest-image selection.

## Why Not Plain SHA Tags

Tags such as `sha-5838cfa` and `sha-e655b0b` are identifiers. Alphabetical or
numerical ordering of those strings does not map to commit time, branch order,
or deploy freshness.

## Why Not Only Latest

`latest` can work with `digestReflectionPolicy: Always`, but it hides which Git
commit produced the running container unless the digest is traced through
external metadata. It is also easy to confuse humans during rollback or incident
review because many builds share the same tag.

## Workflow Loop Prevention

When Flux commits image updates into `kustomization/**`, GitHub Actions should
ignore those paths if image-update commits should not trigger another build:

```yaml
on:
  push:
    branches:
      - master
    paths-ignore:
      - "kustomization/**"
```

Without this, a Flux image update commit can trigger a new image build, which
can trigger another Flux image update commit.
