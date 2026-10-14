# Lint Rules Catalog

| Rule | Why | Bad example | Added | Origin |
| --- | --- | --- | --- | --- |
| `no-bare-todo` (Go semgrep) | Untracked TODO/FIXME/HACK rot. A linked issue gives an owner and exit. | `// TODO: fix later` | 2026-10-14 | coding skill comment rules |
| `nolintlint` require-explanation (Go) | A silent suppression hides why the linter was wrong. | `//nolint:errcheck` | 2026-10-14 | coding skill comment rules |
| `unwrap_used` / `expect_used` (Rust) | Panics in non-test code crash the process. Return or handle the error. | `let v = map.get(k).unwrap();` | 2026-10-14 | general Rust practice |
| `no-latest-image-tag` (k8s semgrep) | `:latest` is not reproducible and defeats rollbacks. | `image: nginx:latest` | 2026-10-14 | general k8s practice |
| kubeconform strict (k8s) | Typos in manifests apply silently or fail at deploy time. | `replicas: "3"` | 2026-10-14 | general k8s practice |
