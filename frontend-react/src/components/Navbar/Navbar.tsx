import { Group, Button, Text, Avatar } from '@mantine/core';
import { Link, useLocation } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';

export default function Navbar() {
  const location = useLocation();
  const { isAuthenticated, user, logout } = useAuth0();

  return (
    <Group h="100%" px="md" justify="space-between">
      <Link to="/" style={{ textDecoration: 'none' }}>
        <Text size="xl" fw={700} c="blue.9">Horse Tracker</Text>
      </Link>

      <Group>
        {location.pathname !== '/' && (
          <Button
            component={Link}
            to="/"
            variant="subtle"
          >
            Horses
          </Button>
        )}

        {isAuthenticated && (
          <Group>
            {user?.picture && (
              <Avatar src={user.picture} radius="xl" />
            )}
            <Text>{user?.name}</Text>
            <Button 
              variant="subtle" 
              onClick={() => logout({ 
                logoutParams: { 
                  returnTo: window.location.origin 
                } 
              })}
            >
              Logout
            </Button>
          </Group>
        )}
      </Group>
    </Group>
  );
}
