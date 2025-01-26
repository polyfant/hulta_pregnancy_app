import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Auth0Provider } from '@auth0/auth0-react';
import AppContent from './components/AppContent';

const queryClient = new QueryClient();

export default function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<MantineProvider
				theme={{
					colorScheme: 'dark',
					primaryColor: 'blue',
					colors: {
						dark: [
							'#C1C2C5',
							'#A6A7AB',
							'#909296',
							'#5c5f66',
							'#373A40',
							'#2C2E33',
							'#25262b',
							'#1A1B1E',
							'#141517',
							'#101113',
						],
					},
				}}
			>
				<Notifications />
				<AppContent />
			</MantineProvider>
		</QueryClientProvider>
	);
}
