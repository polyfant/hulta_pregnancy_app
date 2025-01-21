import { Grid, Card, Text, Group, RingProgress, Stack, Badge } from '@mantine/core';
import { Horse, FirstAid, Calendar, Clock } from '@phosphor-icons/react';
import { FC } from 'react';

interface DashboardStats {
  pregnancyProgress: number;
  daysRemaining: number;
  stage: string;
  nextCheckup: string;
}

export const HorseStatusDashboard: FC = () => {
  // This would typically come from your API/state management
  const stats: DashboardStats = {
    pregnancyProgress: 65,
    daysRemaining: 119,
    stage: 'Mid Stage',
    nextCheckup: '3 days',
  };

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
                <Text fw={700} size="xl">Today</Text>
              </div>
              <Clock size={32} weight="fill" />
            </Group>
            <Text size="sm" c="dimmed">2 hours ago</Text>
          </Stack>
        </Card>
      </Grid.Col>
    </Grid>
  );
};