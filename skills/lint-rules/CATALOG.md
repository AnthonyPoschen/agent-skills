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
| `color-no-hex` / `color-named` / `function-disallowed-list` (CSS Stylelint) | Hard-coded colours scatter the palette so a recolour or theme switch misses values. Keep raw colours in the theme file; everywhere else use `var(--token)`. | `color: #f4f6f8` | 2026-10-17 | agents missing shared colour tokens during theme recolours | not measured |
| `lint/no-hex-color` (CSS Stylelint) | Same ban as `color-no-hex`, with a configurable theme-file exemption, so a raw hex cannot hide in a shorthand the built-in misses. | `outline: 1px solid #abc` | 2026-10-20 | palette enforcement on stylesheets | not measured |
| `scale-unlimited/declaration-strict-value` (CSS Stylelint) | Colour properties must be a token, `currentColor`, or `transparent`, so a named colour or function cannot sit on `color` while hex is banned elsewhere. | `color: red` | 2026-10-20 | palette enforcement on stylesheets | not measured |
| `lint/no-old-scheme` (CSS Stylelint) | Deprecated token names keep a retired palette alive after a recolour. The deprecated list and prefix are configurable. | `color: var(--palette-accent)` | 2026-10-20 | retired palette tokens still referenced after a recolour | not measured |
| `lint/colors-in-theme` (CSS Stylelint) | A colour custom property defined outside the theme file is a second palette. Colour prefixes and the theme path are configurable. | `--palette-local: #abcdef` | 2026-10-20 | colour tokens declared beside components | not measured |
| `lint/no-placeholder-color` (CSS Stylelint) | Per-field placeholder colours drift. One `::placeholder` rule in the theme file owns that token. | `input::placeholder { color: var(--palette-ink) }` | 2026-10-20 | placeholder colour set on individual fields | not measured |
| `lint/no-clipped-heading` (CSS Stylelint) | A fixed height or width plus `overflow: hidden`, or `line-height` under 1.15, clips heading glyphs. | `h1 { height: 24px; overflow: hidden }` | 2026-10-20 | display headings cut off by a fixed box | not measured |
| `lint/no-layout-animation` (CSS Stylelint) | Animating width, height, insets, margin, or padding shifts layout (CLS). Animate transform, opacity, filter, or colour. | `transition: height 200ms ease` | 2026-10-20 | layout properties in transitions | not measured |
| `lint/unknown-token` (CSS Stylelint) | `var(--prefix-*)` missing from the theme file is a silent fallback, including when `var()` has a fallback. Prefix and theme path are configurable. | `var(--palette-missing)` | 2026-10-20 | typos in token names | not measured |
| `lint/no-page-component-css` (CSS Stylelint) | Component classes and raw colours in page CSS drift from the shared stylesheet and the theme. Class list and exempt files are configurable. | `.card { color: var(--ink) }` in a page | 2026-10-20 | component CSS copied into pages | not measured |
| `declaration-no-important` (CSS Stylelint) | `!important` breaks the cascade so later rules cannot override it. | `color: var(--ink) !important` | 2026-10-20 | CSS hygiene | not measured |
| `declaration-block-no-duplicate-properties` (CSS Stylelint) | A repeated property leaves the earlier declaration dead, and a later edit can miss it. | `color` set twice in one block | 2026-10-20 | CSS hygiene | not measured |
| `no-duplicate-selectors` (CSS Stylelint) | One selector split across blocks drifts when only half is updated. | `.example` twice | 2026-10-20 | CSS hygiene | not measured |
| `declaration-block-no-shorthand-property-overrides` (CSS Stylelint) | A longhand before its shorthand never applies. | `margin-top` then `margin` | 2026-10-20 | CSS hygiene | not measured |
| `property-no-unknown` (CSS Stylelint) | An unknown property is dropped by the browser, so the declaration does nothing. | `colour: var(--ink)` | 2026-10-20 | CSS hygiene | not measured |
| `declaration-property-value-allowed-list` font-size / z-index (CSS Stylelint) | Raw font-size and z-index scatter type and stacking. Use a token (`0`, `-1`, and `auto` stay allowed for z-index). | `font-size: 14px` | 2026-10-20 | type and stacking tokens | not measured |
| inline-color (HTML/JS scan) | A colour in `style=`, `element.style`, or canvas `fillStyle` / `strokeStyle` / `shadowColor` is missed by a stylesheet recolour. | `style="color: #abc"` | 2026-10-20 | colours outside stylesheets | not measured |
| svg-palette (SVG scan) | An SVG fill or stroke whose hex is not in the theme file is off-palette. | `fill="#00ff00"` | 2026-10-20 | SVG assets with their own colours | not measured |
| `sqlclosecheck` (Go golangci) | An unclosed `database/sql` rows or statement leaks a connection. | `rows, err := db.Query(...)` with no `Close` | 2026-10-20 | Go SQL call sites | not measured |
| `rowserrcheck` (Go golangci) | Ignoring `rows.Err` hides a truncated result that looks complete. | loop `rows.Next` and return nil | 2026-10-20 | Go SQL call sites | not measured |
| `forbidigo` (Go golangci) | `os.Getenv` outside `config/`, `fmt.Print` / `log.Fatal` outside `main.go` and `cmd/`, and the default HTTP client have no single source, logger, or timeout. `_test.go` may print. | `os.Getenv("APP_ENV")` in a library | 2026-10-20 | env, logging, and HTTP call sites | not measured |

