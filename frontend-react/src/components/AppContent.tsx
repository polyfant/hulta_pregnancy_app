import { useAuth0 } from '@auth0/auth0-react';
import { AppShell, Button, Center, Group, LoadingOverlay, Text, Title } from '@mantine/core';
import { Suspense } from 'react';
import { Link } from 'react-router-dom';
import AppRoutes from '../routes';

export default function AppContent() {
	const { isLoading, isAuthenticated, error, loginWithRedirect, logout, user } = useAuth0();

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
		<AppShell
			header={{ height: 60 }}
			padding='md'
			styles={(theme) => ({
				main: {
					backgroundColor: theme.colors.dark[8],
				},
				header: {
					backgroundColor: theme.colors.dark[7],
				}
			})}
		>
			<AppShell.Header>
				<Group h="100%" px="md" justify="space-between">
					<Title 
						order={2} 
						c="white" 
						component={Link} 
						to="/"
						style={{ 
							textDecoration: 'none',
							cursor: 'pointer'
						}}
					>
						Horse Tracker
					</Title>
					<Group>
						<Text c="white">Welcome, {user?.name}</Text>
						<Button 
							variant="filled"
							color="brown"
							onClick={() => logout()}
						>
							Logout
						</Button>
					</Group>
				</Group>
			</AppShell.Header>

			<AppShell.Main>
				<Suspense fallback={<LoadingOverlay visible />}>
					<AppRoutes />
				</Suspense>
			</AppShell.Main>
		</AppShell>
	);
}
