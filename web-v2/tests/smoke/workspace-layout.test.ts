import { existsSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..', '..')

describe('workspace layout', () => {
  it('has root workspace files', () => {
    expect(existsSync(resolve(root, 'package.json'))).toBe(true)
    expect(existsSync(resolve(root, 'pnpm-workspace.yaml'))).toBe(true)
    expect(existsSync(resolve(root, 'turbo.json'))).toBe(true)
    expect(existsSync(resolve(root, 'biome.json'))).toBe(true)
  })
})
