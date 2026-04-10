import react from '@vitejs/plugin-react'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'

export default {
	plugins: process.env.VITEST ? [react()] : [tanstackStart(), react()],
	test: {
		environment: 'jsdom',
	},
} as any
