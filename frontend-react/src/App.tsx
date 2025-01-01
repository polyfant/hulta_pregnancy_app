import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import AppContent from './components/AppContent';
import { theme } from './theme';

const queryClient = new QueryClient();

export default function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<MantineProvider theme={theme}>
				<Notifications />
				<AppContent />
			</MantineProvider>
		</QueryClientProvider>
	);
}
