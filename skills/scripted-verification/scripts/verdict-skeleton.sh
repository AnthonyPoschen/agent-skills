#!/usr/bin/env bash
# Verdict script skeleton. Replace measure() and the thresholds.
# Exit codes: 0 PASS, 1 FAIL, 2 NEEDS-HUMAN, 3 check error.
set -euo pipefail

out="${TMPDIR:-/tmp}/verify-$(date +%Y%m%d-%H%M%S)"
while [ $# -gt 0 ]; do
  case "$1" in
    --out) out="$2"; shift 2 ;;
    *) echo "usage: $0 [--out DIR]" >&2; exit 3 ;;
  esac
done
mkdir -p "$out"

verdict() { echo "$1 $2 evidence=$out"; exit "$3"; }
trap 'verdict ERROR "check failed to run" 3' ERR

measure() {
  # Run the real check. Keep raw output in $out; echo one number.
  ./run-check > "$out/raw.log" 2>&1
  awk '/score/ {print $2}' "$out/raw.log"
}

pass_max=0.02   # above the measured noise floor
fail_min=0.10   # clearly broken
score=$(measure)
echo "$score" > "$out/score.txt"

if awk "BEGIN{exit !($score <= $pass_max)}"; then
  verdict PASS "score=$score" 0
elif awk "BEGIN{exit !($score >= $fail_min)}"; then
  verdict FAIL "score=$score" 1
else
  verdict NEEDS-HUMAN "score=$score in borderline band" 2
fi
