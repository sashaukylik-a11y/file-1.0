#!/usr/bin/env bash
set -euo pipefail
python3 tools/generate_assets.py
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags='-H=windowsgui -s -w' -o SNAPWAVE_HORROR_V8.exe .
echo 'Built SNAPWAVE_HORROR_V8.exe'
