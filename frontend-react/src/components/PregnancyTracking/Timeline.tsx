import { Timeline as MantineTimeline, Text, Card, Stack, ThemeIcon } from '@mantine/core';
import { Horse, Baby, FirstAid, Calendar } from '@phosphor-icons/react';
import { FC } from 'react';

interface TimelineEvent {
  id: string;
  title: string;
  description: string;
  date: string;
  icon: JSX.Element;
  color: string;
}

export const PregnancyTimeline: FC = () => {
  const events: TimelineEvent[] = [
    {
      id: '1',
      title: 'Early Stage',
      description: 'Critical development period. Regular check-ups essential.',
      date: 'Days 1-100',
      icon: <Horse size={16} weight="fill" />,
      color: 'blue'
    },
    {
      id: '2',
      title: 'Mid Stage',
      description: 'Steady growth phase. Monitor nutrition carefully.',
      date: 'Days 101-240',
      icon: <FirstAid size={16} weight="fill" />,
      color: 'cyan'
    },
    {
      id: '3',
      title: 'Late Stage',
      description: 'Rapid foal development. Prepare for pre-foaling.',
      date: 'Days 241-320',
      icon: <Calendar size={16} weight="fill" />,
      color: 'teal'
    },
    {
      id: '4',
      title: 'Pre-foaling',
      description: 'Final preparations. Monitor closely for signs.',
      date: 'Days 321-340',
      icon: <Baby size={16} weight="fill" />,
      color: 'brand'
    }
  ];

  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="lg">
        <Text fw={700} size="xl">Pregnancy Timeline</Text>
        
        <MantineTimeline active={1} bulletSize={24} lineWidth={2}>
          {events.map((event) => (
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