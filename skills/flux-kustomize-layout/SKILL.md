---
name: flux-kustomize-layout
description: >
  Set up or review FluxCD and Kustomize repository layouts. Use when the user
  asks for a FluxCD repo layout, Kubernetes GitOps structure, kustomization
  scaffold, base/overlay setup, or dev/prod environment overlays.
---

# Flux Kustomize Layout

## Goal

Create predictable FluxCD-friendly Kustomize layouts with one shared base and
explicit dev/prod overlays. Keep the repository structure boring, discoverable,
and easy for CI or Flux reconciliation to validate.

## Required Layout

Always scaffold this structure unless the user explicitly asks otherwise:

```text
kustomization/
  base/
    kustomization.yaml
  overlays/
    dev/
      kustomization.yaml
    prod/
      kustomization.yaml
```

## Rules

- Put shared Kubernetes manifests in `kustomization/base/`.
- Include every base manifest from `kustomization/base/kustomization.yaml`.
- Make `dev` and `prod` overlays include `../../base`.
- Put environment-specific patches, generators, images, labels, or namespace
  settings in the matching overlay folder.
- Keep base manifests environment-neutral unless the user asks for a different
  environment model.
- Prefer `apiVersion: kustomize.config.k8s.io/v1beta1` and
  `kind: Kustomization` for Kustomize files.
- Use lowercase file names with clear resource names, such as
  `deployment.yaml`, `service.yaml`, `configmap.yaml`, and `ingress.yaml`.
- Do not put application manifests directly in overlay folders when they should
  be shared through base.

## Minimal Files

Base `kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - deployment.yaml
  - service.yaml
```

Overlay `kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
```

## Workflow

1. Inspect existing Kubernetes, Kustomize, or Flux files before editing.
2. Create missing directories and `kustomization.yaml` files in the required
   layout.
3. Move shared manifests into `kustomization/base/` when safe.
4. Add all shared manifests to the base resources list.
5. Add `../../base` to each overlay resources list.
6. Put dev/prod differences in overlays as patches or Kustomize fields.
7. If Flux `Kustomization` custom resources are requested, keep them separate
   from Kustomize `kustomization.yaml` files and name the distinction clearly.
8. Validate with `kustomize build` or `kubectl kustomize` when available;
   otherwise report that validation was not run.
