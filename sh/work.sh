#!/usr/bin/env bash
# find-configmaps-multicluster.sh  v0.2

set -euo pipefail
IFS=$'\n\t'

# 依赖检查
command -v kubectl >/dev/null || { echo >&2 "kubectl not found"; exit 1; }
command -v jq       >/dev/null || { echo >&2 "jq not found"; exit 1; }

DOMAIN="${1:-marvel-interface.pix.prod}"

for CTX in $(kubectl config get-contexts -o name); do
  echo "===== Context: $CTX ====="

  MATCHES=$(kubectl --context="$CTX" get configmap -A -o json |
    jq -r --arg dom "$DOMAIN" '
      .items[]
      | select(
          ((.data // {})      | tostring | test($dom; "i")) or
          ((.binaryData // {})| tostring | test($dom; "i"))
        )
      | "\(.metadata.namespace)/\(.metadata.name)"
    ' | sort -u)

  if [[ -n "$MATCHES" ]]; then
    echo "$MATCHES"
  else
    echo "(no match)"
  fi

  echo
done