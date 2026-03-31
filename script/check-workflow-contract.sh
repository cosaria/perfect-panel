#!/bin/sh

set -eu

fail() {
  echo "ERROR: $1" >&2
  exit 1
}

for path in web/.github/workflows server/.github/workflows; do
  if [ -d "$path" ] && find "$path" -type f | grep -q .; then
    fail "Subdirectory workflows are forbidden: $path"
  fi
done

[ -f .github/workflows/monorepo-check.yml ] || fail "Missing .github/workflows/monorepo-check.yml"
[ -f .github/workflows/web-release.yml ] || fail "Missing .github/workflows/web-release.yml"
[ -f .github/workflows/web-publish-release-assets.yml ] || fail "Missing .github/workflows/web-publish-release-assets.yml"
[ -f .github/workflows/server-release.yml ] || fail "Missing .github/workflows/server-release.yml"
[ -f .github/workflows/server-develop.yml ] || fail "Missing .github/workflows/server-develop.yml"
[ -f .github/workflows/server-swagger.yml ] || fail "Missing .github/workflows/server-swagger.yml"

grep -Fq 'web-v\${version}' web/.releaserc.js || fail "web semantic-release tag must be namespaced"
grep -Fq "startsWith(github.event.release.tag_name, 'web-v')" .github/workflows/web-publish-release-assets.yml || fail "web publish workflow must gate on web-v tags"
grep -Fq -- "- 'v*'" .github/workflows/server-release.yml || fail "server release workflow must listen to v* tags"
grep -Fq 'name: Monorepo Summary' .github/workflows/monorepo-check.yml || fail "monorepo summary check must exist"

echo "Workflow contract checks passed."
