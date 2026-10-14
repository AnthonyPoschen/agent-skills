# Kubernetes Rules

Validate rendered manifests with kubeconform in CI:

```sh
kustomize build overlays/prod | kubeconform -strict -summary \
  -schema-location default \
  -schema-location 'https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json'
```

Run `semgrep --config semgrep.yml` for the custom rules.
