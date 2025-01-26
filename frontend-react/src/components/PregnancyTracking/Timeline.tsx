import { Timeline as MantineTimeline, Text, Card, Stack, ThemeIcon, Skeleton } from '@mantine/core';
import { Horse, Baby, FirstAid, Calendar } from '@phosphor-icons/react';
import { FC } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useHorsesApi } from '../../api/horses';

interface TimelineEvent {
  id: string;
  title: string;
  description: string;
  date: string;
  icon: JSX.Element;
  color: string;
}

interface PregnancyTimelineProps {
  horseId: number;
}

export const PregnancyTimeline: FC<PregnancyTimelineProps> = ({ horseId }) => {
  const api = useHorsesApi();
  const { data: events, isLoading, error } = useQuery<TimelineEvent[]>({
    queryKey: ['pregnancyEvents', horseId],
    queryFn: () => api.getPregnancyEvents(horseId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });

  if (isLoading) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Stack gap="lg">
          <Skeleton height={30} width="60%" />
          <Stack gap="md">
            {[1, 2, 3, 4].map((item) => (
              <Skeleton key={item} height={80} />
            ))}
          </Stack>
        </Stack>
      </Card>
    );
  }

  if (error || !events || events.length === 0) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Text c="red">No pregnancy events found</Text>
      </Card>
    );
  }

  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="lg">
        <Text fw={700} size="xl">Pregnancy Timeline</Text>
        
        <MantineTimeline active={events.length} bulletSize={24} lineWidth={2}>
          {events.map((event, index) => (
            <MantineTimeline.Item
              key={event.id}
              bullet={
                <ThemeIcon size={24} radius="xl" color={event.color}>
                  {event.icon}
                </ThemeIcon>
              }
              title={
                <Text fw={500} size="md">
                  {event.title}
                </Text>
              }
            >
              <Text size="sm" c="dimmed" mb={4}>
                {event.date}
              </Text>
              <Text size="sm">
                {event.description}
              </Text>
            </MantineTimeline.Item>
          ))}
        </MantineTimeline>
      </Stack>
    </Card>
  );
};