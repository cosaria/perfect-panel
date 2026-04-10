import { existsSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

const root = resolve(import.meta.dirname, '..', '..')

describe('openapi pipeline', () => {
	it('defines generated api-client package', () => {
		expect(existsSync(resolve(root, 'packages/api-client/package.json'))).toBe(true)
		expect(existsSync(resolve(root, 'packages/api-client/scripts/generate.mjs'))).toBe(true)
	})
})
