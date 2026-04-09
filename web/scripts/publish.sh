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
  # Copy the built static SPA assets into the archive payload
  PROJECT_BUILD_DIR=$OUT_DIR/$PROJECT
  mkdir -p $PROJECT_BUILD_DIR/$PROJECT_PATH
  cp -r $PROJECT_PATH/dist $PROJECT_BUILD_DIR/$PROJECT_PATH/
  if [ -d "$PROJECT_PATH/public" ]; then
    cp -r $PROJECT_PATH/public $PROJECT_BUILD_DIR/$PROJECT_PATH/
  fi
  if [ -f "$PROJECT_PATH/.env.template" ]; then
    cp $PROJECT_PATH/.env.template $PROJECT_BUILD_DIR/$PROJECT_PATH/.env
  fi

  # Generate ecosystem.config.js for the project
  ECOSYSTEM_CONFIG="$PROJECT_BUILD_DIR/ecosystem.config.js"
  cat > $ECOSYSTEM_CONFIG << EOL
module.exports = {
  apps: [
    {
      name: "$PROJECT",
      script: "bun",
      args: "run --cwd $PROJECT_PATH preview -- --host 0.0.0.0 --port $DEFAULT_PORT",
      watch: false,
      instances: 1,
      exec_mode: "fork",
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
