##!/bin/bash

# Update dependencies in web workspace root
echo "Updating dependencies in web workspace root..."

bun update --latest

# Update dependencies in packages and apps directories
for dir in ./packages/* ./apps/*; do
  if [ -d "$dir" ]; then
    echo "Updating dependencies in $dir..."
    (cd "$dir" && bun update --latest)
  fi
done
