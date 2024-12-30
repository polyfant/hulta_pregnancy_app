import { useParams, Link, useNavigate } from 'react-router-dom';
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
import { FamilyTree } from './FamilyTree';

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

export function HorseDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: horse, isLoading, error } = useQuery<Horse>({
    queryKey: ['horse', id],
    queryFn: async () => {
      const response = await fetch(`/api/horses/${id}`);
      if (!response.ok) throw new Error('Failed to fetch horse details');
      return response.json();
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {
      const response = await fetch(`/api/horses/${id}`, {
        method: 'DELETE'
      });
      if (!response.ok) throw new Error('Failed to delete horse');
    },
    onSuccess: () => {
      notifications.show({
        title: 'Success',
        message: 'Horse deleted successfully',
        color: 'green'
      });
      navigate('/');
    },
    onError: (error) => {
      notifications.show({
        title: 'Error',
        message: error.message,
        color: 'red'
      });
    }
  });

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this horse? This action cannot be undone.')) {
      deleteMutation.mutate();
    }
  };

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
          <Group>
            <IconHorse size={30} />
            <Title order={2}>{horse.name}</Title>
          </Group>
          <Group>
            <Button
              component={Link}
              to={`/horses/${id}/edit`}
              variant="light"
              leftSection={<IconEdit size="1rem" />}
            >
              Edit
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
              size="lg"
              aria-label="Delete horse"
              onClick={handleDelete}
              loading={deleteMutation.isPending}
            >
              <IconTrash size="1rem" />
            </ActionIcon>
          </Group>
        </Group>

        <Grid>
          <Grid.Col span={4}>
            <Image
              src={horse.imageUrl}
              alt={horse.name}
              radius="md"
              fallbackSrc="/placeholder-horse.jpg"
            />
          </Grid.Col>

          <Grid.Col span={8}>
            <Stack gap="md">
              <Group>
                <Badge
                  color={horse.gender === 'STALLION' || horse.gender === 'GELDING' ? 'blue' : 'pink'}
                  variant="light"
                  size="lg"
                  leftSection={horse.gender === 'STALLION' || horse.gender === 'GELDING' ? <IconMars size="0.8rem" /> : <IconVenus size="0.8rem" />}
                >
                  {horse.gender === 'STALLION' ? 'Stallion' : horse.gender === 'GELDING' ? 'Gelding' : 'Mare'}
                </Badge>
                {horse.conceptionDate && (
                  <Badge color="grape" variant="light" size="lg">
                    Pregnant
                  </Badge>
                )}
              </Group>

              <Grid>
                <Grid.Col span={6}>
                  <Text size="sm" c="dimmed">Breed</Text>
                  <Text>{horse.breed}</Text>
                </Grid.Col>
                <Grid.Col span={6}>
                  <Text size="sm" c="dimmed">Date of Birth</Text>
                  <Text>{new Date(horse.dateOfBirth).toLocaleDateString()}</Text>
                </Grid.Col>
                {horse.weight && (
                  <Grid.Col span={6}>
                    <Text size="sm" c="dimmed">Weight</Text>
                    <Text>{horse.weight} kg</Text>
                  </Grid.Col>
                )}
                {horse.age && (
                  <Grid.Col span={6}>
                    <Text size="sm" c="dimmed">Age</Text>
                    <Text>{horse.age}</Text>
                  </Grid.Col>
                )}
                {horse.motherId && (
                  <Grid.Col span={6}>
                    <Text size="sm" c="dimmed">Mother</Text>
                    <Text>{horse.externalMother}</Text>
                  </Grid.Col>
                )}
                {horse.fatherId && (
                  <Grid.Col span={6}>
                    <Text size="sm" c="dimmed">Father</Text>
                    <Text>{horse.externalFather}</Text>
                  </Grid.Col>
                )}
              </Grid>
            </Stack>
          </Grid.Col>
        </Grid>
      </Card>

      <Tabs defaultValue="health">
        <Tabs.List>
          <Tabs.Tab
            value="health"
            leftSection={<IconStethoscope size="1rem" />}
          >
            Health Records
          </Tabs.Tab>
          <Tabs.Tab
            value="vitals"
            leftSection={<IconHeart size="1rem" />}
          >
            Vital Signs
          </Tabs.Tab>
          <Tabs.Tab
            value="family"
            leftSection={<IconTree size="1rem" />}
          >
            Family Tree
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="health" pt="xl">
          <Text>Health records coming soon...</Text>
        </Tabs.Panel>

        <Tabs.Panel value="vitals" pt="xl">
          <Text>Vital signs tracking coming soon...</Text>
        </Tabs.Panel>

        <Tabs.Panel value="family" pt="xl">
          <Stack gap="md">
            <Title order={3}>Family Tree</Title>
            <FamilyTree horse={horse} />
          </Stack>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
