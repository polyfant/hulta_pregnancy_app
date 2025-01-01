import { useAuth0 } from '@auth0/auth0-react';
import { AppShell, Button, Center, LoadingOverlay, Text } from '@mantine/core';
import AppRoutes from '../routes';
import Navbar from './Navbar/Navbar';

export default function AppContent() {
	const { isLoading, isAuthenticated, error, loginWithRedirect } = useAuth0();

	console.log('Auth State:', { isLoading, isAuthenticated, error }); // Debug log

	if (error) {
		return (
			<Center h='100vh'>
				<div>
					<Text c='red' mb='md'>
						Error: {error.message}
					</Text>
					<Button onClick={() => loginWithRedirect()}>
						Retry Login
					</Button>
				</div>
			</Center>
		);
	}

	if (isLoading) {
		return (
			<Center h='100vh'>
				<LoadingOverlay visible={true} />
			</Center>
		);
	}

	if (!isAuthenticated) {
		return (
			<Center h='100vh'>
				<Button onClick={() => loginWithRedirect()}>Log In</Button>
			</Center>
		);
	}

	return (
		<AppShell header={{ height: 60 }} padding='md'>
			<AppShell.Header>
				<Navbar />
			</AppShell.Header>
			<AppShell.Main>
				<AppRoutes />
			</AppShell.Main>
		</AppShell>
	);
}
