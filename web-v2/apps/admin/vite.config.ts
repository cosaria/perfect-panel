import { tanstackStart } from '@tanstack/react-start/plugin/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: process.env.VITEST ? [react()] : [tanstackStart(), react()],
	test: {
		environment: 'jsdom',
	},
} as any);
