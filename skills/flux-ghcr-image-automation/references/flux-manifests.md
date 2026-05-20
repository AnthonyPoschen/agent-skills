# Flux Manifests

## ImageRepository

Point `ImageRepository` at the GHCR image repository, not a tag:

```yaml
apiVersion: image.toolkit.fluxcd.io/v1
kind: ImageRepository
metadata:
  name: my-app
spec:
  image: ghcr.io/owner/my-app
  interval: 1m
```

## ImagePolicy For Branch-SHA-Timestamp Tags

Use a named capture group for the timestamp and select the largest timestamp:

```yaml
apiVersion: image.toolkit.fluxcd.io/v1
kind: ImagePolicy
metadata:
  name: my-app
spec:
  imageRepositoryRef:
    name: my-app
  filterTags:
    pattern: "^master-[a-fA-F0-9]+-(?P<ts>[1-9][0-9]*)$"
    extract: "$ts"
  policy:
    numerical:
      order: asc
  digestReflectionPolicy: IfNotPresent
```

Use `digestReflectionPolicy: IfNotPresent` for immutable timestamped tags. Use
`Always` only when intentionally tracking a mutable tag such as `latest`.

## Deployment Setter Comment

Flux only updates image fields that carry a setter comment:

```yaml
containers:
  - name: app
    image: ghcr.io/owner/my-app:master-5838cfa-1779289200 # {"$imagepolicy": "app-namespace:my-app"}
```

The value inside `$imagepolicy` must be:

```text
<ImagePolicy namespace>:<ImagePolicy name>
```

This must match the namespace after Flux applies the resource. A parent Flux
`Kustomization` with `targetNamespace` can move namespaced resources into the
target namespace.

## ImageUpdateAutomation

Configure `ImageUpdateAutomation` to update the Git path that contains the
setter comment:

```yaml
apiVersion: image.toolkit.fluxcd.io/v1
kind: ImageUpdateAutomation
metadata:
  name: my-app
spec:
  interval: 1m
  sourceRef:
    kind: GitRepository
    name: my-app
    namespace: flux-system
  git:
    checkout:
      ref:
        branch: master
    commit:
      author:
        name: fluxcdbot
        email: fluxcdbot@users.noreply.github.com
      messageTemplate: "chore(images): update my-app image [skip ci]"
    push:
      branch: master
  update:
    path: ./kustomization/base
    strategy: Setters
```

Set `sourceRef.namespace` when the `GitRepository` lives outside the
`ImageUpdateAutomation` namespace.

## Namespace Checklist

When automation does not push commits:

- Find where `ImageRepository`, `ImagePolicy`, and `ImageUpdateAutomation`
  actually exist.
- Find where the referenced `GitRepository` exists.
- Compare `ImageUpdateAutomation.spec.sourceRef.name` and
  `spec.sourceRef.namespace` with the live `GitRepository`.
- Compare the image setter comment with the live `ImagePolicy` namespace/name.
- Check service account permissions if source refs cross namespaces.
