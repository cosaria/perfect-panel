#!/bin/bash

cd apps/user

bunx --bun shadcn@canary add -a -o
bunx --bun shadcn@canary add https://ui.aceternity.com/registry/timeline.json -o
bunx --bun shadcn@canary add https://ui.aceternity.com/registry/text-generate-effect.json -o
bunx --bun shadcn@canary add https://ui.aceternity.com/registry/hover-border-gradient.json -o

bunx --bun shadcn@canary add https://magicui.design/r/styles/default/orbiting-circles.json -o
bunx --bun shadcn@canary add https://magicui.design/r/styles/default/hyper-text.json -o
