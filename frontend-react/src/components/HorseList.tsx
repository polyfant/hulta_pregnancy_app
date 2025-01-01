import React from 'react';
import { useState } from 'react';
import { 
  Text, 
  Title, 
  Stack, 
  Group, 
  Button, 
  TextInput, 
  Card, 
  Grid, 
  Badge, 
  LoadingOverlay, 
  ActionIcon, 
  SimpleGrid, 
  Alert
} from '@mantine/core';
import { 
  MdAlertCircle, 
  FiPlus, 
  FiSearch, 
  FiEdit, 
  FiTrash2, 
  MdMale, 
  MdFemale, 
  MdPets
} from '@/utils/icons';
import { Link } from 'react-router-dom';

interface Horse {
  id: string;
  name: string;
  breed?: string;
  color?: string;
  gender: 'male' | 'female';
  birthDate?: string;
  isPregnant?: boolean;
}

const HorseList = () => {
  const [searchQuery, setSearchQuery] = useState('');

  const { data, isLoading, error } = useQuery<Horse[]>({
    queryKey: ['horses'],
    queryFn: async () => {
      console.log('Fetching horses...');
      try {
        const response = await fetch('/api/horses');
        console.log('API Response:', response);
        
        if (!response.ok) {
          const errorData = await response.text();
          console.error('API Error:', errorData);
          throw new Error(`Failed to fetch horses: ${errorData}`);
        }
        
        const data = await response.json();
        console.log('Horses data:', data);
        return data;
      } catch (error) {
        console.error('Error fetching horses:', error);
        throw error;
      }
    },
    retry: 1,
    staleTime: 30000,
    refetchOnWindowFocus: false
  });

  // Ensure data is never null
  const horses = data || [];

  // Filter horses based on search query
  const filteredHorses = horses.filter(horse =>
    horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (horse.breed && horse.breed.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  if (error) {
    return (
      <Alert icon={<MdAlertCircle size="1rem" />} title="Error" color="red">
        Failed to load horses. Please try again later.
      </Alert>
    );
  }

  return (
    <Stack>
      <Group justify="space-between" align="center">
        <Title order={2}>Horses</Title>
        <Button
          component={Link}
          to="/add-horse"
          variant="filled"
          styles={(theme) => ({
            root: {
              color: theme.colors.green[4],
              backgroundColor: theme.colors.dark[7],
              '&:hover': {
                backgroundColor: theme.colors.dark[8],
              }
            }
          })}
          leftSection={<FiPlus size="1rem" />}
        >
          Add Horse
        </Button>
      </Group>

      <TextInput
        icon={<FiSearch size="1rem" />}
        placeholder="Search horses..."
        value={searchQuery}
        onChange={(event) => setSearchQuery(event.currentTarget.value)}
      />

      <div style={{ position: 'relative' }}>
        <LoadingOverlay visible={isLoading} />
        <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="md">
          {filteredHorses.map((horse) => (
            <Card key={horse.id} shadow="sm" padding="lg" radius="md" withBorder>
              <Group justify="space-between" mb="xs">
                <Text fw={500}>{horse.name}</Text>
                <ActionIcon
                  variant="light"
                  color={horse.gender === 'male' ? 'blue' : 'pink'}
                  title={horse.gender}
                >
                  {horse.gender === 'male' ? <MdMale size="1.2rem" /> : <MdFemale size="1.2rem" />}
                </ActionIcon>
              </Group>

              <Group gap="xs">
                {horse.breed && (
                  <Badge color="blue" variant="light">
                    {horse.breed}
                  </Badge>
                )}
                {horse.isPregnant && (
                  <Badge color="grape" variant="light">
                    Pregnant
                  </Badge>
                )}
              </Group>

              <Button
                component={Link}
                to={`/horses/${horse.id}`}
                variant="light"
                color="blue"
                fullWidth
                mt="md"
                radius="md"
                leftSection={<MdPets size="1rem" />}
              >
                View Details
              </Button>
            </Card>
          ))}
        </SimpleGrid>
      </div>
    </Stack>
  );
};

export default HorseList;
