import { FiPlus } from 'react-icons/md';
import { FiPlus, FiSearch, FiEdit, FiTrash2 } from 'react-icons/fi';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import {
  Card,
  Stack,
  Group,
  Title,
  Text,
  Button,
  Progress,
  Grid,
  LoadingOverlay
} from '@mantine/core';

import { notifications } from '@mantine/notifications';
import { Horse } from '../../types/horse';
import { PregnancyStatus, PregnancyEvent } from '../../types/pregnancy';
import { StartPregnancyDialog } from './StartPregnancyDialog';
import { EndPregnancyDialog } from './EndPregnancyDialog';
import { formatDate } from '../../utils/dateUtils';

interface PregnancyTrackingProps {}

const STAGES = {
  EARLY: { label: 'Early Stage', progress: 25 },
  MIDDLE: { label: 'Middle Stage', progress: 50 },
  LATE: { label: 'Late Stage', progress: 75 },
  NEARTERM: { label: 'Near Term', progress: 90 },
  FOALING: { label: 'Foaling', progress: 100 }
};

export default function PregnancyTracking({}: PregnancyTrackingProps) {
  const { id } = useParams();
  const [startDialogOpened, setStartDialogOpen] = useState(false);
  const [endDialogOpened, setEndDialogOpen] = useState(false);
  const queryClient = useQueryClient();

  const { data: horse, isLoading: horseLoading } = useQuery<Horse>({
    queryKey: ['horse', id],
    enabled: !!id,
  });

  const { data: pregnancyStatus } = useQuery<PregnancyStatus>({
    queryKey: ['pregnancyStatus', id],
    enabled: !!id,
  });

  if (horseLoading) {
    return <LoadingOverlay visible />;
  }

  if (!horse) {
    return <Text>Horse not found</Text>;
  }

  const startPregnancyMutation = useMutation({
    mutationFn: async (date: string) => {
      const response = await fetch(`/api/horses/${id}/pregnancy/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ conceptionDate: date }),
      });
      if (!response.ok) throw new Error('Failed to start pregnancy tracking');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horse', id] });
      queryClient.invalidateQueries({ queryKey: ['pregnancyStatus', id] });
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
      const response = await fetch(`/api/horses/${id}/pregnancy/end`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error('Failed to end pregnancy tracking');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horse', id] });
      queryClient.invalidateQueries({ queryKey: ['pregnancyStatus', id] });
      notifications.show({
        title: 'Success',
        message: 'Ended pregnancy tracking',
        color: 'green',
      });
      setEndDialogOpen(false);
    },
  });

  const getStageProgress = (stage: PregnancyStatus['currentStage']) => {
    return STAGES[stage]?.progress || 0;
  };

  return (
    <Stack spacing="lg">
      <Card withBorder>
        <Stack>
          <Group position="apart">
            <Title order={2}>Pregnancy Tracking</Title>
            {pregnancyStatus?.isPregnant ? (
              <Button
                color="red"
                leftIcon={<IconX />}
                onClick={() => setEndDialogOpen(true)}
              >
                End Pregnancy
              </Button>
            ) : (
              <Button
                color="blue"
                leftIcon={<FiPlus />}
                onClick={() => setStartDialogOpen(true)}
              >
                Start Pregnancy
              </Button>
            )}
          </Group>

          {pregnancyStatus?.isPregnant && (
            <>
              <Progress
                value={getStageProgress(pregnancyStatus.currentStage)}
                label={STAGES[pregnancyStatus.currentStage]?.label}
                size="xl"
                radius="xl"
              />
              <Grid>
                <Grid.Col span={6}>
                  <Stack spacing="xs">
                    <Text size="sm" color="dimmed">Conception Date</Text>
                    <Text>{formatDate(pregnancyStatus.conceptionDate)}</Text>
                  </Stack>
                </Grid.Col>
                <Grid.Col span={6}>
                  <Stack spacing="xs">
                    <Text size="sm" color="dimmed">Expected Due Date</Text>
                    <Text>{formatDate(pregnancyStatus.expectedDueDate)}</Text>
                  </Stack>
                </Grid.Col>
                <Grid.Col span={12}>
                  <Stack spacing="xs">
                    <Text size="sm" color="dimmed">Days in Pregnancy</Text>
                    <Text>{pregnancyStatus.daysInPregnancy} days</Text>
                  </Stack>
                </Grid.Col>
              </Grid>
            </>
          )}
        </Stack>
      </Card>

      <StartPregnancyDialog
        opened={startDialogOpened}
        onClose={() => setStartDialogOpen(false)}
        onSubmit={(date) => startPregnancyMutation.mutate(date)}
        isLoading={startPregnancyMutation.isPending}
      />

      <EndPregnancyDialog
        opened={endDialogOpened}
        onClose={() => setEndDialogOpen(false)}
        onSubmit={(data) => endPregnancyMutation.mutate(data)}
        isLoading={endPregnancyMutation.isPending}
      />
    </Stack>
  );
}

