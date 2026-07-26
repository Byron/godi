#!/bin/sh
set -eu

old=${1:?usage: compare-performance.sh OLD_GODI NEW_GODI}
new=${2:?usage: compare-performance.sh OLD_GODI NEW_GODI}
work=${TMPDIR:-/tmp}/godi-performance-$$
trap 'rm -rf "$work"' EXIT INT TERM
mkdir -p "$work/source"
dd if=/dev/zero of="$work/source/large.bin" bs=1m count=1024 2>/dev/null

for binary in "$old" "$new"; do
  rm -f "$work/source"/godi_*.gobz
  echo "== $binary =="
  /usr/bin/time "$binary" --verbosity=off seal "$work/source"
done
