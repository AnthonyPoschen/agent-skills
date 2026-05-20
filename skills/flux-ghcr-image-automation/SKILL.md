---
name: flux-ghcr-image-automation
description: >
  Build, review, debug, or fix GitHub Actions to GHCR to FluxCD image
  automation loops. Use whenever the user mentions Flux image automation,
  ImageRepository, ImagePolicy, ImageUpdateAutomation, GHCR, Docker metadata
  tags, Kubernetes pods not rolling after an image push, mutable latest tags,
  digest updates, or GitHub Actions publishing container images for Flux-managed
  Kubernetes workloads.
---

# Flux GHCR Image Automation

## Goal

Create a reliable GitOps image release loop:

1. GitHub Actions builds and pushes a container image to GHCR.
2. Flux image-reflector-controller discovers the new tag or digest.
3. Flux image-automation-controller commits the selected image reference back to
   Git.
4. Flux kustomize-controller applies that commit and rolls the Kubernetes
   workload.

Prefer traceable, orderable image tags over mutable-only `latest` deployments.
The operator should be able to answer "which Git commit produced this running
pod?" from the image tag.

## When Starting

Inspect the repo and cluster shape before editing:

- `.github/workflows/*.yml`
- `Dockerfile`
- `kustomization/**`
- any `ImageRepository`, `ImagePolicy`, `ImageUpdateAutomation`, `Deployment`,
  and Flux `Kustomization` resources
- the current branch name and default branch
- whether the Flux app `Kustomization` uses `targetNamespace`

If a live cluster is available and the user is debugging a rollout, inspect
live state before changing files:

```sh
flux get images all -n <image-resource-namespace>
kubectl get kustomization -A -o wide
kubectl describe imageupdateautomation <name> -n <namespace>
kubectl get deploy <name> -n <namespace> -o yaml
kubectl get gitrepository -A -o wide
```

## Design Rules

- Do not rely on plain `sha-*` tags for "latest commit" selection. SHA strings
  identify commits but do not sort by time.
- Prefer branch, short SHA, and timestamp tags for branch builds:

```text
<branch>-<short-git-sha>-<unix-timestamp>
```

- Configure Flux `ImagePolicy` to filter the branch tag pattern, extract the
  timestamp, and use numerical ordering.
- Keep Flux image update commits from triggering fresh image builds by ignoring
  the manifest paths that Flux mutates, commonly `kustomization/**`.
- Keep the inline Flux setter comment on every image field Flux should update.
- Make the image setter comment reference the namespace where the `ImagePolicy`
  actually exists after Flux applies the manifests.
- When a Flux `Kustomization` has `targetNamespace`, remember that namespaced
  resources without an explicit namespace may land in that target namespace.
- Set `ImageUpdateAutomation.spec.sourceRef.namespace` when the referenced
  `GitRepository` is in a different namespace.
- Keep `ImageUpdateAutomation.update.path` pointed at the directory that
  contains the file with the image setter comment.

## References

Load only the reference needed for the task:

- [tagging-policy.md](references/tagging-policy.md) for GitHub Actions Docker
  metadata tags and orderable Flux policies.
- [flux-manifests.md](references/flux-manifests.md) for ImageRepository,
  ImagePolicy, ImageUpdateAutomation, setter comments, and namespace rules.
- [diagnostics.md](references/diagnostics.md) for live cluster debugging and
  rollout verification.

## Implementation Workflow

1. Identify the workload image and the registry repository.
2. Confirm which Git branch should drive deployments.
3. Ensure the GitHub Actions workflow publishes an orderable, commit-traceable
   tag for that branch.
4. Ensure `ImageRepository` scans the correct GHCR image.
5. Ensure `ImagePolicy` selects the intended branch builds.
6. Ensure the workload image field has a correct Flux setter comment.
7. Ensure `ImageUpdateAutomation` can read and push to the correct
   `GitRepository`.
8. Validate the rendered manifests.
9. If a cluster is available, use server-side dry-run and live Flux status.

## Validation

Run the strongest practical checks:

```sh
kubectl kustomize <overlay>
kubectl apply --dry-run=server -k <overlay>
git diff --check
```

For live rollout debugging, confirm each stage:

```sh
flux get images all -n <image-resource-namespace>
kubectl describe imageupdateautomation <name> -n <namespace>
kubectl get gitrepository <name> -n <namespace> -o wide
kubectl get kustomization <name> -n <namespace> -o wide
kubectl get deploy <name> -n <namespace> -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
kubectl get pods -n <namespace> -l <selector> -o wide
```

Report exactly where the loop is stopped: image build, registry scan, image
policy selection, Git update commit, Kustomization apply, or Deployment rollout.
