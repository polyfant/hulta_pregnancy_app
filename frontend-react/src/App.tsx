import React from 'react';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Auth0Provider } from '@auth0/auth0-react';
import AppContent from './components/AppContent';
import ErrorBoundary from './components/ErrorBoundary';
import { theme } from './theme';

const queryClient = new QueryClient();

export default function App() {
	return (
		<Auth0Provider
			domain={import.meta.env.VITE_AUTH0_DOMAIN}
			clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
			authorizationParams={{
				redirect_uri: window.location.origin,
				audience: import.meta.env.VITE_AUTH0_AUDIENCE,
			}}
		>
			<QueryClientProvider client={queryClient}>
				<MantineProvider theme={theme}>
					<Notifications />
					<ErrorBoundary>
						<AppContent />
					</ErrorBoundary>
				</MantineProvider>
			</QueryClientProvider>
		</Auth0Provider>
	);
}
