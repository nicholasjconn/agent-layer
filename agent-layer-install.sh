#!/usr/bin/env bash
set -euo pipefail

VERSION="latest"

usage() {
  echo "Usage: $0 [--version <tag>]" >&2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      if [[ $# -lt 2 ]]; then
        usage
        exit 1
      fi
      VERSION="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required to install Agent Layer." >&2
  exit 1
fi

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)
    OS="darwin"
    ;;
  Linux)
    OS="linux"
    ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
 esac

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
 esac

ASSET="al-${OS}-${ARCH}"
BASE_URL="https://github.com/nicholasjconn/agent-layer/releases"
if [[ "$VERSION" == "latest" ]]; then
  URL="${BASE_URL}/latest/download/${ASSET}"
  SUMS_URL="${BASE_URL}/latest/download/SHA256SUMS"
else
  URL="${BASE_URL}/download/${VERSION}/${ASSET}"
  SUMS_URL="${BASE_URL}/download/${VERSION}/SHA256SUMS"
fi

curl -fsSL "$URL" -o ./al
chmod +x ./al

verify_checksum_shasum() {
  local sums_file
  sums_file="$(mktemp)"
  if ! curl -fsSL "$SUMS_URL" -o "$sums_file"; then
    echo "Warning: unable to download SHA256SUMS; skipping checksum verification." >&2
    rm -f "$sums_file"
    return 0
  fi
  local expected
  expected="$(grep " ${ASSET}$" "$sums_file" | awk '{print $1}')"
  rm -f "$sums_file"
  if [[ -z "$expected" ]]; then
    echo "Warning: checksum for ${ASSET} not found; skipping verification." >&2
    return 0
  fi
  echo "${expected}  ./al" | shasum -a 256 -c -
}

verify_checksum_sha256sum() {
  local sums_file
  sums_file="$(mktemp)"
  if ! curl -fsSL "$SUMS_URL" -o "$sums_file"; then
    echo "Warning: unable to download SHA256SUMS; skipping checksum verification." >&2
    rm -f "$sums_file"
    return 0
  fi
  local expected
  expected="$(grep " ${ASSET}$" "$sums_file" | awk '{print $1}')"
  rm -f "$sums_file"
  if [[ -z "$expected" ]]; then
    echo "Warning: checksum for ${ASSET} not found; skipping verification." >&2
    return 0
  fi
  echo "${expected}  ./al" | sha256sum -c -
}

if command -v shasum >/dev/null 2>&1; then
  verify_checksum_shasum
elif command -v sha256sum >/dev/null 2>&1; then
  verify_checksum_sha256sum
else
  echo "Warning: shasum/sha256sum not available; skipping checksum verification." >&2
fi

./al install
