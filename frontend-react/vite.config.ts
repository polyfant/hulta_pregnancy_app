import react from '@vitejs/plugin-react';
import path from 'node:path';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [react()],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './src'),
		},
	},
	server: {
		port: 3000,
		proxy: {
			'/api': {
				target: 'https://api.hulta-foaltracker.app',
				changeOrigin: true,
				secure: true,
				rewrite: (path) => path.replace(/^\/api/, '/api/v1'),
			},
		},
		hmr: true,
		watch: {
			usePolling: false,
		},
	},
	build: {
		minify: false, // for development
		sourcemap: true,
	},
	optimizeDeps: {
		include: [
			'react',
			'react-dom',
			'@mantine/core',
			'@tanstack/react-query',
		],
	},
});
