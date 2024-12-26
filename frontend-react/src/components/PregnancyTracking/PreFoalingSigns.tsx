import React, { useState, useEffect } from 'react';
import {
  Box,
  Title,
  Paper,
  Stack,
  Switch,
  Text,
  Alert,
  Loader,
  Group,
  Badge,
  Textarea,
  Button,
  Modal
} from '@mantine/core';
import { format } from 'date-fns';
import { PreFoalingSign } from '../../types/pregnancy';

interface PreFoalingSignsProps {
  horseId: string;
}

const PreFoalingSigns: React.FC<PreFoalingSignsProps> = ({ horseId }) => {
  const [signs, setSigns] = useState<PreFoalingSign[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSign, setSelectedSign] = useState<PreFoalingSign | null>(null);
  const [notes, setNotes] = useState('');
  const [openDialog, setOpenDialog] = useState(false);

  const fetchSigns = async () => {
    try {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/foaling-signs`);
      if (!response.ok) throw new Error('Failed to fetch pre-foaling signs');
      const data = await response.json();
      setSigns(data);
    } catch (err) {
      setError('Failed to load pre-foaling signs');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSigns();
  }, [horseId]);

  const handleToggleSign = async (sign: PreFoalingSign) => {
    if (!sign.observed) {
      setSelectedSign(sign);
      setOpenDialog(true);
    } else {
      await updateSign(sign, false);
    }
  };

  const handleSubmitNotes = async () => {
    if (selectedSign) {
      await updateSign(selectedSign, true);
      setOpenDialog(false);
      setSelectedSign(null);
      setNotes('');
    }
  };

  const updateSign = async (sign: PreFoalingSign, observed: boolean) => {
    try {
      const response = await fetch(`/api/horses/${horseId}/pregnancy/foaling-signs/${sign.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          observed,
          notes: observed ? notes : undefined,
          dateObserved: observed ? new Date().toISOString() : undefined
        })
      });
      
      if (!response.ok) throw new Error('Failed to update pre-foaling sign');
      
      fetchSigns();
    } catch (err) {
      setError('Failed to update pre-foaling sign');
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
      <Title order={3} mb="md">Pre-Foaling Signs</Title>

      <Stack>
        {signs.map((sign) => (
          <Paper key={sign.id} p="md" withBorder>
            <Group position="apart">
              <Box>
                <Group spacing="xs">
                  <Text fw={500}>{sign.signName}</Text>
                  {sign.observed && (
                    <Badge color="green">
                      Observed on {format(new Date(sign.dateObserved!), 'MMM d, yyyy')}
                    </Badge>
                  )}
                </Group>
                {sign.notes && (
                  <Text size="sm" c="dimmed" mt={4}>
                    {sign.notes}
                  </Text>
                )}
              </Box>
              <Switch
                checked={sign.observed}
                onChange={() => handleToggleSign(sign)}
                size="lg"
              />
            </Group>
          </Paper>
        ))}
      </Stack>

      <Modal
        opened={openDialog}
        onClose={() => {
          setOpenDialog(false);
          setSelectedSign(null);
          setNotes('');
        }}
        title="Add Notes"
        size="md"
      >
        <Stack>
          <Text>
            Add notes for: {selectedSign?.signName}
          </Text>

          <Textarea
            label="Notes"
            placeholder="Enter any additional observations or notes"
            value={notes}
            onChange={(event) => setNotes(event.currentTarget.value)}
            minRows={3}
          />

          <Button onClick={handleSubmitNotes} fullWidth>
            Save
          </Button>
        </Stack>
      </Modal>
    </Box>
  );
};

export default PreFoalingSigns;
