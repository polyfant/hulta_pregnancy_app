import { Card, Progress, Text, Stack, Group, Badge, Tooltip, Skeleton } from '@mantine/core';
import { Horse } from '@phosphor-icons/react';
import { FC } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useHorsesApi } from '../../api/horses';
import { PregnancyStatus } from '../../types/pregnancy';

interface StageVisualizationProps {
  horseId: number;
}

const stageColors = {
  'Early': 'blue',
  'Mid': 'cyan',
  'Late': 'teal',
  'Pre-foaling': 'brand'
} as const;

const stageDescriptions = {
  'Early': 'First trimester - Critical development period',
  'Mid': 'Second trimester - Steady growth phase',
  'Late': 'Third trimester - Rapid foal development',
  'Pre-foaling': 'Final preparation for birth'
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
          <Tooltip label={stageDescriptions[status.stage]}>
            <Badge 
              size="lg" 
              variant="filled" 
              color={stageColors[status.stage]}
              leftSection={<Horse size={16} weight="fill" />}
            >
              {status.stage}
            </Badge>
          </Tooltip>
        </Group>
        
        <Stack gap="xs">
          <Progress 
            size="xl" 
            value={progress} 
            color={stageColors[status.stage]}
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