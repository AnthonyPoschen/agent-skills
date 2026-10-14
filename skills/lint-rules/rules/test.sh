#!/usr/bin/env bash
# Runs every rule's failing example. Needs semgrep (pip) and cargo clippy.
set -euo pipefail
cd "$(dirname "$0")"
semgrep --test --metrics=off .
out=$(cd rust/fixture && cargo clippy --quiet --all-targets --target-dir "${TMPDIR:-/tmp}/lint-rules-fixture" 2>&1 || true)
want=("#panic" "#print_stdout" "#unwrap_used" "#cast_possible_truncation" "#disallowed_methods" "disallowed method \`std::env::var\`" "disallowed method \`std::thread::sleep\`")
for w in "${want[@]}"; do
  grep -qF -- "$w" <<<"$out" || { echo "rust fixture: missing $w"; exit 1; }
done
echo "rust fixture: all clippy rules fire"
