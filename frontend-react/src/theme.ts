import { createTheme, MantineTheme, rem } from '@mantine/core';

// Theme inspired by Hulta Equestrian
export const theme = createTheme({
  primaryColor: 'brand',
  colors: {
    // Main brand colors
    brand: [
      '#FFE4CC', // 0: Lightest
      '#FFD4B3', // 1
      '#FFC499', // 2
      '#FFB380', // 3
      '#FFA266', // 4
      '#FF914D', // 5: Primary brand color
      '#FF8033', // 6
      '#FF6F1A', // 7
      '#FF5E00', // 8
      '#E65400', // 9: Darkest
    ],
    // Dark theme colors
    dark: [
      '#F8F9FA', // 0: Lightest - text
      '#E9ECEF', // 1
      '#DEE2E6', // 2
      '#CED4DA', // 3
      '#212529', // 4
      '#1A1D20', // 5
      '#141619', // 6: Main background
      '#101214', // 7: Card background
      '#0A0B0C', // 8: Deeper elements
      '#050506', // 9: Darkest - modals
    ],
    // Accent colors for various states
    accent: [
      '#E9FFF2', // 0: Success lightest
      '#CCFFE0', // 1
      '#B3FFD1', // 2
      '#99FFC2', // 3
      '#80FFB3', // 4
      '#66FFA3', // 5: Success primary
      '#4DFF94', // 6
      '#33FF85', // 7
      '#1AFF76', // 8
      '#00FF66', // 9: Success darkest
    ],
  },
  components: {
    AppShell: {
      styles: (theme: MantineTheme) => ({
        main: {
          background: theme.colors.dark[6],
          color: theme.colors.dark[0],
        },
      }),
    },
    Card: {
      styles: (theme: MantineTheme) => ({
        root: {
          backgroundColor: theme.colors.dark[7],
          color: theme.colors.dark[0],
          border: `1px solid ${theme.colors.dark[8]}`,
          transition: 'transform 0.2s ease, box-shadow 0.2s ease',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: `0 4px 8px ${theme.colors.dark[9]}`,
          },
        },
      }),
    },
    Button: {
      styles: (theme: MantineTheme) => ({
        root: {
          transition: 'transform 0.2s ease',
          '&:active': {
            transform: 'translateY(1px)',
          },
        },
      }),
    },
    Modal: {
      styles: (theme: MantineTheme) => ({
        content: {
          backgroundColor: theme.colors.dark[7],
        },
        header: {
          backgroundColor: theme.colors.dark[8],
          color: theme.colors.dark[0],
        },
      }),
    },
    Tabs: {
      styles: (theme: MantineTheme) => ({
        tab: {
          '&[data-active]': {
            color: theme.colors.brand[5],
            borderColor: theme.colors.brand[5],
          },
        },
      }),
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
