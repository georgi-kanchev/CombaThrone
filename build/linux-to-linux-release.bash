#!/usr/bin/env bash
set -euo pipefail

# Create the release directory if it doesn't exist
mkdir -p linux-release

echo "Building Linux Release..."
# Output the binary into the new folder
go build -ldflags="-s -w" -o linux-release/game ../

echo "Build complete: ./linux-release/game"
echo "Press Enter to exit..."
read
