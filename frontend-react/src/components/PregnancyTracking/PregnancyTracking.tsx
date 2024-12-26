import React, { useState, useEffect } from 'react';
import { 
  Box,
  Title,
  Paper,
  Grid,
  Button,
  Loader,
  Alert,
  Notification,
  Container
} from '@mantine/core';
import { useParams } from 'react-router-dom';
import PregnancyStatus from './PregnancyStatus';
import PregnancyEvents from './PregnancyEvents';
import PregnancyGuidelines from './PregnancyGuidelines';
import PreFoalingSigns from './PreFoalingSigns';
import StartPregnancyDialog from './StartPregnancyDialog';
import EndPregnancyDialog from './EndPregnancyDialog';

const PregnancyTracking: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState<any>(null);
  const [openStartDialog, setOpenStartDialog] = useState(false);
  const [openEndDialog, setOpenEndDialog] = useState(false);
  const [notification, setNotification] = useState<{ show: boolean; message: string; type: 'success' | 'error' }>({
    show: false,
    message: '',
    type: 'success'
  });

  const fetchPregnancyStatus = async () => {
    if (!id) return;
    try {
      const response = await fetch(`/api/horses/${id}/pregnancy/status`);
      if (!response.ok) throw new Error('Failed to fetch pregnancy status');
      const data = await response.json();
      setStatus(data);
    } catch (err) {
      setError('Failed to load pregnancy status');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPregnancyStatus();
  }, [id]);

  const handleStartPregnancy = async (conceptionDate: string) => {
    if (!id) return;
    try {
      const response = await fetch(`/api/horses/${id}/pregnancy/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ conceptionDate })
      });
      
      if (!response.ok) throw new Error('Failed to start pregnancy tracking');
      
      setNotification({
        show: true,
        message: 'Successfully started pregnancy tracking',
        type: 'success'
      });
      fetchPregnancyStatus();
    } catch (err) {
      setNotification({
        show: true,
        message: 'Failed to start pregnancy tracking',
        type: 'error'
      });
    }
    setOpenStartDialog(false);
  };

  const handleEndPregnancy = async (outcome: string) => {
    if (!id) return;
    try {
      const response = await fetch(`/api/horses/${id}/pregnancy/end`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ outcome })
      });
      
      if (!response.ok) throw new Error('Failed to end pregnancy tracking');
      
      setNotification({
        show: true,
        message: 'Successfully ended pregnancy tracking',
        type: 'success'
      });
      fetchPregnancyStatus();
    } catch (err) {
      setNotification({
        show: true,
        message: 'Failed to end pregnancy tracking',
        type: 'error'
      });
    }
    setOpenEndDialog(false);
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
    <Container size="xl" p="md">
      <Grid>
        <Grid.Col span={12}>
          <Box style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
            <Title order={2}>Pregnancy Tracking</Title>
            {status?.isPregnant ? (
              <Button 
                variant="filled" 
                color="red"
                onClick={() => setOpenEndDialog(true)}
              >
                End Pregnancy Tracking
              </Button>
            ) : (
              <Button 
                variant="filled"
                onClick={() => setOpenStartDialog(true)}
              >
                Start Pregnancy Tracking
              </Button>
            )}
          </Box>
        </Grid.Col>

        {status?.isPregnant && (
          <>
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Paper shadow="xs" p="md">
                <PregnancyStatus status={status} />
              </Paper>
            </Grid.Col>

            <Grid.Col span={{ base: 12, md: 6 }}>
              <Paper shadow="xs" p="md">
                <PregnancyGuidelines horseId={id} />
              </Paper>
            </Grid.Col>

            <Grid.Col span={{ base: 12, md: 6 }}>
              <Paper shadow="xs" p="md">
                <PregnancyEvents horseId={id} />
              </Paper>
            </Grid.Col>

            <Grid.Col span={{ base: 12, md: 6 }}>
              <Paper shadow="xs" p="md">
                <PreFoalingSigns horseId={id} />
              </Paper>
            </Grid.Col>
          </>
        )}
      </Grid>

      <StartPregnancyDialog
        opened={openStartDialog}
        onClose={() => setOpenStartDialog(false)}
        onSubmit={handleStartPregnancy}
      />

      <EndPregnancyDialog
        opened={openEndDialog}
        onClose={() => setOpenEndDialog(false)}
        onSubmit={handleEndPregnancy}
      />

      {notification.show && (
        <Notification
          color={notification.type === 'success' ? 'green' : 'red'}
          onClose={() => setNotification({ ...notification, show: false })}
          style={{
            position: 'fixed',
            bottom: '1rem',
            right: '1rem',
            zIndex: 1000
          }}
        >
          {notification.message}
        </Notification>
      )}
    </Container>
  );
};

export default PregnancyTracking;
