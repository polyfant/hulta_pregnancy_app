import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Card,
  Stack,
  Group,
  Title,
  Text,
  Button,
  Progress,
  Grid
} from '@mantine/core';
import {
  IconBabyCarriage,
  IconPlus,
  IconX,
  IconCalendarEvent,
  IconAlertTriangle
} from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { Horse } from '../../types/horse';
import { PregnancyProgress } from '../../types/pregnancy';
import { StartPregnancyDialog } from './StartPregnancyDialog';
import { EndPregnancyDialog } from './EndPregnancyDialog';
import { formatDate } from '../../utils/dateUtils';

interface PregnancyTrackingProps {
  horse: Horse;
}

interface PregnancyStatus {
  daysInPregnancy: number;
  expectedDueDate: Date;
  currentStage: string;
  stageDescription: string;
  progressPercentage: number;
  upcomingEvents: string[];
  warningSignsList: string[];
  preFoalingProgress: number;
  preFoalingSignsCount: number;
}

export function PregnancyTracking({ horse }: PregnancyTrackingProps) {
  const [startDialogOpened, setStartDialogOpen] = useState(false);
  const [endDialogOpened, setEndDialogOpen] = useState(false);
  const queryClient = useQueryClient();

  const { data: pregnancyStatus } = useQuery<PregnancyStatus>({
    queryKey: ['pregnancyStatus', horse.id],
    enabled: horse.isPregnant,
  });

  const startPregnancyMutation = useMutation({
    mutationFn: async (date: string) => {
      const response = await fetch(`/api/horses/${horse.id}/pregnancy/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ startDate: date }),
      });
      if (!response.ok) throw new Error('Failed to start pregnancy tracking');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horse', horse.id] });
      notifications.show({
        title: 'Success',
        message: 'Started pregnancy tracking',
        color: 'green',
      });
      setStartDialogOpen(false);
    },
  });

  const endPregnancyMutation = useMutation({
    mutationFn: async (data: { outcome: string; foalingDate: string }) => {
      const response = await fetch(`/api/horses/${horse.id}/pregnancy/end`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error('Failed to end pregnancy tracking');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horse', horse.id] });
      notifications.show({
        title: 'Success',
        message: 'Ended pregnancy tracking',
        color: 'green',
      });
      setEndDialogOpen(false);
    },
  });

  const handleStartPregnancy = (date: string) => {
    startPregnancyMutation.mutate(date);
  };

  const handleEndPregnancy = (data: { outcome: string; foalingDate: string }) => {
    endPregnancyMutation.mutate(data);
  };

  const openStartDialog = () => setStartDialogOpen(true);
  const closeStartDialog = () => setStartDialogOpen(false);
  const openEndDialog = () => setEndDialogOpen(true);
  const closeEndDialog = () => setEndDialogOpen(false);

  return (
    <Stack gap="lg">
      <Group justify="space-between" mb="md">
        <Group>
          <IconBabyCarriage size={30} />
          <Title order={2}>Pregnancy Tracking</Title>
        </Group>
        <Group>
          {!horse.isPregnant ? (
            <Button
              onClick={openStartDialog}
              leftSection={<IconPlus size="1rem" />}
              loading={startPregnancyMutation.isPending}
            >
              Start Tracking
            </Button>
          ) : (
            <Button
              onClick={openEndDialog}
              leftSection={<IconX size="1rem" />}
              loading={endPregnancyMutation.isPending}
              color="red"
            >
              End Tracking
            </Button>
          )}
        </Group>
      </Group>

      {horse.isPregnant && pregnancyStatus && (
        <>
          <Stack gap={8}>
            <Title order={3}>Pregnancy Progress</Title>
            <Text size="sm" c="dimmed">
              {pregnancyStatus.daysInPregnancy} days into pregnancy
            </Text>
            <Progress
              value={pregnancyStatus.progressPercentage}
              size="xl"
              aria-label="Pregnancy progress"
            />
            <Group justify="space-between" mt="xs">
              <Text size="sm">Conception Date</Text>
              <Text size="sm">{formatDate(pregnancyStatus.expectedDueDate)}</Text>
            </Group>
          </Stack>

          <Card withBorder>
            <Stack gap="md">
              <Title order={3}>Current Stage: {pregnancyStatus.currentStage}</Title>
              <Text>
                {pregnancyStatus.stageDescription}
              </Text>
            </Stack>
          </Card>

          <Grid>
            <Grid.Col span={6}>
              <Card withBorder h="100%">
                <Stack gap="md">
                  <Group>
                    <IconCalendarEvent size="1.5rem" />
                    <Title order={3}>Upcoming Events</Title>
                  </Group>
                  {pregnancyStatus.upcomingEvents?.map((event, index) => (
                    <Text key={index}>{event}</Text>
                  )) || <Text c="dimmed">No upcoming events</Text>}
                </Stack>
              </Card>
            </Grid.Col>

            <Grid.Col span={6}>
              <Card withBorder h="100%">
                <Stack gap="md">
                  <Group>
                    <IconAlertTriangle size="1.5rem" />
                    <Title order={3}>Warning Signs</Title>
                  </Group>
                  {pregnancyStatus.warningSignsList?.map((sign, index) => (
                    <Text key={index}>{sign}</Text>
                  )) || <Text c="dimmed">No warning signs to monitor at this stage</Text>}
                </Stack>
              </Card>
            </Grid.Col>
          </Grid>

          <Stack gap="md" mt="xl">
            <Group justify="space-between">
              <Title order={3}>Pre-Foaling Signs</Title>
              <Button
                onClick={openStartDialog}
                leftSection={<IconPlus size="1rem" />}
                variant="light"
              >
                Record Sign
              </Button>
            </Group>
            <Progress
              value={pregnancyStatus.preFoalingProgress}
              size="xl"
              aria-label="Pre-foaling progress"
            />
            <Group justify="space-between" mt="xs">
              <Text size="sm">Pre-Foaling Signs Observed</Text>
              <Text size="sm">{pregnancyStatus.preFoalingSignsCount} signs</Text>
            </Group>
          </Stack>
        </>
      )}

      <StartPregnancyDialog
        opened={startDialogOpened}
        onClose={closeStartDialog}
        onConfirm={handleStartPregnancy}
        isLoading={startPregnancyMutation.isPending}
      />

      <EndPregnancyDialog
        opened={endDialogOpened}
        onClose={closeEndDialog}
        onConfirm={handleEndPregnancy}
        isLoading={endPregnancyMutation.isPending}
      />
    </Stack>
  );
}
