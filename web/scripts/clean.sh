#!/bin/bash

echo "Cleaning all node_modules, bun.lock, and bun.lockb files..."

# Remove all node_modules directories
find . -name "node_modules" -type d -prune -exec rm -rf '{}' +

# Remove all bun.lock files
find . -name "bun.lock" -type f -delete

# Remove all bun.lockb files
find . -name "bun.lockb" -type f -delete