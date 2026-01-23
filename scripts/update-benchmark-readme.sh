#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BENCH_DIR="$ROOT_DIR/benchmarks"
README="$ROOT_DIR/README.md"

START_MARKER="<!-- BENCHMARK_START -->"
END_MARKER="<!-- BENCHMARK_END -->"

cd "$BENCH_DIR"

echo "Running benchmarks..."
go test -bench=. -benchmem -count=1 -benchtime=1s > benchmark_results.txt

echo "Formatting results..."
awk '
BEGIN {
  print "| Benchmark | ns/op | B/op | allocs/op |"
  print "|-----------|------:|-----:|----------:|"
}
/^Benchmark/ {
  name = $1
  gsub(/^Benchmark/, "", name)
  gsub(/-[0-9]+$/, "", name)
  ns_op = $3
  b_op = $5
  allocs = $7
  printf "| %s | %s | %s | %s |\n", name, ns_op, b_op, allocs
}
' benchmark_results.txt > benchmark_table.md

echo "Updating README..."
awk -v start="$START_MARKER" -v end="$END_MARKER" -v tablefile="$BENCH_DIR/benchmark_table.md" '
BEGIN { skip = 0 }
index($0, start) {
  print
  print ""
  while ((getline line < tablefile) > 0) print line
  print ""
  skip = 1
  next
}
index($0, end) {
  skip = 0
}
!skip { print }
' "$README" > "$README.tmp" && mv "$README.tmp" "$README"

rm -f benchmark_results.txt benchmark_table.md

echo "Done!"
