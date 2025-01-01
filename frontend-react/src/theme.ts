import { createTheme, rem } from '@mantine/core';

// Theme inspired by Hulta Equestrian
export const theme = createTheme({
	primaryColor: 'dark',
	colors: {
		// Custom color palette inspired by Hulta's dark theme
		dark: [
			'#C1C2C5', // 0: Lightest - text
			'#A6A7AB', // 1
			'#909296', // 2
			'#5C5F66', // 3
			'#373A40', // 4
			'#2C2E33', // 5
			'#25262B', // 6
			'#1A1B1E', // 7: Main background
			'#141517', // 8
			'#101113', // 9: Darkest - deeper backgrounds
		],
		accent: [
			'#F8E4D0', // 0
			'#E6C4A0', // 1
			'#D4A470', // 2
			'#C28440', // 3
			'#B06410', // 4: Primary accent
			'#955409', // 5
			'#7A4407', // 6
			'#5F3405', // 7
			'#442403', // 8
			'#291401', // 9
		],
		// ... rest of your color configurations
	},
	components: {
		AppShell: {
			styles: {
				main: {
					background: 'var(--mantine-color-dark-8)',
					color: 'white',
				},
			},
		},
		Card: {
			styles: {
				root: {
					backgroundColor: 'var(--mantine-color-dark-7)',
					color: 'white',
				},
			},
		},
		Text: {
			styles: {
				root: {
					color: 'white',
				},
			},
		},
		Title: {
			styles: {
				root: {
					color: 'white',
				},
			},
		},
		Button: {
			defaultProps: {
				variant: 'filled',
				color: 'accent',
			},
			styles: {
				root: {
					'&[dataFilled]': {
						backgroundColor: 'var(--mantine-color-dark-7)',
						color: 'var(--mantine-color-accent-4)',
					},
				},
			},
		},
		TextInput: {
			defaultProps: {
				variant: 'filled',
			},
		},
		InputWrapper: {
			styles: {
				required: {
					color: 'var(--mantine-color-red-6)',
					marginLeft: '0.2rem',
				},
			},
		},
	},
	fontSizes: {
		xs: rem(12),
		sm: rem(14),
		md: rem(16),
		lg: rem(18),
		xl: rem(20),
	},
	radius: {
		xs: rem(2),
		sm: rem(4),
		md: rem(8),
		lg: rem(16),
		xl: rem(32),
	},
	spacing: {
		xs: rem(8),
		sm: rem(12),
		md: rem(16),
		lg: rem(24),
		xl: rem(32),
	},
});
