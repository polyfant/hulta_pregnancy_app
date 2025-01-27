import { Card, Progress, Text, Stack, Group, Badge, Tooltip, Skeleton } from '@mantine/core';
import { Horse } from '@phosphor-icons/react';
import { FC } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useHorsesApi } from '../../api/horses';
import { PregnancyStage, PregnancyStatus } from '../../types/pregnancy';

interface StageVisualizationProps {
  horseId: number;
}

const stageColors: Record<PregnancyStage, string> = {
  'EARLY': 'blue',
  'MIDDLE': 'cyan',
  'LATE': 'teal',
  'NEARTERM': 'indigo',
  'FOALING': 'brand'
} as const;

const stageDescriptions: Record<PregnancyStage, string> = {
  'EARLY': 'First trimester - Critical development period',
  'MIDDLE': 'Second trimester - Steady growth phase',
  'LATE': 'Third trimester - Rapid foal development',
  'NEARTERM': 'Final preparation phase',
  'FOALING': 'Birth imminent'
} as const;

export const StageVisualization: FC<StageVisualizationProps> = ({ horseId }) => {
  const api = useHorsesApi();
  const { data: status, isLoading, error } = useQuery<PregnancyStatus>({
    queryKey: ['pregnancy', horseId],
    queryFn: () => api.getPregnancyStatus(horseId)
  });

  if (isLoading) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Stack gap="md">
          <Skeleton height={30} width="60%" />
          <Skeleton height={40} />
          <Stack gap="xs">
            <Skeleton height={20} />
            <Skeleton height={20} />
          </Stack>
        </Stack>
      </Card>
    );
  }

  if (error || !status) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Text c="red">Error loading pregnancy status</Text>
      </Card>
    );
  }

  const progress = (status.currentDay / status.totalDays) * 100;
  
  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="md">
        <Group justify="space-between" align="center">
          <div>
            <Text fw={700} size="xl">Pregnancy Progress</Text>
            <Text c="dimmed" size="sm">{status.currentDay} days completed</Text>
          </div>
          <Tooltip label={stageDescriptions[status.currentStage]}>
          <Badge 
  size="lg" 
  variant="filled" 
  color={stageColors[status.currentStage]}
  leftSection={<Horse size={16} weight="fill" />}
>
  {status.currentStage}
</Badge>
          </Tooltip>
        </Group>
        
        <Stack gap="xs">
          <Progress 
            size="xl" 
            value={progress} 
            color={stageColors[status.currentStage]}
            radius="xl"
            striped
            animated
          />
          
          <Group justify="space-between" align="center">
            <Text size="sm" c="dimmed">Start</Text>
            <Text size="sm" fw={500}>{Math.round(progress)}% Complete</Text>
            <Text size="sm" c="dimmed">Due Date</Text>
          </Group>

          <Group justify="space-between" mt="xs">
            <Text fw={500}>{status.totalDays - status.currentDay} days remaining</Text>
            <Text c="dimmed" size="sm">Due: {new Date(status.dueDate).toLocaleDateString()}</Text>
          </Group>
        </Stack>
      </Stack>
    </Card>
  );
};