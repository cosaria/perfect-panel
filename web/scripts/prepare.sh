#!/bin/bash

# Skip entire script in CI environment
if [ "$CI" = "true" ]; then
  echo "CI environment detected, skipping script execution."
  exit 0
fi

# Exit if not a Git repository
if [ ! -d ".git" ]; then
  echo "Not a Git repository. Exiting."
  exit 1
fi

# Function to set up a Husky hook
setup_husky_hook() {
  local hook_name=$1
  local hook_content=$2

  if [ ! -f ".husky/$hook_name" ]; then
    echo "Setting up $hook_name hook..."
    echo "$hook_content" > ".husky/$hook_name" || echo "Failed to set up $hook_name hook. Skipping."
    chmod +x ".husky/$hook_name" || echo "Failed to make $hook_name hook executable. Skipping."
  else
    echo "$hook_name hook is already set up."
  fi
}

# Ensure Husky is installed and initialized
if [ ! -d ".husky" ]; then
  echo "Setting up Husky..."
  npx husky || echo "Failed to set up Husky. Skipping."
else
  echo "Husky is already set up."
fi

# Set up pre-commit hook
setup_husky_hook "pre-commit" "#!/bin/sh
. \"\$(dirname \"\$0\")/_/husky.sh\"

npx --no-install lint-staged"

# Set up commit-msg hook
setup_husky_hook "commit-msg" "#!/bin/sh
. \"\$(dirname \"\$0\")/_/husky.sh\"

npx --no-install commitlint --edit \"\$1\""

# Function to globally install an npm package if not installed
install_global_package() {
  local package_name=$1

  if ! npm list -g --depth=0 "$package_name" > /dev/null 2>&1; then
    echo "Installing $package_name globally..."
    npm install -g "$package_name" || echo "Failed to install $package_name globally. Skipping."
  else
    echo "$package_name is already installed."
  fi
}

# Check and install required global npm packages
install_global_package "@lobehub/i18n-cli"
install_global_package "@lobehub/commit-cli"

# Run lobe-commit interactively
echo "Running lobe-commit -i..."
lobe-commit -i || echo "lobe-commit failed. Skipping."
