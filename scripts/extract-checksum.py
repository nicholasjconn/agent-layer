#!/usr/bin/env python3
"""Extract SHA256 checksum for a specific file from checksums.txt.

Usage: extract-checksum.py <checksums-file> <target-filename>

Supports checksum formats:
  <hash>  <filename>
  sha256:<hash>  <filename>

Prints the hash to stdout and exits 0 on success, exits 1 if not found.
"""
import re
import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) != 3:
        usage = f"Usage: {sys.argv[0]} <checksums-file> <target-filename>"
        print(usage, file=sys.stderr)
        return 1

    checksums_path = Path(sys.argv[1])
    target = sys.argv[2]

    if not checksums_path.exists():
        print(f"Error: {checksums_path} not found", file=sys.stderr)
        return 1

    for line in checksums_path.read_text().splitlines():
        line = line.strip()
        if not line:
            continue
        parts = line.split()
        if len(parts) < 2:
            continue
        h = parts[0]
        f = parts[-1].lstrip("./")
        if f == target or f == target.lstrip("./"):
            h = re.sub(r"^sha256:", "", h)
            print(h)
            return 0

    print(f"Error: {target} not found in {checksums_path}", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
