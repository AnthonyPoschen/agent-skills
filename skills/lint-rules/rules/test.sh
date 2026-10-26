#!/usr/bin/env bash
# Runs every rule's failing example. Needs semgrep (pip), cargo clippy, and
# node/npm for the CSS Stylelint fixtures.
set -euo pipefail
cd "$(dirname "$0")"
semgrep --test --metrics=off .
out=$(cd rust/fixture && cargo clippy --quiet --all-targets --target-dir "${TMPDIR:-/tmp}/lint-rules-fixture" 2>&1 || true)
want=("#panic" "#print_stdout" "#unwrap_used" "#cast_possible_truncation" "#disallowed_methods" "disallowed method \`std::env::var\`" "disallowed method \`std::thread::sleep\`")
for w in "${want[@]}"; do
  grep -qF -- "$w" <<<"$out" || { echo "rust fixture: missing $w"; exit 1; }
done
echo "rust fixture: all clippy rules fire"
(
  cd css
  if [[ ! -d node_modules ]]; then
    npm ci --ignore-scripts
  fi
)
css_out=$(cd css && node --disable-warning=stylelint:002 --test test.js 2>&1 || true)
css_rules=(
  color-no-hex
  color-named
  function-disallowed-list
  scale-unlimited/declaration-strict-value
  lint/no-old-scheme
  lint/no-hex-color
  lint/colors-in-theme
  lint/no-placeholder-color
  lint/no-clipped-heading
  lint/no-layout-animation
  lint/unknown-token
  lint/no-page-component-css
  declaration-no-important
  declaration-block-no-duplicate-properties
  no-duplicate-selectors
  declaration-block-no-shorthand-property-overrides
  property-no-unknown
  declaration-property-value-allowed-list
  inline-color
  svg-palette
)
for rule in "${css_rules[@]}"; do
  grep -qF -- "css-rule fail ${rule}" <<<"$css_out" || { echo "css: failing fixture did not fail for $rule"; exit 1; }
done
for pass in good.css old-scheme-good.css hex-good.css builtin-color-good.css colors-in-theme-good.css placeholder-good.css clipped-heading-good.css layout-animation-good.css unknown-token-good.css page-component-good.css hygiene-good.css inline-color-good.html svg-palette-good.svg; do
  grep -qF -- "css-rule pass ${pass}" <<<"$css_out" || { echo "css: passing fixture did not pass: $pass"; exit 1; }
done
if grep -q "not ok" <<<"$css_out"; then
  echo "$css_out"
  echo "css: node test reported a failure"
  exit 1
fi
echo "css fixtures: each rule fails its bad fixture and passes its good fixture"
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "go fixtures: golangci-lint is not installed" >&2
  exit 1
fi
go_bad=$(cd go/fixture && golangci-lint run ./bad/... 2>&1 || true)
for w in sqlclosecheck rowserrcheck forbidigo; do
  grep -qF -- "$w" <<<"$go_bad" || { echo "go fixture: missing $w"; echo "$go_bad"; exit 1; }
done
echo "go fixture: sqlclosecheck, rowserrcheck, and forbidigo fire"
go_good=$(cd go/fixture && golangci-lint run ./good/... 2>&1 || true)
if grep -qE "sqlclosecheck|rowserrcheck|forbidigo" <<<"$go_good"; then
  echo "go fixture: passing files were flagged"
  echo "$go_good"
  exit 1
fi
echo "go fixture: config/, cmd/, and _test.go exclusions pass"
