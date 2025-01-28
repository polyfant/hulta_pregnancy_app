
npm installnpm installimport { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'node:path';

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
				target: 'http://localhost:8080',
				changeOrigin: true,
				secure: false,
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
