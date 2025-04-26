#!/usr/bin/env bash
set -euo pipefail

CLI=../cli/driftdune   # adjust path if needed

echo "==> Writing sample embeddings"
printf '[1,0,0]\n' > baseline.json
printf '[0,1,0]\n' > current.json

echo "==> Running detect pipeline"
$CLI detect -b baseline.json -c current.json

echo "==> Verifying report.json exists"
if [ ! -f report.json ]; then
  echo "❌ report.json not found!" >&2
  exit 1
fi

drift=$(jq .drift_score report.json)
echo "==> Got drift_score = $drift"

if [[ "$drift" != "1" && "$drift" != "1.0" ]]; then
  echo "❌ Expected drift_score 1.0, got $drift" >&2
  exit 1
fi

echo "✅ Integration test passed"

