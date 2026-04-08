#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

violations=()

# 显式边界清单：这些路径一旦重新出现在仓库根目录，就说明前端 workspace 的根结构回流了。
for path in apps packages scripts bun.lock turbo.json biome.json tsconfig.json; do
  if [ -e "$ROOT_DIR/$path" ]; then
    violations+=("$path")
  fi
done

if [ "${#violations[@]}" -gt 0 ]; then
  echo "错误：仓库根目录检测到不应回流的前端 workspace 关键根路径。"
  for path in "${violations[@]}"; do
    echo "违规路径：$ROOT_DIR/$path"
  done
  exit 1
fi

# 不检查 package.json、.gitignore 之类根目录允许存在的文件：
# 这些文件本来就是仓库根的常驻约定，是否存在不能说明 workspace 根结构回流。
echo "通过：仓库根目录未发现前端 workspace 关键根路径回流。"
