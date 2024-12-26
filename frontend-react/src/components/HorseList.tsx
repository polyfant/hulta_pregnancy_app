import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Table,
  Group,
  Button,
  Modal,
  TextInput,
  Select,
  Stack,
  Text,
  ActionIcon
} from '@mantine/core';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { IconEdit, IconTrash } from '@tabler/icons-react';

interface Horse {
  id: number;
  name: string;
  gender: 'MARE' | 'STALLION' | 'GELDING';
  breed: string;
  birthYear: number;
  color: string;
}

interface CreateHorseInput {
  name: string;
  gender: 'MARE' | 'STALLION' | 'GELDING';
  breed: string;
  birthYear: number;
  color: string;
}

const HorseList: React.FC = () => {
  const queryClient = useQueryClient();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState<CreateHorseInput>({
    name: '',
    gender: 'MARE',
    breed: '',
    birthYear: new Date().getFullYear(),
    color: ''
  });

  const { data: horses, isLoading, error } = useQuery<Horse[]>({
    queryKey: ['horses'],
    queryFn: async () => {
      const response = await fetch('/api/horses');
      if (!response.ok) {
        throw new Error('Failed to fetch horses');
      }
      return response.json();
    }
  });

  const createHorseMutation = useMutation({
    mutationFn: async (newHorse: CreateHorseInput) => {
      const response = await fetch('/api/horses', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newHorse),
      });
      if (!response.ok) {
        throw new Error('Failed to create horse');
      }
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horses'] });
      setIsModalOpen(false);
      setFormData({
        name: '',
        gender: 'MARE',
        breed: '',
        birthYear: new Date().getFullYear(),
        color: ''
      });
    },
  });

  const deleteHorseMutation = useMutation({
    mutationFn: async (horseId: number) => {
      const response = await fetch(`/api/horses/${horseId}`, {
        method: 'DELETE',
      });
      if (!response.ok) {
        throw new Error('Failed to delete horse');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['horses'] });
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createHorseMutation.mutate(formData);
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {(error as Error).message}</div>;

  return (
    <div>
      <Group justify="space-between" mb="md">
        <Text size="xl" fw={700}>Horses</Text>
        <Button onClick={() => setIsModalOpen(true)}>Add Horse</Button>
      </Group>

      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Gender</Table.Th>
            <Table.Th>Breed</Table.Th>
            <Table.Th>Birth Year</Table.Th>
            <Table.Th>Color</Table.Th>
            <Table.Th>Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {horses?.map((horse) => (
            <Table.Tr key={horse.id}>
              <Table.Td>{horse.name}</Table.Td>
              <Table.Td>{horse.gender}</Table.Td>
              <Table.Td>{horse.breed}</Table.Td>
              <Table.Td>{horse.birthYear}</Table.Td>
              <Table.Td>{horse.color}</Table.Td>
              <Table.Td>
                <Group>
                  <ActionIcon 
                    component={Link} 
                    to={`/horses/${horse.id}`}
                    variant="subtle"
                    color="blue"
                  >
                    <IconEdit size={16} />
                  </ActionIcon>
                  <ActionIcon
                    onClick={() => deleteHorseMutation.mutate(horse.id)}
                    variant="subtle"
                    color="red"
                  >
                    <IconTrash size={16} />
                  </ActionIcon>
                  {horse.gender === 'MARE' && (
                    <Button
                      component={Link}
                      to={`/horses/${horse.id}/pregnancy`}
                      size="xs"
                      variant="light"
                    >
                      Pregnancy
                    </Button>
                  )}
                </Group>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>

      <Modal
        opened={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Add New Horse"
      >
        <form onSubmit={handleSubmit}>
          <Stack>
            <TextInput
              label="Name"
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />

            <Select
              label="Gender"
              required
              data={[
                { value: 'MARE', label: 'Mare' },
                { value: 'STALLION', label: 'Stallion' },
                { value: 'GELDING', label: 'Gelding' }
              ]}
              value={formData.gender}
              onChange={(value) => setFormData({ ...formData, gender: value as 'MARE' | 'STALLION' | 'GELDING' })}
            />

            <TextInput
              label="Breed"
              required
              value={formData.breed}
              onChange={(e) => setFormData({ ...formData, breed: e.target.value })}
            />

            <TextInput
              label="Birth Year"
              type="number"
              required
              value={formData.birthYear.toString()}
              onChange={(e) => setFormData({ ...formData, birthYear: parseInt(e.target.value) })}
            />

            <TextInput
              label="Color"
              required
              value={formData.color}
              onChange={(e) => setFormData({ ...formData, color: e.target.value })}
            />

            <Button type="submit" loading={createHorseMutation.isPending}>
              Add Horse
            </Button>
          </Stack>
        </form>
      </Modal>
    </div>
  );
};

export default HorseList;
