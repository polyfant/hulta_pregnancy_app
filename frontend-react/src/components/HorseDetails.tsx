import { useParams, Link, useNavigate } from 'react-router-dom';
import { useCallback } from 'react';
import {
  Card,
  Image,
  Text,
  Group,
  Badge,
  Button,
  ActionIcon,
  Stack,
  Title,
  Grid,
  LoadingOverlay,
  Tabs,
  Paper,
  Box
} from '@mantine/core';
import {
  IconEdit,
  IconTrash,
  IconHorse,
  IconCalendar,
  IconRuler,
  IconWeight,
  IconVaccine,
  IconNotes,
  IconBabyCarriage,
  IconMars,
  IconVenus,
  IconStethoscope,
  IconHeart,
  IconTree
} from '@tabler/icons-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { lazy, Suspense } from 'react';

// Lazy load FamilyTree component
const FamilyTree = lazy(() => import('./FamilyTree/FamilyTree'));
const PregnancyStatus = lazy(() => import('./PregnancyTracking/PregnancyStatus'));

interface Horse {
  id: number;
  name: string;
  breed: string;
  gender: 'MARE' | 'STALLION' | 'GELDING';
  dateOfBirth: string;
  weight?: number;
  age?: string;
  conceptionDate?: string;
  motherId?: number;
  fatherId?: number;
  externalMother?: string;
  externalFather?: string;
  created_at?: string;
  updated_at?: string;
  imageUrl?: string;
}

const HorseDetails = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: horse, isLoading, error } = useQuery<Horse>({
    queryKey: ['horse', id],
    queryFn: async () => {
      const response = await fetch(`/api/horses/${id}`);
      if (!response.ok) throw new Error('Failed to fetch horse details');
      return response.json();
    },
    staleTime: 30000,
    gcTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {
      const response = await fetch(`/api/horses/${id}`, {
        method: 'DELETE'
      });
      if (!response.ok) throw new Error('Failed to delete horse');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horses'] });
      notifications.show({
        title: 'Success',
        message: 'Horse deleted successfully',
        color: 'green'
      });
      navigate('/');
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Error',
        message: error.message,
        color: 'red'
      });
    }
  });

  const handleDelete = useCallback(() => {
    if (window.confirm('Are you sure you want to delete this horse? This action cannot be undone.')) {
      deleteMutation.mutate();
    }
  }, [deleteMutation]);

  if (isLoading) {
    return (
      <Paper p="xl" pos="relative">
        <LoadingOverlay visible />
      </Paper>
    );
  }

  if (error || !horse) {
    return (
      <Paper p="xl">
        <Text c="red">Error loading horse details</Text>
      </Paper>
    );
  }

  return (
    <Stack gap="lg">
      <Card withBorder>
        <Group justify="space-between" mb="md">
          <Group gap="sm">
            <IconHorse size={30} />
            <Title order={2}>{horse.name}</Title>
          </Group>
          <Group gap="sm">
            <Button
              component={Link}
              to={`/horses/${id}/edit`}
              variant="filled"
              leftSection={<IconEdit size="1rem" />}
              styles={(theme) => ({
                root: {
                  color: theme.colors.green[4],
                  backgroundColor: theme.colors.dark[7],
                  '&:hover': {
                    backgroundColor: theme.colors.dark[8],
                  }
                }
              })}
            >
              Edit Horse
            </Button>
            {horse.gender === 'MARE' && (
              <Button
                component={Link}
                to={`/horses/${id}/pregnancy`}
                variant="light"
                color="blue"
                leftSection={<IconBabyCarriage size="1rem" />}
              >
                Pregnancy Tracking
              </Button>
            )}
            <ActionIcon
              color="red"
              variant="light"
              onClick={handleDelete}
              loading={deleteMutation.isPending}
            >
              <IconTrash size="1rem" />
            </ActionIcon>
          </Group>
        </Group>

        <Tabs defaultValue="details">
          <Tabs.List>
            <Tabs.Tab value="details" leftSection={<IconNotes size="1rem" />}>
              Details
            </Tabs.Tab>
            <Tabs.Tab value="health" leftSection={<IconStethoscope size="1rem" />}>
              Health
            </Tabs.Tab>
            {horse.gender === 'MARE' && horse.conceptionDate && (
              <Tabs.Tab value="pregnancy" leftSection={<IconBabyCarriage size="1rem" />}>
                Pregnancy
              </Tabs.Tab>
            )}
            <Tabs.Tab value="family" leftSection={<IconTree size="1rem" />}>
              Family Tree
            </Tabs.Tab>
          </Tabs.List>

          <Box mt="md">
            <Tabs.Panel value="details">
              <Paper p="md" withBorder>
                <Grid>
                  <Grid.Col span={{ base: 12, md: 6 }}>
                    <Stack>
                      <Group>
                        <IconCalendar size="1rem" />
                        <Text>Born: {new Date(horse.dateOfBirth).toLocaleDateString()}</Text>
                      </Group>
                      <Group>
                        {horse.gender === 'STALLION' ? (
                          <IconMars size="1rem" color="blue" />
                        ) : (
                          <IconVenus size="1rem" color="pink" />
                        )}
                        <Text>Gender: {horse.gender}</Text>
                      </Group>
                      {horse.breed && (
                        <Group>
                          <IconHorse size="1rem" />
                          <Text>Breed: {horse.breed}</Text>
                        </Group>
                      )}
                      {horse.weight && (
                        <Group>
                          <IconWeight size="1rem" />
                          <Text>Weight: {horse.weight} kg</Text>
                        </Group>
                      )}
                    </Stack>
                  </Grid.Col>
                  <Grid.Col span={{ base: 12, md: 6 }}>
                    {horse.imageUrl && (
                      <Image
                        src={horse.imageUrl}
                        alt={horse.name}
                        radius="md"
                        fit="cover"
                      />
                    )}
                  </Grid.Col>
                </Grid>
              </Paper>
            </Tabs.Panel>

            <Tabs.Panel value="health">
              <Paper p="md" withBorder>
                <Stack>
                  <Group>
                    <IconHeart size="1rem" />
                    <Text>Health Status: Healthy</Text>
                  </Group>
                  <Group>
                    <IconVaccine size="1rem" />
                    <Text>Last Vaccination: Up to date</Text>
                  </Group>
                </Stack>
              </Paper>
            </Tabs.Panel>

            {horse.gender === 'MARE' && horse.conceptionDate && (
              <Tabs.Panel value="pregnancy">
                <Paper p="md" withBorder>
                  <Suspense fallback={<LoadingOverlay visible />}>
                    <PregnancyStatus horseId={horse.id} />
                  </Suspense>
                </Paper>
              </Tabs.Panel>
            )}

            <Tabs.Panel value="family">
              <Paper p="md" withBorder>
                <Suspense fallback={<LoadingOverlay visible />}>
                  <FamilyTree horseId={horse.id} />
                </Suspense>
              </Paper>
            </Tabs.Panel>
          </Box>
        </Tabs>
      </Card>
    </Stack>
  );
};

export default HorseDetails;
