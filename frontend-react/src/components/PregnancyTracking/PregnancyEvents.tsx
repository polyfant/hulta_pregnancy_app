import { FiPlus, FiSearch, FiEdit, FiTrash2, FiActivity, FiCalendar, FiAlertTriangle } from 'react-icons/fi';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Card,
  Title,
  Text,
  Stack,
  Group,
  Button,
  Modal,
  TextInput,
  Textarea,
  Select,
  LoadingOverlay,
  Timeline,
  ThemeIcon
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { notifications } from '@mantine/notifications';


interface PregnancyEventsProps {
  horseId: string;
}

interface PregnancyEvent {
  id: number;
  date: string;
  type: 'checkup' | 'milestone' | 'warning';
  title: string;
  description: string;
}

export function PregnancyEvents({ horseId }: PregnancyEventsProps) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newEvent, setNewEvent] = useState<Partial<PregnancyEvent>>({
    type: 'checkup',
    date: new Date().toISOString(),
  });

  const queryClient = useQueryClient();

  const { data: events = [], isLoading } = useQuery<PregnancyEvent[]>({
    queryKey: ['pregnancy-events', horseId],
    queryFn: async () => {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/events`);
      if (!response.ok) throw new Error('Failed to fetch pregnancy events');
      return response.json();
    }
  });

  const addEventMutation = useMutation({
    mutationFn: async (event: Omit<PregnancyEvent, 'id'>) => {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(event)
      });
      if (!response.ok) throw new Error('Failed to add event');
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pregnancy-events', horseId] });
      setIsModalOpen(false);
      setNewEvent({ type: 'checkup', date: new Date().toISOString() });
      notifications.show({
        title: 'Success',
        message: 'Event added successfully',
        color: 'green'
      });
    },
    onError: (error) => {
      notifications.show({
        title: 'Error',
        message: error.message,
        color: 'red'
      });
    }
  });

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'checkup':
        return <FiActivity size="1rem" />;
      case 'milestone':
        return <FiCalendar size="1rem" />;
      case 'warning':
        return <FiAlertTriangle size="1rem" />;
      default:
        return <FiCalendar size="1rem" />;
    }
  };

  const getEventColor = (type: string) => {
    switch (type) {
      case 'checkup':
        return 'blue';
      case 'milestone':
        return 'green';
      case 'warning':
        return 'red';
      default:
        return 'gray';
    }
  };

  const handleSubmit = () => {
    if (!newEvent.title || !newEvent.type || !newEvent.date) {
      notifications.show({
        title: 'Error',
        message: 'Please fill in all required fields',
        color: 'red'
      });
      return;
    }

    addEventMutation.mutate(newEvent as Omit<PregnancyEvent, 'id'>);
  };

  return (
    <Card withBorder>
      <LoadingOverlay visible={isLoading} />
      
      <Stack gap="md">
        <Group justify="space-between" mb="md">
          <Title order={3}>Pregnancy Events</Title>
          <Button
            leftSection={<FiPlus size="1rem" />}
            onClick={() => setIsModalOpen(true)}
          >
            Add Event
          </Button>
        </Group>

        {events.length === 0 ? (
          <Text c="dimmed" ta="center" py="xl">
            No events recorded yet
          </Text>
        ) : (
          <Timeline active={events.length - 1}>
            {events.map((event) => (
              <Timeline.Item
                key={event.id}
                title={event.title}
                bullet={
                  <ThemeIcon size={24} color={getEventColor(event.type)} radius="xl">
                    {getEventIcon(event.type)}
                  </ThemeIcon>
                }
              >
                <Text size="sm" mt={4}>{event.description}</Text>
                <Text size="xs" mt={4} c="dimmed">
                  {new Date(event.date).toLocaleDateString()}
                </Text>
              </Timeline.Item>
            ))}
          </Timeline>
        )}

        <Modal
          opened={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          title="Add Pregnancy Event"
        >
          <Stack>
            <Select
              label="Event Type"
              required
              data={[
                { value: 'checkup', label: 'Checkup' },
                { value: 'milestone', label: 'Milestone' },
                { value: 'warning', label: 'Warning' }
              ]}
              value={newEvent.type}
              onChange={(value) => setNewEvent({ ...newEvent, type: value as PregnancyEvent['type'] })}
            />

            <DatePickerInput
              label="Date"
              required
              value={newEvent.date ? new Date(newEvent.date) : null}
              onChange={(date) => setNewEvent({ ...newEvent, date: date?.toISOString() })}
              maxDate={new Date()}
            />

            <TextInput
              label="Title"
              required
              value={newEvent.title || ''}
              onChange={(e) => setNewEvent({ ...newEvent, title: e.target.value })}
            />

            <Textarea
              label="Description"
              value={newEvent.description || ''}
              onChange={(e) => setNewEvent({ ...newEvent, description: e.target.value })}
            />

            <Button
              onClick={handleSubmit}
              loading={addEventMutation.isPending}
              mt="md"
            >
              Add Event
            </Button>
          </Stack>
        </Modal>
      </Stack>
    </Card>
  );
}
