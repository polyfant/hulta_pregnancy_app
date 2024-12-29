import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Text,
  Group,
  Stack,
  Timeline,
  Badge,
  Title,
  LoadingOverlay
} from '@mantine/core';
import {
  IconCalendarEvent,
  IconStethoscope,
  IconBabyCarriage,
  IconAlertTriangle
} from '@tabler/icons-react';

interface PregnancyStatusProps {
  horseId: number;
}

interface PregnancyStatus {
  currentStage: string;
  nextMilestone: {
    date: string;
    event: string;
  };
  recentEvents: Array<{
    date: string;
    event: string;
    type: 'milestone' | 'checkup' | 'warning';
  }>;
  upcomingEvents: Array<{
    date: string;
    event: string;
    type: 'milestone' | 'checkup' | 'warning';
  }>;
}

export function PregnancyStatus({ horseId }: PregnancyStatusProps) {
  const { data, isLoading } = useQuery<PregnancyStatus>({
    queryKey: ['pregnancy-status', horseId],
    queryFn: async () => {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/status`);
      if (!response.ok) throw new Error('Failed to fetch pregnancy status');
      return response.json();
    }
  });

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'milestone':
        return <IconCalendarEvent size="1.2rem" />;
      case 'checkup':
        return <IconStethoscope size="1.2rem" />;
      case 'warning':
        return <IconAlertTriangle size="1.2rem" />;
      default:
        return <IconBabyCarriage size="1.2rem" />;
    }
  };

  const getEventColor = (type: string) => {
    switch (type) {
      case 'milestone':
        return 'blue';
      case 'checkup':
        return 'green';
      case 'warning':
        return 'yellow';
      default:
        return 'gray';
    }
  };

  if (isLoading) {
    return (
      <Card withBorder radius="md" pos="relative" h={400}>
        <LoadingOverlay visible />
      </Card>
    );
  }

  if (!data) {
    return (
      <Card withBorder radius="md">
        <Text c="dimmed">No pregnancy status available</Text>
      </Card>
    );
  }

  return (
    <Stack gap="lg">
      <Card withBorder radius="md">
        <Stack gap="xs">
          <Group justify="space-between">
            <Text c="dimmed">Current Stage:</Text>
            <Badge size="lg">{data.currentStage}</Badge>
          </Group>
          <Group justify="space-between">
            <Text c="dimmed">Next Milestone:</Text>
            <Group spacing="xs">
              <Text>{new Date(data.nextMilestone.date).toLocaleDateString()}</Text>
              <Text>-</Text>
              <Text>{data.nextMilestone.event}</Text>
            </Group>
          </Group>
        </Stack>
      </Card>

      <Card withBorder radius="md">
        <Title order={3} mb="md">Timeline</Title>
        <Timeline active={data.recentEvents.length - 1} bulletSize={24}>
          {data.recentEvents.map((event, index) => (
            <Timeline.Item
              key={index}
              bullet={getEventIcon(event.type)}
              title={event.event}
              color={getEventColor(event.type)}
            >
              <Text size="sm" mt={4}>
                {new Date(event.date).toLocaleDateString()}
              </Text>
            </Timeline.Item>
          ))}
        </Timeline>
      </Card>

      <Card withBorder radius="md">
        <Title order={3} mb="md">Upcoming Events</Title>
        <Timeline bulletSize={24}>
          {data.upcomingEvents.map((event, index) => (
            <Timeline.Item
              key={index}
              bullet={getEventIcon(event.type)}
              title={event.event}
              color={getEventColor(event.type)}
            >
              <Text size="sm" mt={4}>
                {new Date(event.date).toLocaleDateString()}
              </Text>
            </Timeline.Item>
          ))}
        </Timeline>
      </Card>
    </Stack>
  );
}
