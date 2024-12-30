import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  TextInput,
  Select,
  NumberInput,
  Button,
  Card,
  Title,
  Stack,
  Group,
  Text,
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { IconHorse } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';

interface HorseInput {
  name: string;
  gender: 'MARE' | 'STALLION' | 'GELDING';
  breed: string;
  dateOfBirth: Date;
  weight?: number;
  conceptionDate?: Date;
  motherId?: number;
  fatherId?: number;
  externalMother?: string;
  externalFather?: string;
}

export function AddHorse() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [formData, setFormData] = useState<HorseInput>({
    name: '',
    gender: 'MARE',
    breed: '',
    dateOfBirth: new Date(),
    weight: undefined
  });

  const mutation = useMutation({
    mutationFn: async (newHorse: HorseInput) => {
      const response = await fetch('/api/horses', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...newHorse,
          dateOfBirth: newHorse.dateOfBirth.toISOString(),
        }),
      });
      
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create horse');
      }
      
      return response.json();
    },
    onSuccess: () => {
      // Invalidate the horses query to refetch the list
      queryClient.invalidateQueries({ queryKey: ['horses'] });
      notifications.show({
        title: 'Success',
        message: 'Horse added successfully',
        color: 'green'
      });
      navigate('/');
    },
    onError: (error) => {
      notifications.show({
        title: 'Error',
        message: error.message,
        color: 'red'
      });
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate(formData);
  };

  return (
    <Card withBorder radius="md" p="xl" maw={600} mx="auto">
      <Stack gap="lg">
        <Group justify="space-between" mb="md">
          <Title order={2}>
            <Group gap="xs">
              <IconHorse size={30} />
              <Text>Add New Horse</Text>
            </Group>
          </Title>
        </Group>

        <form onSubmit={handleSubmit}>
          <Stack gap="md">
            <TextInput
              label="Name"
              placeholder="Enter horse's name"
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />

            <Select
              label="Gender"
              placeholder="Select gender"
              required
              value={formData.gender}
              onChange={(value) => setFormData({ ...formData, gender: value as 'MARE' | 'STALLION' | 'GELDING' })}
              data={[
                { value: 'MARE', label: 'Mare' },
                { value: 'STALLION', label: 'Stallion' },
                { value: 'GELDING', label: 'Gelding' }
              ]}
            />

            <TextInput
              label="Breed"
              placeholder="Enter horse's breed"
              required
              value={formData.breed}
              onChange={(e) => setFormData({ ...formData, breed: e.target.value })}
            />

            <DatePickerInput
              label="Date of Birth"
              placeholder="Select date"
              required
              value={formData.dateOfBirth}
              onChange={(date) => setFormData({ ...formData, dateOfBirth: date || new Date() })}
              maxDate={new Date()}
            />

            <NumberInput
              label="Weight (kg)"
              placeholder="Enter weight"
              min={0}
              value={formData.weight}
              onChange={(value) => {
                const numValue = typeof value === 'string' ? parseFloat(value) : value;
                setFormData({ ...formData, weight: numValue || undefined });
              }}
              max={1000}
            />

            <Group justify="flex-end" mt="xl">
              <Button type="submit" loading={mutation.isPending}>
                Add Horse
              </Button>
            </Group>
          </Stack>
        </form>
      </Stack>
    </Card>
  );
}
