import { Group, Title, Button, Image } from '@mantine/core';
import { Link, useLocation } from 'react-router-dom';
import { IconPlus, IconHorse } from '@tabler/icons-react';
import logo from '../../assets/he_logga_stor.png';

export function Navbar() {
  const location = useLocation();
  const isActive = (path: string) => location.pathname === path;

  return (
    <Group h="100%" px="md" justify="space-between">
      <Group>
        <Link to="/" style={{ textDecoration: 'none' }}>
          <Group>
            <Image src={logo} alt="Horse Tracker Logo" h={40} w="auto" />
            <Title order={2} c="blue.9">Horse Tracker</Title>
          </Group>
        </Link>
      </Group>

      <Group>
        <Button
          component={Link}
          to="/"
          variant={isActive('/') ? 'filled' : 'light'}
          color="blue"
          leftSection={<IconHorse size="1rem" />}
        >
          Horses
        </Button>

        <Button
          component={Link}
          to="/add-horse"
          variant={isActive('/add-horse') ? 'filled' : 'light'}
          color="green"
          leftSection={<IconPlus size="1rem" />}
        >
          Add Horse
        </Button>
      </Group>
    </Group>
  );
}
