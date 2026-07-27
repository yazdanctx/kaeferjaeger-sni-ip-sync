#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <domain>" >&2
  echo "Example: $0 dell.com" >&2
  exit 1
fi

DOMAIN="$1"
INPUT="final.txt"
OUTPUT_DIR="output/${DOMAIN//\./_}"
TIMESTAMP=$(date +%s%3N)

if [ ! -f "$INPUT" ]; then
  echo "Error: $INPUT not found. Run the sync tool first." >&2
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

OUTPUT_FILE="$OUTPUT_DIR/results_${TIMESTAMP}.txt"

grep -F ".$DOMAIN" "$INPUT" \
  | awk -F'-- ' '{print $2}' \
  | tr ' [' '\n\n' \
  | sed 's/\]//g' \
  | grep -F ".$DOMAIN" \
  | sort -u > "$OUTPUT_FILE"

echo "Extracted to $OUTPUT_FILE"
