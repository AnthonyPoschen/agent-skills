---
name: flux
description: >
  FluxCD GitOps work. Use for GitHub Actions to GHCR to Flux image automation:
  ImageRepository, ImagePolicy, ImageUpdateAutomation, Docker metadata tags,
  mutable latest tags, digest updates, or Kubernetes pods not rolling after an
  image push (flux-ghcr-image-automation). Use for FluxCD and Kustomize repo
  layouts: GitOps structure, kustomization scaffolds, base/overlay setup, or
  dev/prod environment overlays (flux-kustomize-layout).
---

# Flux

Load only the reference for the current task.

## Routing

- Image builds published to GHCR, Flux image policies and automation, tag
  strategy, or rollouts that do not happen: read
  [image-automation](./references/image-automation.md). It points to
  [tagging-policy](./references/tagging-policy.md),
  [flux-manifests](./references/flux-manifests.md), and
  [diagnostics](./references/diagnostics.md).
- Repository layout, kustomization scaffolds, bases, and environment
  overlays: read [kustomize-layout](./references/kustomize-layout.md).

An image automation change often also touches the layout. Load both when a
task spans them.
