#!/usr/bin/env bash
set -euo pipefail
python3 tools/generate_assets.py
python3 tools/audit.py
