import react from '@vitejs/plugin-react';
import path from 'node:path';
import { defineConfig, PluginOption } from 'vite';

// Explicitly type the plugins array
const plugins: PluginOption[] = [react()];

export default defineConfig({
	plugins: plugins,
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