## Verdicts, 2026-10-14

Measured on shallow clones of genos, genos-cli and genos-controller (Go) and
genos-engine (Rust). Go counts exclude vendor; "non-test" excludes
`_test.go`. Rust counts are unique sites, split by target. Kubernetes counts
are on `kustomize build` output of all six overlays.

| Rule | Hits | Verdict |
| --- | --- | --- |
| `no-error-string-match` | 11 non-test (genos) | Keep. All sampled hits branch on error text. |
| `errorf-wrap-with-w` | 12 non-test (controller) | Keep. Cuts the chain; the `%w: %v` sentinel idiom stays allowed. |
| `no-default-http-client` | 5 non-test (genos 3, controller 2) | Keep. Every hit is an untimed outbound call. |
| `no-ignored-unmarshal-error` | 5 non-test (genos) | Keep. |
| `bufio-scanner-set-buffer` | 2 (genos dotenv, controller log scan) | Keep. Log scan is real; dotenv is low risk but cheap to fix. |
| `no-sleep-poll-loop` | 3 → 0 | Tightened. All 3 were 3 s deadline-bounded waits for a socket; loops bounded by `time.Now().Before`/`time.Until` are now exempt. |
| `no-local-replace` | 0 | Keep. No noise; guards a CI-breaking mistake. |
| `no-bare-todo` | 0 | Keep. No noise. |
| golangci `errorlint` | 60 Errorf + 4 others → 10 | Tightened. `errorf: false`: 46 hits were the allowed `%w: %v` idiom, the rest duplicated `errorf-wrap-with-w`. Comparison and type-assertion checks stay (10 hits, 5 in genos-cli tests). |
| golangci `bodyclose` | 0 | Keep. No noise. |
| golangci `noctx` | 908 in tests, 7 non-test | Tightened. Excluded in `_test.go`. Non-test hits are `exec.Command` and `Client.Get` without a context, which is real. |
| golangci `contextcheck` | 5 non-test | Keep. Low volume; suppress a deliberate detached context with an explained `//nolint`. |
| Rust `unwrap_used`/`expect_used`/`panic`/`todo`/`unimplemented` | lib 37; tests/ 139, build.rs 15, bin+example 2 | Scoped to library crates; build.rs and `tests/` allow at file level. Lib hits are the real target. |
| Rust `print_stdout`/`print_stderr`/`dbg_macro` | lib 6; bin 28, example 10, tests 5 | Scoped to library crates. Bins print by design. |
| Rust `disallowed_methods` `std::env::var(_os)` | lib 17 (13 tuning knobs, 4 platform discovery); bin 7, build.rs 3, tests 1 | Scoped to library crates. Platform discovery (HOME, XDG_*, WAYLAND_DISPLAY) allowed per item with a reason. |
| Rust `disallowed_methods` `std::thread::sleep` | lib 7 (5 fixed 0.3–5 s waits in window, 2 short polls in mcp); tests 1 | Keep, library crates only. Fixed waits are the target. |
| Rust `cast_possible_truncation`/`_wrap`/`sign_loss` | lib 445; tests 90, bin 22, example 9 | Keep as warn, library crates only, gate new code first. |
| `no-latest-image-tag` | 3 raw → 0 rendered | Keep. Raw hits were kustomize placeholders (`app`, `controller`) that render to pinned GHCR tags; the registry-port fix is covered by its test. |
| `container-resource-limits` | 0 rendered | Keep. All containers already set memory limits. |
