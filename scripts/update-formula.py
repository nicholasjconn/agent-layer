#!/usr/bin/env python3
"""Update URL and SHA256 in a Homebrew formula file.

Usage: update-formula.py <formula-file> <new-url> <new-sha256>

Replaces exactly one `url "..."` line and exactly one `sha256 "..."` line.
Exits non-zero if replacements cannot be made.
"""
import re
import sys
from pathlib import Path

URL_PATTERN = re.compile(r'^(\s*url\s+").*("\s*)$', re.M)
SHA_PATTERN = re.compile(r'^(\s*sha256\s+").*("\s*)$', re.M)


def main() -> int:
    if len(sys.argv) != 4:
        usage = f"Usage: {sys.argv[0]} <formula-file> <new-url> <new-sha256>"
        print(usage, file=sys.stderr)
        return 1

    formula_path = Path(sys.argv[1])
    new_url = sys.argv[2]
    new_sha = sys.argv[3]

    if not formula_path.exists():
        print(f"Error: {formula_path} not found", file=sys.stderr)
        return 1

    text = formula_path.read_text()

    url_matches = URL_PATTERN.findall(text)
    if len(url_matches) != 1:
        print(f"Error: expected 1 url line, found {len(url_matches)}", file=sys.stderr)
        return 1

    sha_matches = SHA_PATTERN.findall(text)
    if len(sha_matches) != 1:
        msg = f"Error: expected 1 sha256 line, found {len(sha_matches)}"
        print(msg, file=sys.stderr)
        return 1

    text = URL_PATTERN.sub(r'\g<1>' + new_url + r'\g<2>', text)
    text = SHA_PATTERN.sub(r'\g<1>' + new_sha + r'\g<2>', text)
    formula_path.write_text(text)

    return 0


if __name__ == "__main__":
    sys.exit(main())
