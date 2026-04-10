/// <reference types="vitest/config" />

import { tanstackStart } from '@tanstack/react-start/plugin/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

export default defineConfig(() => ({
	plugins: process.env.VITEST ? [react()] : [tanstackStart(), react()],
	test: {
		environment: 'jsdom',
	},
}));
