#!/usr/bin/env bash
# Shell tests for the lint-rules scripts.
set -euo pipefail
cd "$(dirname "$0")"
./baseline-only-shrinks.test.sh
./max-file-size.test.sh
