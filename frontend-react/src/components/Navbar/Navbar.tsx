import { useAuth0 } from '@auth0/auth0-react';
import { Avatar, Button, Group, Text, Tooltip } from '@mantine/core';
import { Baby } from '@phosphor-icons/react';
import { Link, useLocation } from 'react-router-dom';

export default function Navbar() {
	const location = useLocation();
	const { user, logout } = useAuth0();

	return (
		<Group h='100%' px='md' justify='space-between'>
			<Link to='/' style={{ textDecoration: 'none' }}>
				<Text size='xl' fw={700} c='blue.9'>
					Horse Tracker
				</Text>
			</Link>

			<Group>
				{location.pathname !== '/' && (
					<Button component={Link} to='/' variant='subtle'>
						Horses
					</Button>
				)}

				<Tooltip label='Pregnancy Overview'>
					<Button
						component={Link}
						to='/pregnancies'
						variant='subtle'
						leftSection={<Baby size='1.2rem' />}
					>
						Pregnancies
					</Button>
				</Tooltip>

				<Group>
					{user?.picture && <Avatar src={user.picture} radius='xl' />}
					<Text>{user?.name}</Text>
					<Button
						onClick={() =>
							logout({
								logoutParams: {
									returnTo: window.location.origin,
								},
							})
						}
					>
						Logout
					</Button>
				</Group>
			</Group>
		</Group>
	);
}
