import { Grid, Card, Text, Group, RingProgress, Stack, Badge, Skeleton } from '@mantine/core';
import { Horse, FirstAid, Calendar, Clock } from '@phosphor-icons/react';
import { FC } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useHorsesApi } from '../../api/horses';

interface DashboardStats {
  pregnancyProgress: number;
  daysRemaining: number;
  stage: string;
  nextCheckup: string;
  lastUpdated: string;
}

interface HorseStatusDashboardProps {
  horseId: number;
}

export const HorseStatusDashboard: FC<HorseStatusDashboardProps> = ({ horseId }) => {
  const api = useHorsesApi();
  const { data: stats, isLoading, error } = useQuery<DashboardStats>({
    queryKey: ['horseStatusDashboard', horseId],
    queryFn: () => api.getDashboardStats(horseId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });

  if (isLoading) {
    return (
      <Grid gutter="md">
        {[1, 2, 3, 4].map((col) => (
          <Grid.Col key={col} span={{ base: 12, md: 6, lg: 3 }}>
            <Card shadow="sm" padding="lg" radius="md">
              <Skeleton height={150} />
            </Card>
          </Grid.Col>
        ))}
      </Grid>
    );
  }

  if (error || !stats) {
    return (
      <Grid gutter="md">
        <Grid.Col span={12}>
          <Card shadow="sm" padding="lg" radius="md">
            <Text c="red">Failed to load horse status</Text>
          </Card>
        </Grid.Col>
      </Grid>
    );
  }

  return (
    <Grid gutter="md">
      <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
        <Card shadow="sm" padding="lg" radius="md">
          <Stack gap="md">
            <Group justify="space-between">
              <div>
                <Text fw={500} size="sm" c="dimmed">Pregnancy Progress</Text>
                <Text fw={700} size="xl">{stats.pregnancyProgress}%</Text>
              </div>
              <RingProgress
                size={80}
                thickness={8}
                sections={[{ value: stats.pregnancyProgress, color: 'brand.5' }]}
                label={
                  <Horse size={20} weight="fill" />
                }
              />
            </Group>
            <Badge color="brand" variant="light">
              {stats.stage}
            </Badge>
          </Stack>
        </Card>
      </Grid.Col>

      <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
        <Card shadow="sm" padding="lg" radius="md">
          <Stack gap="md">
            <Group justify="space-between">
              <div>
                <Text fw={500} size="sm" c="dimmed">Days Remaining</Text>
                <Text fw={700} size="xl">{stats.daysRemaining}</Text>
              </div>
              <Calendar size={32} weight="fill" />
            </Group>
            <Text size="sm" c="dimmed">Until Expected Due Date</Text>
          </Stack>
        </Card>
      </Grid.Col>

      <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
        <Card shadow="sm" padding="lg" radius="md">
          <Stack gap="md">
            <Group justify="space-between">
              <div>
                <Text fw={500} size="sm" c="dimmed">Next Checkup</Text>
                <Text fw={700} size="xl">{stats.nextCheckup}</Text>
              </div>
              <FirstAid size={32} weight="fill" />
            </Group>
            <Text size="sm" c="dimmed">Veterinary Visit</Text>
          </Stack>
        </Card>
      </Grid.Col>

      <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
        <Card shadow="sm" padding="lg" radius="md">
          <Stack gap="md">
            <Group justify="space-between">
              <div>
                <Text fw={500} size="sm" c="dimmed">Last Updated</Text>
                <Text fw={700} size="xl">{stats.lastUpdated}</Text>
              </div>
              <Clock size={32} weight="fill" />
            </Group>
            <Text size="sm" c="dimmed">Horse Status</Text>
          </Stack>
        </Card>
      </Grid.Col>
    </Grid>
  );
};