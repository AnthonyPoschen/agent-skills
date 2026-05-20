# Diagnostics

## Debug The Loop In Order

Check each stage in sequence. Do not jump straight to restarting pods.

## 1. GitHub Actions Build

Confirm the workflow ran for the commit and pushed the expected tags to GHCR.
If using GitHub CLI:

```sh
gh run list --workflow "Publish Docker image" --limit 5
gh run view <run-id> --log
```

Expected tags include the branch tag, `sha-*`, optional `latest`, and the
commit-traceable timestamp tag.

## 2. Registry Scan

```sh
flux get images all -n <image-resource-namespace>
kubectl get imagerepository <name> -n <namespace> -o yaml
```

The `ImageRepository` should be `Ready=True` and show the expected tags in
`status.lastScanResult.latestTags`.

## 3. Policy Selection

```sh
kubectl get imagepolicy <name> -n <namespace> -o yaml
```

Check:

- `status.conditions` is `Ready=True`
- `status.latestRef.tag` is the expected tag
- `status.latestRef.digest` is present when digest reflection is enabled

If the policy selected the wrong image, inspect `filterTags.pattern`,
`filterTags.extract`, and the policy ordering.

## 4. Git Update Commit

```sh
kubectl describe imageupdateautomation <name> -n <namespace>
kubectl get imageupdateautomation <name> -n <namespace> -o yaml
```

Healthy automation reports `Ready=True`, `lastPushCommit`, and `lastPushTime`.

Common failures:

- `referenced git repository does not exist`: add or correct
  `spec.sourceRef.namespace` or `spec.sourceRef.name`.
- repository is up to date but Deployment did not change: the setter comment
  may point to the wrong policy or the `update.path` may not include the file.
- authentication or push errors: check GitRepository credentials and controller
  logs.

## 5. Kustomization Apply

```sh
kubectl get gitrepository <name> -n <namespace> -o wide
kubectl get kustomization <name> -n <namespace> -o wide
```

The `GitRepository` and app `Kustomization` should advance to the automation
commit revision.

## 6. Deployment Rollout

```sh
kubectl get deploy <name> -n <namespace> -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}{.metadata.generation}{"\n"}{.status.observedGeneration}{"\n"}'
kubectl get pods -n <namespace> -l <selector> -o wide
```

The Deployment image should match the Flux-selected image. A pod should roll
when the pod template image string changes.

## Validation Commands

Before committing manifest changes:

```sh
kubectl kustomize <overlay>
kubectl apply --dry-run=server -k <overlay>
git diff --check
```

Server-side dry-run is valuable because it validates CRD fields against the
cluster API server.
