import { useParams, useNavigate } from 'react-router-dom';
import { Paper, Title, LoadingOverlay } from '@mantine/core';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { HorseForm } from './HorseForm';
import { Horse, CreateHorseInput } from '../types/horse';

export function EditHorse() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: horse, isLoading } = useQuery<Horse>({
    queryKey: ['horse', id],
    queryFn: async () => {
      const response = await fetch(`/api/horses/${id}`);
      if (!response.ok) throw new Error('Failed to fetch horse details');
      return response.json();
    },
  });

  const updateMutation = useMutation({
    mutationFn: async (updatedHorse: CreateHorseInput) => {
      const response = await fetch(`/api/horses/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedHorse),
      });
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Failed to update horse');
      }
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horses'] });
      queryClient.invalidateQueries({ queryKey: ['horse', id] });
      notifications.show({
        title: 'Success',
        message: 'Horse updated successfully',
        color: 'green',
      });
      navigate(`/horses/${id}`);
    },
    onError: (error: Error) => {
      notifications.show({
        title: 'Error',
        message: error.message,
        color: 'red',
      });
    },
  });

  const handleSubmit = (data: CreateHorseInput) => {
    updateMutation.mutate(data);
  };

  if (isLoading) {
    return <LoadingOverlay visible />;
  }

  return (
    <Paper p="md">
      <Title order={2} mb="md">Edit Horse: {horse?.name}</Title>
      {horse && <HorseForm onSubmit={handleSubmit} initialValues={horse} />}
    </Paper>
  );
}

export default EditHorse;
