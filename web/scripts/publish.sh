#!/bin/bash

# Set up directories
OUT_DIR="out"

# Clean up any existing build artifacts
rm -rf $OUT_DIR
mkdir -p $OUT_DIR

# Declare an array of projects to build
PROJECTS=(
  "ppanel-admin-web:apps/admin:3001"
  "ppanel-user-web:apps/user:3002"
)

# Step 1: Install dependencies
bun install || {
  echo "Dependency installation failed"
  exit 1
}

# Step 2: Build each project using Turbo
for ITEM in "${PROJECTS[@]}"; do
  IFS=":" read -r PROJECT PROJECT_PATH DEFAULT_PORT <<< "$ITEM"
  echo "Building project: $PROJECT (Path: $PROJECT_PATH)"
  bun run build --filter=$PROJECT || {
    echo "Build failed for $PROJECT"
    exit 1
  }
  # Copy build output and static resources to the build directory
  PROJECT_BUILD_DIR=$OUT_DIR/$PROJECT
  cp -r $PROJECT_PATH/.next/standalone/. $PROJECT_BUILD_DIR/
  cp -r $PROJECT_PATH/.next/static $PROJECT_BUILD_DIR/$PROJECT_PATH/.next/
  cp -r $PROJECT_PATH/public $PROJECT_BUILD_DIR/$PROJECT_PATH/
  cp -r $PROJECT_PATH/.env.template $PROJECT_BUILD_DIR/$PROJECT_PATH/.env

  # Generate ecosystem.config.js for the project
  ECOSYSTEM_CONFIG="$PROJECT_BUILD_DIR/ecosystem.config.js"
  cat > $ECOSYSTEM_CONFIG << EOL
module.exports = {
  apps: [
    {
      name: "$PROJECT",
      script: "$PROJECT_PATH/server.js",
      interpreter: "bun",
      watch: true,
      instances: "max",
      exec_mode: "cluster",
      env: {
        PORT: $DEFAULT_PORT
      }
    }
  ]
};
EOL
  echo "PM2 configuration created: $ECOSYSTEM_CONFIG"

  # Create a tar.gz archive for each project
  ARCHIVE_NAME="$OUT_DIR/$PROJECT.tar.gz"
  tar -czvf $ARCHIVE_NAME -C $OUT_DIR $PROJECT || {
    echo "Archiving failed for $PROJECT"
    exit 1
  }
  echo "Archive created: $ARCHIVE_NAME"
done

# Final output
echo "All projects have been built and archived successfully."
