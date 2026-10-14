# Lint Rules Catalog

Hits are non-vendored, non-test matches when the rule was added, summed over
the repos it was measured on. "Sampled" means hits were read by hand.

| Rule | Why | Bad example | Added | Origin | Hits |
| --- | --- | --- | --- | --- | --- |
| `no-bare-todo` (Go semgrep) | Untracked TODO/FIXME/HACK rot. A linked issue gives an owner and exit. | `// TODO: fix later` | 2026-10-14 | coding skill comment rules | 0 in 2 Go repos |
| `nolintlint` require-explanation (Go) | A silent suppression hides why the linter was wrong. | `//nolint:errcheck` | 2026-10-14 | coding skill comment rules | not measured |
| `no-error-string-match` (Go semgrep) | Error text changes and wrapping breaks the match; behaviour silently flips. | `strings.Contains(err.Error(), "not found")` | 2026-10-14 | web backend branching on provider error messages | 11 in 1 of 2 Go repos; sampled 11/11 real |
| `errorf-wrap-with-w` (Go semgrep) | `%v`/`%s` with no `%w` cuts the error chain, so `errors.Is/As` fail upstream. `"%w: %v"` with a sentinel stays allowed. | `fmt.Errorf("rollback failed: %s", err)` | 2026-10-14 | same message copy-pasted across a family of importers | 12 in 1 of 2; sampled 12/12 real |
| `no-default-http-client` (Go semgrep) | No timeout: a server that accepts and stalls hangs the caller forever. | `http.DefaultClient.Do(req)` | 2026-10-14 | review comment on an unbounded response-header wait; fallback clients in API integrations | 5 in 2 of 2; sampled 5/5 real |
| `no-ignored-unmarshal-error` (Go semgrep) | A dropped decode error leaves a zero-value struct that looks valid. | `_ = json.Unmarshal(b, &v)` | 2026-10-14 | DB JSON columns and API payloads decoded blind | 5 in 1 of 2; sampled 5/5 real |
| `bufio-scanner-set-buffer` (Go semgrep) | Scanner dies on a line over 64 KiB and the rest of the stream is lost. | `s := bufio.NewScanner(r)` with no `Buffer` | 2026-10-14 | review comment: log streaming stopped at one oversized line | 2 in 2 of 2 |
| `no-sleep-poll-loop` (Go semgrep) | Sleep-in-loop ignores cancellation and has no bound or backoff. | `for !ready() { time.Sleep(50 * time.Millisecond) }` | 2026-10-14 | process supervisors polling for a child to start | 3 in 1 of 2; sampled 3/3 real |
| `no-local-replace` (go.mod semgrep) | A local-path replace builds only on one machine and breaks CI. | `replace example.com/lib => ../lib` | 2026-10-14 | corrected repeatedly in agent sessions (about 20 sessions) | 0 committed now |
| errorlint, bodyclose, noctx, contextcheck (golangci) | Stock checks for wrapping, leaked bodies and dropped contexts; back the semgrep rules above. | `resp, _ := http.Get(u)` | 2026-10-14 | general Go practice | not measured (golangci not run) |
| `unwrap_used` / `expect_used` (Rust) | Panics in non-test code crash the process. Return or handle the error. | `let v = map.get(k).unwrap();` | 2026-10-14 | general Rust practice | ~180 unwrap, ~300 expect (all targets) |
| `panic` / `todo` / `unimplemented` (Rust) | Same as above for explicit panics in library code. | `panic!("not a mesh")` | 2026-10-14 | renderer library panicking on bad input | 2 in library code (17 incl. tests/build.rs, allowed) |
| `print_stdout` / `print_stderr` / `dbg_macro` (Rust) | Libraries report through return values or a logger; stray prints corrupt CLI/protocol output. | `println!("frame {i}")` in a lib crate | 2026-10-14 | debug prints left in library crates | 18 in library code (114 incl. bins/tests) |
| `cast_possible_truncation` / `_wrap` / `sign_loss` (Rust, warn) | `as` silently truncates or wraps. Use `From`/`TryFrom`. | `len as u32` | 2026-10-14 | review comments on unchecked integer overflow | ~675 `as` casts in library code; gate new code only |
| `disallowed-methods std::env::var` (Rust clippy.toml) | Library code reading env hides inputs from callers and tests; parse once in `main` and pass config. | `std::env::var("TIER_MS")` in a lib | 2026-10-14 | tuning knobs read from env deep inside a renderer | 17 in library code |
| `disallowed-methods std::thread::sleep` (Rust clippy.toml) | Fixed sleeps waste time when fast and flake when slow; wait on an event or bounded deadline. | `thread::sleep(Duration::from_millis(4000))` | 2026-10-14 | window/audio libraries waiting with fixed sleeps | 10 in library code |
| `no-latest-image-tag` (k8s semgrep) | `:latest` or no tag is not reproducible and defeats rollbacks. Handles `registry:5000/app`. | `image: registry:5000/app` | 2026-10-14 | general k8s practice; registry-port gap fixed | 3 (all kustomize placeholders: lint rendered output) |
| `container-resource-limits` (k8s semgrep) | No memory limit lets one leaking pod OOM the node. | container with no `resources.limits.memory` | 2026-10-14 | general k8s practice | 0 in 2 repos' manifests |
| kubeconform strict (k8s) | Typos in manifests apply silently or fail at deploy time. | `replicas: "3"` | 2026-10-14 | general k8s practice | not measured |
