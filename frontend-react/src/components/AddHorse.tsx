import { useState, useEffect } from 'react';
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
  Switch,
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

const AddHorse = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [formData, setFormData] = useState<HorseInput>({
    name: '',
    gender: 'MARE',
    breed: '',
    dateOfBirth: new Date(),
    weight: undefined,
    conceptionDate: undefined,
    fatherId: undefined,
    externalFather: ''
  });

  const [useExternalFather, setUseExternalFather] = useState(false);
  const [availableStallions, setAvailableStallions] = useState<Array<{ value: string, label: string }>>([]);

  useEffect(() => {
    const fetchStallions = async () => {
      try {
        const response = await fetch('/api/horses');
        if (!response.ok) throw new Error('Failed to fetch horses');
        const horses = await response.json();
        const stallions = horses
          .filter((horse: any) => horse.gender === 'STALLION')
          .map((horse: any) => ({
            value: horse.id.toString(),
            label: horse.name
          }));
        setAvailableStallions(stallions);
      } catch (error) {
        console.error('Error fetching stallions:', error);
      }
    };
    fetchStallions();
  }, []);

  const mutation = useMutation({
    mutationFn: async (data: HorseInput) => {
      const response = await fetch('/api/horses', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        throw new Error('Failed to add horse');
      }
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horses'] });
      notifications.show({
        title: 'Success',
        message: 'Horse added successfully',
        color: 'green',
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

            {formData.gender === 'MARE' && (
              <>
                <DatePickerInput
                  label="Conception Date"
                  placeholder="Select date if pregnant"
                  value={formData.conceptionDate}
                  onChange={(date) => setFormData({ ...formData, conceptionDate: date || undefined })}
                  maxDate={new Date()}
                />

                <Switch
                  label="External Father"
                  checked={useExternalFather}
                  onChange={(event) => {
                    setUseExternalFather(event.currentTarget.checked);
                    setFormData({
                      ...formData,
                      fatherId: undefined,
                      externalFather: ''
                    });
                  }}
                />

                {useExternalFather ? (
                  <TextInput
                    label="Father's Name (External)"
                    placeholder="Enter external father's name"
                    value={formData.externalFather}
                    onChange={(e) => setFormData({ ...formData, externalFather: e.target.value })}
                  />
                ) : (
                  <Select
                    label="Father"
                    placeholder="Select father"
                    data={availableStallions}
                    value={formData.fatherId?.toString()}
                    onChange={(value) => setFormData({ 
                      ...formData, 
                      fatherId: value ? parseInt(value) : undefined 
                    })}
                    clearable
                  />
                )}
              </>
            )}

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

export default AddHorse;
