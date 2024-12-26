import React, { useState, useEffect } from 'react';
import {
  Box,
  Title,
  Stack,
  Button,
  ActionIcon,
  Modal,
  Select,
  Textarea,
  Loader,
  Alert,
  Timeline,
  Text,
  Group
} from '@mantine/core';
import { IconPlus, IconCalendar } from '@tabler/icons-react';
import { format } from 'date-fns';
import { PregnancyEvent } from '../../types/pregnancy';

const eventTypes = [
  { value: 'CONCEPTION', label: 'Conception' },
  { value: 'VET_CHECK', label: 'Veterinary Check' },
  { value: 'VACCINATION', label: 'Vaccination' },
  { value: 'ULTRASOUND', label: 'Ultrasound' },
  { value: 'BEHAVIORAL_CHANGE', label: 'Behavioral Change' },
  { value: 'COMPLICATION', label: 'Complication' },
  { value: 'PRE_FOALING_CHANGES', label: 'Pre-Foaling Changes' },
  { value: 'FOALING', label: 'Foaling' },
  { value: 'POST_FOALING_CHECK', label: 'Post-Foaling Check' }
];

interface PregnancyEventsProps {
  horseId: string;
}

const PregnancyEvents: React.FC<PregnancyEventsProps> = ({ horseId }) => {
  const [events, setEvents] = useState<PregnancyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [newEvent, setNewEvent] = useState({
    eventType: '',
    description: '',
    notes: ''
  });

  const fetchEvents = async () => {
    try {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/events`);
      if (!response.ok) throw new Error('Failed to fetch events');
      const data = await response.json();
      setEvents(data);
    } catch (err) {
      setError('Failed to load pregnancy events');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [horseId]);

  const handleAddEvent = async () => {
    try {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newEvent)
      });
      
      if (!response.ok) throw new Error('Failed to add event');
      
      setOpenDialog(false);
      setNewEvent({ eventType: '', description: '', notes: '' });
      fetchEvents();
    } catch (err) {
      setError('Failed to add pregnancy event');
    }
  };

  if (loading) {
    return (
      <Box style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '200px' }}>
        <Loader />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert color="red" title="Error">
        {error}
      </Alert>
    );
  }

  return (
    <Box>
      <Group justify="space-between" mb="md">
        <Title order={3}>Pregnancy Events</Title>
        <ActionIcon 
          variant="filled" 
          color="blue" 
          onClick={() => setOpenDialog(true)}
          size="lg"
        >
          <IconPlus size={20} />
        </ActionIcon>
      </Group>

      <Timeline active={events.length} bulletSize={24}>
        {events.map((event) => (
          <Timeline.Item
            key={event.id}
            bullet={<IconCalendar size={12} />}
            title={
              <Group gap="xs">
                <Text fw={500}>
                  {eventTypes.find(t => t.value === event.eventType)?.label || event.eventType}
                </Text>
                <Text size="sm" c="dimmed">
                  {format(new Date(event.date), 'MMM d, yyyy')}
                </Text>
              </Group>
            }
          >
            <Text size="sm">{event.description}</Text>
            {event.notes && (
              <Text size="sm" c="dimmed" mt={4}>
                {event.notes}
              </Text>
            )}
          </Timeline.Item>
        ))}
      </Timeline>

      <Modal
        opened={openDialog}
        onClose={() => setOpenDialog(false)}
        title="Add Pregnancy Event"
        size="md"
      >
        <Stack>
          <Select
            label="Event Type"
            placeholder="Select event type"
            data={eventTypes}
            value={newEvent.eventType}
            onChange={(value) => setNewEvent({ ...newEvent, eventType: value || '' })}
            required
          />

          <Textarea
            label="Description"
            placeholder="Enter event description"
            value={newEvent.description}
            onChange={(event) => setNewEvent({ ...newEvent, description: event.currentTarget.value })}
            required
          />

          <Textarea
            label="Notes (Optional)"
            placeholder="Enter additional notes"
            value={newEvent.notes}
            onChange={(event) => setNewEvent({ ...newEvent, notes: event.currentTarget.value })}
            minRows={3}
          />

          <Button
            onClick={handleAddEvent}
            disabled={!newEvent.eventType || !newEvent.description}
            fullWidth
          >
            Add Event
          </Button>
        </Stack>
      </Modal>
    </Box>
  );
};

export default PregnancyEvents;
