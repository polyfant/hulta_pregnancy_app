import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Table, Button, Group, Text, Modal } from '@mantine/core';
import { horsesApi } from '../api/horses';
import { useState } from 'react';
import { Horse } from '../types/horse';
import { AddHorseForm } from './AddHorseForm';

export function HorseList() {
    const queryClient = useQueryClient();
    const [editingHorse, setEditingHorse] = useState<Horse | null>(null);

    const { data: horses, isLoading, error } = useQuery({
        queryKey: ['horses'],
        queryFn: horsesApi.getAll,
    });

    const deleteMutation = useMutation({
        mutationFn: horsesApi.delete,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['horses'] });
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }) => horsesApi.update(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['horses'] });
            setEditingHorse(null);
        },
    });

    if (isLoading) return <Text>Loading horses...</Text>;
    if (error) return <Text color="red">Error loading horses</Text>;

    return (
        <>
            <Modal 
                opened={!!editingHorse} 
                onClose={() => setEditingHorse(null)}
                title="Edit Horse"
            >
                {editingHorse && (
                    <AddHorseForm
                        initialValues={editingHorse}
                        onSubmit={(values) => updateMutation.mutate({ id: editingHorse.id, data: values })}
                        submitButtonText="Update Horse"
                    />
                )}
            </Modal>

            <Table>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Name</Table.Th>
                        <Table.Th>Breed</Table.Th>
                        <Table.Th>Gender</Table.Th>
                        <Table.Th>Date of Birth</Table.Th>
                        <Table.Th>Actions</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {horses?.map((horse) => (
                        <Table.Tr key={horse.id}>
                            <Table.Td>{horse.name}</Table.Td>
                            <Table.Td>{horse.breed}</Table.Td>
                            <Table.Td>{horse.gender}</Table.Td>
                            <Table.Td>{horse.dateOfBirth}</Table.Td>
                            <Table.Td>
                                <Group>
                                    <Button 
                                        variant="outline"
                                        size="xs"
                                        onClick={() => setEditingHorse(horse)}
                                    >
                                        Edit
                                    </Button>
                                    <Button 
                                        variant="outline" 
                                        color="red" 
                                        size="xs"
                                        onClick={() => deleteMutation.mutate(horse.id)}
                                    >
                                        Delete
                                    </Button>
                                </Group>
                            </Table.Td>
                        </Table.Tr>
                    ))}
                </Table.Tbody>
            </Table>
        </>
    );
}
