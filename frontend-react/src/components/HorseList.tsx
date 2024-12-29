import { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Card,
  Text,
  Badge,
  Button,
  Group,
  Stack,
  TextInput,
  Grid,
  LoadingOverlay,
  ActionIcon,
  Title,
  SimpleGrid,
  Alert
} from '@mantine/core';
import {
  IconHorse,
  IconSearch,
  IconPlus,
  IconMars,
  IconVenus
} from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';

interface Horse {
  id: string;
  name: string;
  breed?: string;
  color?: string;
  gender: 'male' | 'female';
  birthDate?: string;
  isPregnant?: boolean;
}

export function HorseList() {
  const [searchQuery, setSearchQuery] = useState('');

  const { data, isLoading, error } = useQuery<Horse[]>({
    queryKey: ['horses'],
    queryFn: async () => {
      const response = await fetch('/api/horses');
      if (!response.ok) throw new Error('Failed to fetch horses');
      return response.json();
    },
    initialData: [] // Set initial data as empty array
  });

  // Ensure data is never null
  const horses = data || [];

  const filteredHorses = horses.filter(horse =>
    horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (horse.breed && horse.breed.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  if (isLoading) {
    return (
      <Card withBorder p="xl" pos="relative">
        <LoadingOverlay visible />
      </Card>
    );
  }

  if (error) {
    return (
      <Card withBorder p="xl">
        <Text c="red">Error loading horses</Text>
      </Card>
    );
  }

  return (
    <Stack gap="lg">
      <Card withBorder>
        <Group justify="space-between" mb="md">
          <Group>
            <IconHorse size={30} />
            <Title order={2}>Horses</Title>
          </Group>
          <Button
            component={Link}
            to="/add-horse"
            leftSection={<IconPlus size="1rem" />}
          >
            Add Horse
          </Button>
        </Group>

        <TextInput
          placeholder="Search horses..."
          value={searchQuery}
          onChange={(event) => setSearchQuery(event.currentTarget.value)}
          mb="md"
          leftSection={<IconSearch size="1rem" />}
        />

        {filteredHorses.length === 0 ? (
          <Text c="dimmed" ta="center" py="xl">
            No horses found
          </Text>
        ) : (
          <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
            {filteredHorses.map((horse) => (
              <Grid.Col key={horse.id} span={{ base: 12, sm: 6, md: 4 }}>
                <Card
                  withBorder
                  component={Link}
                  to={`/horses/${horse.id}`}
                  style={{ textDecoration: 'none', color: 'inherit' }}
                >
                  <Group justify="space-between">
                    <Text fw={500}>{horse.name}</Text>
                    <Group gap="xs">
                      <Badge
                        color={horse.gender === 'male' ? 'blue' : 'pink'}
                        variant="light"
                        leftSection={
                          horse.gender === 'male' 
                            ? <IconMars size="0.8rem" />
                            : <IconVenus size="0.8rem" />
                        }
                      >
                        {horse.gender === 'male' ? 'Stallion' : 'Mare'}
                      </Badge>
                      {horse.isPregnant && (
                        <Badge color="grape" variant="light">
                          Pregnant
                        </Badge>
                      )}
                    </Group>
                  </Group>

                  <Text size="sm" c="dimmed" mt="xs">
                    {horse.breed || 'Unknown breed'}
                    {horse.birthDate && ` • Born ${new Date(horse.birthDate).getFullYear()}`}
                  </Text>
                </Card>
              </Grid.Col>
            ))}
          </SimpleGrid>
        )}
      </Card>
    </Stack>
  );
}
