module.exports = {
	root: true,
	env: {
		browser: true,
		es2021: true,
	},
	extends: [
		'eslint:recommended',
		'plugin:@typescript-eslint/recommended',
		'plugin:react/recommended',
		'plugin:react-hooks/recommended',
	],
	parser: '@typescript-eslint/parser',
	parserOptions: {
		ecmaFeatures: {
			jsx: true,
		},
		ecmaVersion: 'latest',
		sourceType: 'module',
	},
	plugins: ['@typescript-eslint', 'react', 'react-hooks'],
	settings: {
		react: {
			version: 'detect',
		},
	},
	rules: {},
	overrides: [
		{
			files: ['src/**/*.{ts,tsx}'],
			parserOptions: {
				project: './tsconfig.json',
			},
		},
		{
			files: ['vite.config.ts'],
			parserOptions: {
				project: './tsconfig.node.json',
			},
		},
	],
};
